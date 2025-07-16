import libvirt
import json
import requests
import time
import logging
import os
import redis
from datetime import datetime

# Configuração de log
LOG_FILE = '/var/log/bifrost/bifrost-agent.log'
os.makedirs('/var/log/bifrost', exist_ok=True)
logging.basicConfig(filename=LOG_FILE, level=logging.INFO,
                    format='%(asctime)s [%(levelname)s] %(message)s')
logger = logging.getLogger()

# Variáveis de ambiente
API_URL = os.getenv('BIFROST_API_URL', 'http://localhost:8080/api/v1/vms')
REDIS_HOST = os.getenv('REDIS_HOST', 'localhost')
REDIS_PORT = int(os.getenv('REDIS_PORT', '6379'))
REDIS_CHANNEL = os.getenv('REDIS_CHANNEL', 'vm-actions')

# Conexão com Redis
redis_client = redis.Redis(host=REDIS_HOST, port=REDIS_PORT)
try:
    redis_client.ping()
    logger.info(f"Conectado ao Redis em {REDIS_HOST}:{REDIS_PORT}")
except redis.exceptions.ConnectionError as e:
    logger.error(f"Erro ao conectar ao Redis: {e}")
    exit(1)

def coletar_dados_vm(conn):
    vms = []
    for dom in conn.listAllDomains():
        vm_info = {
            "name": dom.name(),
            "uuid": dom.UUIDString(),
            "state": dom.state()[0],
            "cpu_allocation": dom.maxVcpus(),
            "memory_allocation": dom.maxMemory(),
            "disks": [],
            "interfaces": [],
            "metadata": {},
        }

        # Discos
        for disk in dom.XMLDesc().split('<disk type=')[1:]:
            if 'file=' in disk:
                disk_path = disk.split('file=')[1].split("'")[1]
                vm_info["disks"].append({"path": disk_path})

        # Interfaces (somente se VM estiver ativa)
        if dom.isActive() == 1:
            try:
                iface_addrs = dom.interfaceAddresses(libvirt.VIR_DOMAIN_INTERFACE_ADDRESSES_SRC_LEASE, 0)
                for iface in iface_addrs.values():
                    vm_info["interfaces"].append({
                        "mac": iface.get('hwaddr'),
                        "addrs": iface.get('addrs', [])
                    })
            except libvirt.libvirtError as e:
                logger.warning(f"Não foi possível coletar interfaces de {dom.name()}: {e}")

        vms.append(vm_info)
    return vms

def enviar_dados_api(vms):
    payload = {
        "timestamp": datetime.now().astimezone().isoformat(),
        "vms": vms
    }
    try:
        resp = requests.post(API_URL, json=payload)
        if resp.status_code == 200:
            logger.info("Inventário enviado com sucesso.")
        else:
            logger.error(f"Falha ao enviar inventário: {resp.status_code}")
    except Exception as e:
        logger.error(f"Erro ao enviar inventário: {e}")

def executar_acoes_pendentes(conn):
    try:
        # Lê mensagens do canal Redis (publish/subscribe)
        pubsub = redis_client.pubsub()
        pubsub.subscribe(REDIS_CHANNEL)
        logger.info(f"Escutando canal Redis '{REDIS_CHANNEL}'...")

        for message in pubsub.listen():
            if message['type'] != 'message':
                continue

            action_data = json.loads(message['data'])
            uuid = action_data.get('uuid')
            action = action_data.get('action')

            if not uuid or not action:
                logger.warning("Mensagem incompleta recebida no Redis.")
                continue

            dom = conn.lookupByUUIDString(uuid)

            if action == "start":
                if dom.isActive() == 0:
                    dom.create()
                    logger.info(f"VM {dom.name()} iniciada com sucesso.")
                    time.sleep(5)  # aguarda VM subir antes de interagir
                else:
                    logger.info(f"VM {dom.name()} já está em execução.")

            elif action == "stop":
                if dom.isActive() == 1:
                    try:
                        dom.shutdown()
                        logger.info(f"VM {dom.name()} desligada com sucesso.")
                    except libvirt.libvirtError as e:
                        logger.warning(f"Shutdown falhou para {dom.name()}, tentando destroy: {e}")
                        dom.destroy()
                        logger.info(f"VM {dom.name()} forçada com destroy().")
                else:
                    logger.info(f"VM {dom.name()} já está desligada.")
            else:
                logger.warning(f"Ação desconhecida recebida: {action}")

    except Exception as e:
        logger.error(f"Erro ao executar ações pendentes: {e}")

def tarefa_principal():
    conn = None
    try:
        conn = libvirt.open(None)
        if conn is None:
            logger.error("Não foi possível conectar ao libvirt.")
            return

        logger.info("Coletando dados das VMs...")
        vms = coletar_dados_vm(conn)
        enviar_dados_api(vms)

    except Exception as e:
        logger.error(f"Erro na tarefa principal: {e}")
    finally:
        if conn:
            conn.close()

if __name__ == "__main__":
    logger.info("Iniciando Bifrost Agent...")

    # Roda tarefa inicial uma vez
    tarefa_principal()

    # Inicia loop infinito escutando Redis
    try:
        conn = libvirt.open(None)
        if conn is None:
            logger.error("Não foi possível conectar ao libvirt para escutar Redis.")
            exit(1)

        executar_acoes_pendentes(conn)

    except KeyboardInterrupt:
        logger.info("Bifrost Agent finalizado pelo usuário.")
    except Exception as e:
        logger.error(f"Erro fatal no agent: {e}")
    finally:
        if conn:
            conn.close()

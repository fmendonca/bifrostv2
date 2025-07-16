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
REDIS_DB = int(os.getenv('REDIS_DB', '0'))

# Conexão Redis
try:
    redis_client = redis.Redis(host=REDIS_HOST, port=REDIS_PORT, db=REDIS_DB)
    redis_client.ping()
    logger.info(f"Conectado ao Redis em {REDIS_HOST}:{REDIS_PORT}")
except Exception as e:
    logger.error(f"Erro ao conectar no Redis: {e}")
    redis_client = None

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

        # Interfaces (somente se VM estiver rodando)
        if dom.isActive() == 1:
            try:
                iface_addrs = dom.interfaceAddresses(libvirt.VIR_DOMAIN_INTERFACE_ADDRESSES_SRC_LEASE, 0)
                for iface in iface_addrs.values():
                    vm_info["interfaces"].append({
                        "mac": iface.get('hwaddr'),
                        "addrs": iface.get('addrs', [])
                    })
            except libvirt.libvirtError:
                pass  # nem sempre disponível

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
    if not redis_client:
        logger.warning("Redis não conectado, pulando verificação de ações.")
        return

    try:
        for key in redis_client.scan_iter("vm-action:*"):
            uuid = key.decode().split(":")[1]
            action = redis_client.get(key).decode()

            dom = conn.lookupByUUIDString(uuid)
            if action == "start":
                if dom.isActive() == 0:
                    dom.create()
                    logger.info(f"VM {dom.name()} iniciada com sucesso.")
                else:
                    logger.info(f"VM {dom.name()} já está em execução.")
            elif action == "stop":
                if dom.isActive() == 1:
                    try:
                        dom.shutdown()
                        logger.info(f"VM {dom.name()} desligada com sucesso.")
                    except libvirt.libvirtError as e:
                        if "domain is not running" in str(e):
                            logger.info(f"VM {dom.name()} já estava desligada, erro ignorado.")
                        else:
                            logger.error(f"Erro inesperado ao desligar VM {dom.name()}: {e}")
                else:
                    logger.info(f"VM {dom.name()} já está desligada, ignorando comando stop.")
            else:
                logger.warning(f"Ação desconhecida: {action}")

            # Remove a ação após executar
            redis_client.delete(key)
            logger.info(f"Ação {action} para VM {uuid} removida do Redis.")

    except Exception as e:
        logger.error(f"Erro ao buscar ações do Redis: {e}")

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

        logger.info("Verificando ações pendentes...")
        executar_acoes_pendentes(conn)

    except Exception as e:
        logger.error(f"Erro na tarefa principal: {e}")
    finally:
        if conn:
            conn.close()

if __name__ == "__main__":
    logger.info("Iniciando Bifrost Agent...")
    while True:
        tarefa_principal()
        time.sleep(60)  # inventário a cada 5 min

import libvirt
import json
import requests
import time
import logging
import os
import redis
import threading
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
INVENTORY_INTERVAL = int(os.getenv('BIFROST_INVENTORY_INTERVAL', '300'))  # padrão: 5 min

# Conexão Redis
redis_client = redis.Redis(host=REDIS_HOST, port=REDIS_PORT)
try:
    redis_client.ping()
    logger.info(f"✅ Conectado ao Redis em {REDIS_HOST}:{REDIS_PORT}")
except redis.exceptions.ConnectionError as e:
    logger.error(f"❌ Erro ao conectar ao Redis: {e}")
    exit(1)

def coletar_dados_vm(conn):
    vms = []
    for dom in conn.listAllDomains():
        try:
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

            # Interfaces protegido
            if dom.isActive() == 1:
                try:
                    iface_addrs = dom.interfaceAddresses(libvirt.VIR_DOMAIN_INTERFACE_ADDRESSES_SRC_LEASE, 0)
                    for iface in iface_addrs.values():
                        vm_info["interfaces"].append({
                            "mac": iface.get('hwaddr'),
                            "addrs": iface.get('addrs', [])
                        })
                except libvirt.libvirtError as e:
                    logger.warning(f"⚠️ Interfaces não coletadas de {dom.name()}: {e}")

            vms.append(vm_info)

        except Exception as e:
            logger.error(f"❌ Erro ao coletar dados da VM {dom.name()}: {e}")

    return vms

def enviar_dados_api(vms):
    payload = {
        "timestamp": datetime.now().astimezone().isoformat(),
        "vms": vms
    }
    try:
        resp = requests.post(API_URL, json=payload, timeout=10)
        if resp.status_code == 200:
            logger.info("✅ Inventário enviado com sucesso.")
        else:
            logger.error(f"❌ Falha ao enviar inventário: {resp.status_code} - {resp.text}")
    except Exception as e:
        logger.error(f"❌ Erro ao enviar inventário: {e}")

def inventario_loop():
    while True:
        try:
            conn = libvirt.open(None)
            if conn is None:
                logger.error("❌ Não foi possível conectar ao libvirt no inventário.")
            else:
                logger.info("🔄 Coletando inventário...")
                vms = coletar_dados_vm(conn)
                enviar_dados_api(vms)
        except Exception as e:
            logger.error(f"❌ Erro no loop de inventário: {e}")
        finally:
            if conn:
                conn.close()
        time.sleep(INVENTORY_INTERVAL)

def executar_acoes():
    try:
        conn = libvirt.open(None)
        if conn is None:
            logger.error("❌ Não foi possível conectar ao libvirt no worker de ações.")
            return

        pubsub = redis_client.pubsub()
        pubsub.subscribe(REDIS_CHANNEL)
        logger.info(f"🎧 Escutando canal Redis '{REDIS_CHANNEL}'...")

        for message in pubsub.listen():
            if message['type'] != 'message':
                continue

            try:
                action_data = json.loads(message['data'])
                uuid = action_data.get('uuid')
                action = action_data.get('action')

                if not uuid or not action:
                    logger.warning("⚠️ Mensagem incompleta recebida no Redis.")
                    continue

                try:
                    dom = conn.lookupByUUIDString(uuid)
                except libvirt.libvirtError:
                    logger.warning(f"⚠️ VM com UUID {uuid} não encontrada no libvirt.")
                    continue

                if action == "start":
                    if dom.isActive() == 0:
                        dom.create()
                        logger.info(f"✅ VM {dom.name()} iniciada com sucesso.")
                    else:
                        logger.info(f"ℹ️ VM {dom.name()} já está em execução.")

                elif action == "stop":
                    if dom.isActive() == 1:
                        try:
                            dom.shutdown()
                            logger.info(f"✅ VM {dom.name()} desligada com sucesso.")
                        except libvirt.libvirtError as e:
                            logger.warning(f"⚠️ Shutdown falhou para {dom.name()}, tentando destroy: {e}")
                            dom.destroy()
                            logger.info(f"✅ VM {dom.name()} forçada com destroy().")
                    else:
                        logger.info(f"ℹ️ VM {dom.name()} já está desligada.")

                else:
                    logger.warning(f"⚠️ Ação desconhecida recebida: {action}")

            except Exception as e:
                logger.error(f"❌ Erro ao processar mensagem Redis: {e}")

    except Exception as e:
        logger.error(f"❌ Erro no worker de ações: {e}")
    finally:
        if conn:
            conn.close()

if __name__ == "__main__":
    logger.info("🚀 Iniciando Bifrost Agent...")

    # Thread do inventário
    t1 = threading.Thread(target=inventario_loop, daemon=True)
    t1.start()

    # Listener Redis principal
    executar_acoes()

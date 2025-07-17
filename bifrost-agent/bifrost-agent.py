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
API_UPDATE_URL = os.getenv('BIFROST_API_UPDATE_URL', f"{API_URL}/update")
REDIS_HOST = os.getenv('REDIS_HOST', 'localhost')
REDIS_PORT = int(os.getenv('REDIS_PORT', '6379'))
REDIS_CHANNEL = os.getenv('REDIS_CHANNEL', 'vm-actions')
INVENTORY_INTERVAL = int(os.getenv('BIFROST_INVENTORY_INTERVAL', '300'))

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
        vm_info = build_vm_info(dom)
        if vm_info:
            vms.append(vm_info)
    return vms

def coletar_dados_vm_por_uuid(conn, uuid):
    try:
        dom = conn.lookupByUUIDString(uuid)
        vm_info = build_vm_info(dom)
        return vm_info
    except Exception as e:
        logger.error(f"❌ Erro ao coletar VM {uuid} após ação: {e}")
        return None

def build_vm_info(dom):
    try:
        name = dom.name()
        uuid = dom.UUIDString()
        try:
            state_code = dom.state()[0]
            state = {
                1: "running",
                3: "paused",
                5: "shut off",
                7: "crashed"
            }.get(state_code, f"unknown ({state_code})")
        except Exception as e:
            logger.warning(f"⚠️ {name}: falha em state() → {e}")
            state = "unknown"

        try:
            cpu = int(dom.maxVcpus())
        except Exception as e:
            logger.warning(f"⚠️ {name}: falha em maxVcpus() → {e}")
            cpu = 0

        try:
            memory = int(dom.maxMemory())
        except Exception as e:
            logger.warning(f"⚠️ {name}: falha em maxMemory() → {e}")
            memory = 0

        disks = []
        try:
            xml = dom.XMLDesc()
            for disk in xml.split('<disk type=')[1:]:
                if 'file=' in disk:
                    disk_path = disk.split('file=')[1].split("'")[1]
                    disks.append({"path": disk_path})
        except Exception as e:
            logger.warning(f"⚠️ {name}: falha ao coletar discos → {e}")

        interfaces = []
        try:
            if dom.isActive() == 1:
                iface_addrs = dom.interfaceAddresses(libvirt.VIR_DOMAIN_INTERFACE_ADDRESSES_SRC_LEASE, 0)
                for iface in iface_addrs.values():
                    interfaces.append({
                        "mac": iface.get('hwaddr', ''),
                        "addrs": [addr['addr'] for addr in iface.get('addrs', []) if 'addr' in addr]
                    })
        except Exception as e:
            logger.warning(f"⚠️ {name}: falha ao coletar interfaces → {e}")

        return {
            "name": name,
            "uuid": uuid,
            "state": state,
            "cpu_allocation": cpu,
            "memory_allocation": memory,
            "disks": disks,
            "interfaces": interfaces,
            "metadata": {},
        }
    except Exception as e:
        logger.error(f"❌ Erro geral ao coletar VM (skipada): {e}")
        return None

def enviar_dados_api(vms):
    payload = {
        "timestamp": datetime.now().astimezone().isoformat(),
        "vms": vms
    }
    try:
        resp = requests.post(API_URL, json=payload, timeout=15)
        if resp.status_code == 200:
            logger.info("✅ Inventário enviado com sucesso.")
        else:
            logger.error(f"❌ Falha ao enviar inventário: {resp.status_code} - {resp.text}")
    except Exception as e:
        logger.error(f"❌ Erro ao enviar inventário: {e}")

def report_action_to_api(uuid, action, result):
    payload = {
        "uuid": uuid,
        "action": action,
        "result": result,
        "timestamp": datetime.now().astimezone().isoformat()
    }
    try:
        resp = requests.post(API_UPDATE_URL, json=payload, timeout=10)
        if resp.status_code == 200:
            logger.info(f"✅ Atualização enviada para API após {action} em {uuid}.")
        else:
            logger.error(f"❌ Falha ao enviar atualização para API: {resp.status_code} - {resp.text}")
    except Exception as e:
        logger.error(f"❌ Erro ao enviar atualização para API: {e}")

def inventario_loop():
    while True:
        try:
            conn = libvirt.open(None)
            if conn is None:
                logger.error("❌ Não foi possível conectar ao libvirt no inventário.")
            else:
                logger.info("🔄 Coletando inventário completo...")
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
                        logger.info(f"✅ START executado na VM {dom.name()} ({uuid}).")
                        report_action_to_api(uuid, "start", "running")
                    else:
                        logger.info(f"ℹ️ VM {dom.name()} ({uuid}) já estava em execução.")
                        report_action_to_api(uuid, "start", "already_running")

                elif action == "stop":
                    if dom.isActive() == 1:
                        try:
                            dom.shutdown()
                            logger.info(f"✅ STOP executado na VM {dom.name()} ({uuid}).")
                            report_action_to_api(uuid, "stop", "shut off")
                        except libvirt.libvirtError as e:
                            logger.warning(f"⚠️ Shutdown falhou para {dom.name()}, tentando destroy: {e}")
                            dom.destroy()
                            logger.info(f"✅ STOP (destroy) forçado na VM {dom.name()} ({uuid}).")
                            report_action_to_api(uuid, "stop", "forced")
                    else:
                        logger.info(f"ℹ️ VM {dom.name()} ({uuid}) já estava desligada.")
                        report_action_to_api(uuid, "stop", "already_stopped")

                else:
                    logger.warning(f"⚠️ Ação desconhecida recebida: {action}")
                    continue

                # 🔥 Coleta e envia inventário imediato só dessa VM
                vm_info = coletar_dados_vm_por_uuid(conn, uuid)
                if vm_info:
                    enviar_dados_api([vm_info])

            except Exception as e:
                logger.error(f"❌ Erro ao processar mensagem Redis: {e}")

    except Exception as e:
        logger.error(f"❌ Erro no worker de ações: {e}")
    finally:
        if conn:
            conn.close()

if __name__ == "__main__":
    logger.info("🚀 Iniciando Bifrost Agent...")

    # Thread do inventário geral
    t1 = threading.Thread(target=inventario_loop, daemon=True)
    t1.start()

    # Listener Redis principal
    executar_acoes()

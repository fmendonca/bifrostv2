import libvirt
import json
import requests
import time
import logging
import os
import redis
import threading
from datetime import datetime

# Silenciar mensagens C do libvirt no stderr
libvirt.virEventRegisterDefaultImpl()
libvirt.virSetErrorFunc(None, lambda ctx, err: None)

# Configuração de log
LOG_FILE = '/var/log/bifrost/bifrost-agent.log'
os.makedirs('/var/log/bifrost', exist_ok=True)
logging.basicConfig(
    filename=LOG_FILE,
    level=logging.INFO,
    format='%(asctime)s [%(levelname)s] %(message)s'
)
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


def build_vm_info(dom):
    try:
        name = dom.name()
        uuid = dom.UUIDString()

        # Estado
        state_code = dom.state()[0]
        state = {
            1: "running",
            3: "paused",
            5: "shut off",
            7: "crashed"
        }.get(state_code, f"unknown ({state_code})")

        # CPU e Memória
        cpu = dom.maxVcpus() if dom.isActive() else 0
        memory = dom.maxMemory() if dom.isActive() else 0

        # Discos
        disks = []
        xml = dom.XMLDesc()
        for disk in xml.split('<disk type=')[1:]:
            if 'file=' in disk:
                disk_path = disk.split('file=')[1].split("'")[1]
                disks.append({"path": disk_path})

        # Interfaces
        interfaces = []
        if dom.isActive():
            iface_addrs = dom.interfaceAddresses(libvirt.VIR_DOMAIN_INTERFACE_ADDRESSES_SRC_LEASE, 0)
            for iface in iface_addrs.values():
                interfaces.append({
                    "mac": iface.get('hwaddr', ''),
                    "addrs": [addr['addr'] for addr in iface.get('addrs', []) if 'addr' in addr]
                })

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
        logger.debug(f"⚠️ Ignorado erro ao coletar VM: {e}")
        return None


def coletar_dados_vm(conn):
    return [info for dom in conn.listAllDomains() if (info := build_vm_info(dom))]


def enviar_dados_api(vms):
    payload = {"timestamp": datetime.now().astimezone().isoformat(), "vms": vms}
    try:
        resp = requests.post(API_URL, json=payload, timeout=15)
        if resp.status_code == 200:
            logger.info(f"✅ Inventário enviado: {len(vms)} VMs.")
        else:
            logger.error(f"❌ Erro ao enviar inventário: {resp.status_code} - {resp.text}")
    except Exception as e:
        logger.error(f"❌ Erro HTTP ao enviar inventário: {e}")


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
            logger.info(f"✅ Atualização para {uuid}: {action} → {result}")
        else:
            logger.error(f"❌ Falha ao atualizar {uuid}: {resp.status_code} - {resp.text}")
    except Exception as e:
        logger.error(f"❌ Erro HTTP na atualização {uuid}: {e}")


def inventario_loop():
    while True:
        try:
            conn = libvirt.open(None)
            if conn:
                vms = coletar_dados_vm(conn)
                enviar_dados_api(vms)
            else:
                logger.error("❌ Falha ao conectar no libvirt para inventário.")
        except Exception as e:
            logger.error(f"❌ Erro no inventário: {e}")
        finally:
            if conn:
                conn.close()
        time.sleep(INVENTORY_INTERVAL)


def executar_acoes():
    try:
        conn = libvirt.open(None)
        if not conn:
            logger.error("❌ Falha ao conectar no libvirt para ações.")
            return

        pubsub = redis_client.pubsub()
        pubsub.subscribe(REDIS_CHANNEL)
        logger.info(f"🎧 Escutando canal Redis '{REDIS_CHANNEL}'...")

        for message in pubsub.listen():
            if message['type'] != 'message':
                continue

            action_data = json.loads(message['data'])
            uuid = action_data.get('uuid')
            action = action_data.get('action')

            if not uuid or not action:
                logger.warning("⚠️ Mensagem incompleta no Redis.")
                continue

            try:
                dom = conn.lookupByUUIDString(uuid)
                name = dom.name()
            except libvirt.libvirtError:
                logger.warning(f"⚠️ VM UUID {uuid} não encontrada.")
                continue

            result = "unknown"
            if action == "start":
                if not dom.isActive():
                    dom.create()
                    result = "running"
                    logger.info(f"✅ START → {name} ({uuid})")
                else:
                    result = "already_running"
            elif action == "stop":
                if dom.isActive():
                    try:
                        dom.shutdown()
                        result = "shut off"
                        logger.info(f"✅ STOP → {name} ({uuid})")
                    except libvirt.libvirtError:
                        dom.destroy()
                        result = "forced"
                        logger.info(f"✅ FORCED STOP (destroy) → {name} ({uuid})")
                else:
                    result = "already_stopped"
            else:
                logger.warning(f"⚠️ Ação desconhecida: {action}")
                continue

            report_action_to_api(uuid, action, result)

            # Atualiza API só com a VM afetada
            vm_info = build_vm_info(dom)
            if vm_info:
                enviar_dados_api([vm_info])

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

    # Listener Redis principal (fica no foreground)
    executar_acoes()

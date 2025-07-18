import libvirt
import json
import requests
import time
import logging
import os
import redis
import threading
from datetime import datetime

# Config log
LOG_FILE = '/var/log/bifrost/bifrost-agent.log'
os.makedirs('/var/log/bifrost', exist_ok=True)
logging.basicConfig(filename=LOG_FILE, level=logging.INFO,
                    format='%(asctime)s [%(levelname)s] %(message)s')
logger = logging.getLogger()

# Env vars (mandatory)
API_URL = os.getenv('BIFROST_API_URL', 'http://localhost:8080')
AGENT_NAME = os.getenv('BIFROST_AGENT_NAME')
API_KEY = os.getenv('BIFROST_API_KEY')

if not AGENT_NAME or not API_KEY:
    logger.error("❌ BIFROST_AGENT_NAME and BIFROST_API_KEY must be set in environment")
    exit(1)

# Redis setup
REDIS_CHANNEL = f"vm-actions-{AGENT_NAME}"
redis_client = redis.Redis(host=os.getenv('REDIS_HOST', 'localhost'), port=int(os.getenv('REDIS_PORT', '6379')))

def coletar_dados_vm(conn):
    vms = []
    for dom in conn.listAllDomains():
        try:
            state_code = dom.state()[0]
            state = {1: "running", 3: "paused", 5: "shut off", 7: "crashed"}.get(state_code, "unknown")
            disks = []
            try:
                xml = dom.XMLDesc()
                for disk in xml.split('<disk type=')[1:]:
                    if 'file=' in disk:
                        disks.append({"path": disk.split('file=')[1].split("'")[1]})
            except:
                pass
            interfaces = []
            try:
                if dom.isActive():
                    iface_addrs = dom.interfaceAddresses(libvirt.VIR_DOMAIN_INTERFACE_ADDRESSES_SRC_LEASE, 0)
                    for iface in iface_addrs.values():
                        interfaces.append({
                            "mac": iface.get('hwaddr', ''),
                            "addrs": [addr['addr'] for addr in iface.get('addrs', []) if 'addr' in addr]
                        })
            except:
                pass
            vms.append({
                "name": dom.name(),
                "uuid": dom.UUIDString(),
                "state": state,
                "cpu_allocation": dom.maxVcpus() if dom.isActive() else 0,
                "memory_allocation": dom.maxMemory() if dom.isActive() else 0,
                "disks": disks,
                "interfaces": interfaces,
                "metadata": {},
                "host_uuid": AGENT_NAME  # using name as unique id now
            })
        except Exception as e:
            logger.warning(f"⚠️ Failed to collect VM info: {e}")
    return vms

def enviar_dados_api(vms):
    payload = {"timestamp": datetime.now().astimezone().isoformat(), "vms": vms}
    try:
        res = requests.post(f"{API_URL}/api/v1/vms", json=payload, headers={"X-API-KEY": API_KEY}, timeout=15)
        res.raise_for_status()
        logger.info(f"✅ Sent {len(vms)} VMs to API")
    except Exception as e:
        logger.error(f"❌ Failed to send inventory: {e}")

def report_action(uuid, action, result):
    payload = {"uuid": uuid, "action": action, "result": result, "timestamp": datetime.now().astimezone().isoformat()}
    try:
        res = requests.post(f"{API_URL}/api/v1/vms/update", json=payload, headers={"X-API-KEY": API_KEY}, timeout=10)
        res.raise_for_status()
        logger.info(f"✅ Reported {action} → {result} for VM {uuid}")
    except Exception as e:
        logger.error(f"❌ Failed to report action: {e}")

def inventario_loop():
    while True:
        try:
            conn = libvirt.open(None)
            if conn:
                vms = coletar_dados_vm(conn)
                enviar_dados_api(vms)
                conn.close()
            else:
                logger.error("❌ Failed to connect libvirt")
        except Exception as e:
            logger.error(f"❌ Inventory loop error: {e}")
        time.sleep(300)

def executar_acoes():
    pubsub = redis_client.pubsub()
    pubsub.subscribe(REDIS_CHANNEL)
    logger.info(f"🎧 Listening Redis channel {REDIS_CHANNEL}")
    conn = libvirt.open(None)
    if not conn:
        logger.error("❌ Failed to connect libvirt for actions")
        return

    for message in pubsub.listen():
        if message['type'] != 'message':
            continue
        try:
            data = json.loads(message['data'])
            uuid, action = data.get('uuid'), data.get('action')
            dom = conn.lookupByUUIDString(uuid)
            if action == "start" and not dom.isActive():
                dom.create()
                logger.info(f"✅ Started VM {dom.name()}")
                report_action(uuid, action, "running")
            elif action == "stop" and dom.isActive():
                dom.shutdown()
                logger.info(f"✅ Stopped VM {dom.name()}")
                report_action(uuid, action, "shut off")
            else:
                logger.info(f"ℹ️ VM {dom.name()} already {dom.state()[0]}")
        except Exception as e:
            logger.error(f"❌ Action error: {e}")

if __name__ == "__main__":
    logger.info(f"🚀 Starting Bifrost Agent as {AGENT_NAME}")
    threading.Thread(target=inventario_loop, daemon=True).start()
    executar_acoes()

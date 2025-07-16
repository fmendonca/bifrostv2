import libvirt
import json
import requests
import time
import logging
import os
from datetime import datetime

# Configuração de log
LOG_FILE = '/var/log/bifrost/bifrost-agent.log'
os.makedirs('/var/log/bifrost', exist_ok=True)
logging.basicConfig(filename=LOG_FILE, level=logging.INFO,
                    format='%(asctime)s [%(levelname)s] %(message)s')
logger = logging.getLogger()

# Variáveis de ambiente
API_URL = os.getenv('BIFROST_API_URL', 'http://localhost:8080/api/v1/vms')
INFINISPAN_URL = os.getenv('INFINISPAN_URL', 'http://localhost:11222/rest/v2/caches/vm-actions')
INFINISPAN_USER = os.getenv('INFINISPAN_USER', 'user')
INFINISPAN_PASS = os.getenv('INFINISPAN_PASS', 'pass')

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

        # Interfaces
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
            logger.info(f"Inventário enviado com sucesso.")
        else:
            logger.error(f"Falha ao enviar inventário: {resp.status_code}")
    except Exception as e:
        logger.error(f"Erro ao enviar inventário: {e}")

def executar_acoes_pendentes(conn):
    try:
        resp = requests.get(INFINISPAN_URL, auth=(INFINISPAN_USER, INFINISPAN_PASS))
        if resp.status_code != 200 or not resp.text:
            return  # nada para fazer

        actions = json.loads(resp.text)
        uuid = actions.get('uuid')
        action = actions.get('action')

        if not uuid or not action:
            return

        dom = conn.lookupByUUIDString(uuid)
        if action == "start":
            if dom.isActive() == 0:
                dom.create()
                logger.info(f"VM {dom.name()} iniciada com sucesso.")
            else:
                logger.info(f"VM {dom.name()} já está em execução.")
        elif action == "stop":
            if dom.isActive() == 1:
                dom.shutdown()
                logger.info(f"VM {dom.name()} desligada com sucesso.")
            else:
                logger.info(f"VM {dom.name()} já está desligada.")
        else:
            logger.warning(f"Ação desconhecida: {action}")

        # Limpa pending_action no backend
        try:
            clear_url = f"{API_URL}/{uuid}/{action}"
            requests.post(clear_url)
            logger.info(f"Ação {action} confirmada no backend para VM {uuid}.")
        except Exception as e:
            logger.error(f"Erro ao notificar backend: {e}")

    except Exception as e:
        logger.error(f"Erro ao buscar ações do Infinispan: {e}")

def tarefa_principal():
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
        time.sleep(300)  # inventário a cada 5 min

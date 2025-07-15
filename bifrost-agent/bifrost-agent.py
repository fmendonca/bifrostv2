import libvirt
import json
import requests
import sched
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

# Configuração do scheduler
scheduler = sched.scheduler(time.time, time.sleep)

# Endpoint da API (lido de variável ambiente)
API_URL = os.getenv('BIFROST_API_URL', 'http://localhost:8080/api/v1/vms')

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
            disk_path = disk.split('file=')[1].split('\'')[1]
            vm_info["disks"].append({"path": disk_path})

        # Interfaces
        for iface in dom.interfaceAddresses(libvirt.VIR_DOMAIN_INTERFACE_ADDRESSES_SRC_LEASE, 0).values():
            vm_info["interfaces"].append({
                "mac": iface.get('hwaddr'),
                "addrs": iface.get('addrs', [])
            })

        # Metadata
        vm_info["metadata"] = dom.metadata(libvirt.VIR_DOMAIN_METADATA_ELEMENT, None, 0)

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
        resp = requests.get(f"{API_URL}?pending_action=1")
        resp.raise_for_status()
        vms_pendentes = resp.json()
    except Exception as e:
        logger.error(f"Erro ao buscar ações pendentes: {e}")
        return

    for vm in vms_pendentes:
        try:
            dom = conn.lookupByUUIDString(vm['uuid'])
            action = vm.get('pending_action')
            if action == "start":
                if dom.isActive() == 0:
                    dom.create()
                    logger.info(f"VM {vm['name']} iniciada com sucesso.")
                else:
                    logger.info(f"VM {vm['name']} já está em execução.")
                requests.post(f"{API_URL}/{vm['uuid']}/start")
            elif action == "stop":
                if dom.isActive() == 1:
                    dom.shutdown()
                    logger.info(f"VM {vm['name']} desligada com sucesso.")
                else:
                    logger.info(f"VM {vm['name']} já está desligada.")
                requests.post(f"{API_URL}/{vm['uuid']}/stop")
        except Exception as e:
            logger.error(f"Erro ao processar ação '{action}' para VM {vm['name']}: {e}")

def tarefa_principal():
    try:
        conn = libvirt.open(None)
        if conn is None:
            logger.error("Não foi possível conectar ao libvirt.")
            return

        logger.info("Coletando dados das VMs...")
        vms = coletar_dados_vm(conn)
        enviar_dados_api(vms)
        executar_acoes_pendentes(conn)

    except Exception as e:
        logger.error(f"Erro na tarefa principal: {e}")
    finally:
        if conn:
            conn.close()

def agendar_execucao():
    def run_periodically():
        tarefa_principal()
        scheduler.enter(300, 1, run_periodically)

    scheduler.enter(0, 1, run_periodically)
    scheduler.run()

if __name__ == "__main__":
    logger.info("Iniciando Bifrost Agent...")
    agendar_execucao()

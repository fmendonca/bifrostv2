#!/usr/bin/env python3

import libvirt
import json
import time
import requests
import logging
import sys
import os
from datetime import datetime, UTC
import sched
import xml.etree.ElementTree as ET

API_ENDPOINT = os.getenv("API_ENDPOINT", "http://localhost:8080/api/v1/vms")
LIBVIRT_URI = "qemu:///system"
INTERVALO_MINUTOS = 5
LOG_DIR = "/var/log/bifrost"
LOG_FILE = os.path.join(LOG_DIR, "bifrost-agent.log")

# Garante que o diretório de log existe
os.makedirs(LOG_DIR, exist_ok=True)

# Setup de logging: stdout + arquivo
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(message)s",
    handlers=[
        logging.StreamHandler(sys.stdout),
        logging.FileHandler(LOG_FILE)
    ]
)

def coletar_dados_vm(conn):
    dados_vms = []
    try:
        for id in conn.listDomainsID():
            dom = conn.lookupByID(id)
            dados_vms.append(extrair_info_vm(dom))
        for nome in conn.listDefinedDomains():
            dom = conn.lookupByName(nome)
            dados_vms.append(extrair_info_vm(dom))
    except libvirt.libvirtError as e:
        logging.error(f"Erro ao coletar dados das VMs: {e}")
    return dados_vms

def extrair_info_vm(dom):
    try:
        info = dom.info()
        estado_map = {
            libvirt.VIR_DOMAIN_NOSTATE: "no state",
            libvirt.VIR_DOMAIN_RUNNING: "running",
            libvirt.VIR_DOMAIN_BLOCKED: "blocked",
            libvirt.VIR_DOMAIN_PAUSED: "paused",
            libvirt.VIR_DOMAIN_SHUTDOWN: "shutdown",
            libvirt.VIR_DOMAIN_SHUTOFF: "shutoff",
            libvirt.VIR_DOMAIN_CRASHED: "crashed",
            libvirt.VIR_DOMAIN_PMSUSPENDED: "pmsuspended",
        }

        estado = estado_map.get(info[0], "unknown")
        cpu = info[3]
        memoria = info[1]

        # Parse XML for disk paths
        xml_desc = dom.XMLDesc()
        root = ET.fromstring(xml_desc)
        discos = []
        for disk in root.findall("./devices/disk"):
            device = disk.find("target").get("dev")
            source = disk.find("source")
            path = source.get("file") if source is not None else None
            if path:
                discos.append({"device": device, "path": path})

        interfaces = []
        try:
            ifaces = dom.interfaceAddresses(libvirt.VIR_DOMAIN_INTERFACE_ADDRESSES_SRC_AGENT, 0)
            for iface_name, iface_data in ifaces.items():
                interfaces.append({
                    "name": iface_name,
                    "mac": iface_data.get('hwaddr'),
                    "addrs": [addr['addr'] for addr in iface_data.get('addrs') or []]
                })
        except libvirt.libvirtError:
            pass

        try:
            metadados = dom.metadata(libvirt.VIR_DOMAIN_METADATA_ELEMENT, None, 0)
        except libvirt.libvirtError:
            metadados = {}

        return {
            "name": dom.name(),
            "uuid": dom.UUIDString(),
            "state": estado,
            "cpu_allocation": cpu,
            "memory_allocation": memoria,
            "disks": discos,
            "interfaces": interfaces,
            "metadata": metadados
        }

    except libvirt.libvirtError as e:
        logging.error(f"Erro ao extrair info da VM {dom.name()}: {e}")
        return {"name": dom.name(), "error": str(e)}

def enviar_para_api(payload):
    try:
        headers = {'Content-Type': 'application/json'}
        response = requests.post(API_ENDPOINT, json=payload, headers=headers)
        if response.ok:
            logging.info(f"Dados enviados com sucesso para {API_ENDPOINT}. Código {response.status_code}")
        else:
            logging.error(f"Falha ao enviar dados para {API_ENDPOINT}. Código {response.status_code} Resposta: {response.text}")
    except requests.RequestException as e:
        logging.error(f"Erro de comunicação com API {API_ENDPOINT}: {e}")

def tarefa_principal():
    logging.info(f"Iniciando coleta de dados das VMs... Endpoint configurado: {API_ENDPOINT}")
    try:
        conn = libvirt.open(LIBVIRT_URI)
        if conn is None:
            logging.error("Falha ao conectar ao Libvirt.")
            return
        dados_vms = coletar_dados_vm(conn)
        conn.close()
        json_payload = {"timestamp": datetime.now(UTC).isoformat(), "vms": dados_vms}
        logging.info("Coleta concluída, enviando para API...")
        enviar_para_api(json_payload)
    except libvirt.libvirtError as e:
        logging.error(f"Erro geral do Libvirt: {e}")

def agendar_execucao():
    scheduler = sched.scheduler(time.time, time.sleep)
    def run_periodically():
        tarefa_principal()
        scheduler.enter(INTERVALO_MINUTOS * 60, 1, run_periodically)
    scheduler.enter(0, 1, run_periodically)
    scheduler.run()

if __name__ == "__main__":
    logging.info("Iniciando Bifrost Agent...")
    agendar_execucao()

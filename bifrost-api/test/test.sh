#!/bin/bash

API_URL="http://192.168.86.129:8080/api/v1/vms"

echo "📤 Enviando VMs de teste para $API_URL ..."

curl -s -X POST "$API_URL" \
-H "Content-Type: application/json" \
-d '{
  "timestamp": "2025-07-12T22:15:00Z",
  "vms": [
    {
      "name": "test-vm-1",
      "uuid": "11111111-1111-1111-1111-111111111111",
      "state": "running",
      "cpu_allocation": 2,
      "memory_allocation": 4096,
      "disks": [{"device": "vda"}],
      "interfaces": [{"name": "eth0", "mac": "52:54:00:12:34:56", "addrs": ["192.168.122.101"]}],
      "metadata": {"owner": "test-user"}
    },
    {
      "name": "test-vm-2",
      "uuid": "22222222-2222-2222-2222-222222222222",
      "state": "shutoff",
      "cpu_allocation": 4,
      "memory_allocation": 8192,
      "disks": [{"device": "vdb"}],
      "interfaces": [{"name": "eth1", "mac": "52:54:00:65:43:21", "addrs": ["192.168.122.102"]}],
      "metadata": {"owner": "admin"}
    }
  ]
}'

echo -e "\n✅ VMs de teste enviadas com sucesso."

echo "📥 Buscando VMs armazenadas no banco..."

curl -s -X GET "$API_URL" | jq .

echo -e "\n✅ Listagem concluída."

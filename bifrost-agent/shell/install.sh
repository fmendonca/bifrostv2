#!/bin/bash

set -e

# Configurações
AGENT_DIR="/opt/bifrost-agent"
AGENT_SCRIPT="bifrost-agent.py"
SERVICE_FILE="/etc/systemd/system/bifrost-agent.service"
LOG_DIR="/var/log/bifrost"
LOGROTATE_FILE="/etc/logrotate.d/bifrost-agent"

echo "🔧 Instalando Bifrost Agent..."

# 1️⃣ Criar diretório do agente
echo "📁 Criando diretório $AGENT_DIR"
sudo mkdir -p $AGENT_DIR
sudo cp $AGENT_SCRIPT $AGENT_DIR/
sudo chmod +x $AGENT_DIR/$AGENT_SCRIPT

# 2️⃣ Criar diretório de log
echo "📁 Criando diretório $LOG_DIR"
sudo mkdir -p $LOG_DIR
sudo chown $(whoami):$(whoami) $LOG_DIR

# 3️⃣ Instalar systemd service
echo "📝 Instalando systemd service em $SERVICE_FILE"
sudo bash -c "cat > $SERVICE_FILE" <<EOF
[Unit]
Description=Bifrost Agent Python Service
After=network.target

[Service]
Type=simple
WorkingDirectory=$AGENT_DIR
ExecStart=/usr/bin/python3 $AGENT_DIR/$AGENT_SCRIPT
Restart=always
RestartSec=10
Environment="API_ENDPOINT=http://localhost:8080/api/v1/vms"

[Install]
WantedBy=multi-user.target
EOF

# 4️⃣ (Opcional) Criar logrotate config
echo "📝 Instalando logrotate config em $LOGROTATE_FILE"
sudo bash -c "cat > $LOGROTATE_FILE" <<EOF
$LOG_DIR/*.log {
    daily
    rotate 7
    compress
    missingok
    notifempty
    copytruncate
}
EOF

# 5️⃣ Ativar e iniciar serviço
echo "🚀 Habilitando e iniciando serviço"
sudo systemctl daemon-reload
sudo systemctl enable --now bifrost-agent

echo "✅ Bifrost Agent instalado e rodando!"
echo "👉 Status: sudo systemctl status bifrost-agent"
echo "👉 Logs: sudo journalctl -u bifrost-agent -f"
echo "👉 Arquivo de log: $LOG_DIR/bifrost-agent.log"

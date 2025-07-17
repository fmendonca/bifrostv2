
# 📦 Bifrost v2

Orquestrador de hosts de virtualização com KVM.

## ✨ Visão Geral

O **Bifrost v2** é uma plataforma para gerenciamento de hosts de virtualização baseados em KVM, oferecendo:

✅ Cadastro e gerenciamento de hosts físicos  
✅ Listagem e controle de máquinas virtuais (VMs)  
✅ Interface web moderna para administração  
✅ Backend robusto em Golang  
✅ Frontend responsivo em Node.js (React)  
✅ Agentes Python para coleta e execução remota

## 🏗️ Arquitetura

- **Frontend (NodeJS / React):** Interface web para interação com os usuários e administradores.  
- **Backend (Golang):** API REST, controle de usuários, RBAC, orquestração de ações.  
- **Agent (Python):** Coletor e executor local em cada host, usando libvirt.

## 🚀 Instalação

### Pré-requisitos

- Docker ou Podman  
- Redis  
- PostgreSQL  
- Go ≥ 1.22  
- Node.js ≥ 18  
- Python ≥ 3.11

### Passos

1️⃣ Clone o repositório:
```bash
git clone https://seurepositorio/bifrost.git
cd bifrost
```

2️⃣ Build do backend:
```bash
cd backend
go build -o bifrost
```

3️⃣ Build do frontend:
```bash
cd frontend
npm install
REACT_APP_API_URL=http://localhost:8080 npm run build
```

4️⃣ Configure os agentes:
```bash
cd agent
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
```

5️⃣ Suba os containers (opcional):
```bash
docker compose up -d
```

## ⚙️ Configuração

- `.env` para variáveis do backend  
- `.env.local` no frontend para URL da API  
- `config.yaml` no agent para IP, credenciais e endpoints

## 📡 Conectividade

Os agentes Python se conectam ao backend via fila (Redis) e reportam status/execuções.  
O backend fornece API REST consumida pelo frontend.

## 🛡️ Roadmap

- [ ] Dashboard com métricas em tempo real  
- [ ] Migração automatizada de VMs  
- [ ] Integração com storage externo  
- [ ] Autenticação via LDAP/OAuth

## 👥 Contribuição

Pull requests são bem-vindos!  
Veja [CONTRIBUTING.md](CONTRIBUTING.md) para mais informações.

## 📝 Licença

Este projeto está licenciado sob a licença MIT.

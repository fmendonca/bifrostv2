  

# 📦 Bifrost v2

  

Orquestrador de hosts de virtualização com KVM.

  

## ✨ Visão Geral
O **Bifrost** é uma plataforma para gerenciamento de hosts de virtualização baseados em KVM, oferecendo:

  

✅ Cadastro e gerenciamento de hosts físicos
✅ Listagem e controle de máquinas virtuais (VMs)
✅ Interface web moderna para administração
✅ Backend robusto em Golang
✅ Frontend responsivo em Node.js (React)
✅ Agentes Python para coleta e execução remota
  

## 🏗️ Arquitetura
-  **Frontend (NodeJS / React):** Interface web para interação com os usuários e administradores.
-  **Backend (Golang):** API REST, controle de usuários, RBAC, orquestração de ações.
-  **Agent (Python):** Coletor e executor local em cada host, usando libvirt.
 

## 🚀 Instalação
### Pré-requisitos
  

- Docker ou Podman
- Redis
- PostgreSQL
- Go ≥ 1.22
- Node.js ≥ 18
- Python ≥ 3.11
 

### Passos

em Desenvolvimento
**O requisito é que tenha: 
PostgreSQL 15+
REDIS 7+**

Por enquanto 

rode o ./[build-containers.sh](https://github.com/fmendonca/bifrostv2/blob/main/build-containers.sh "build-containers.sh") para build dos containers de backend e front. 
rode depois o ./[run-containers.sh](https://github.com/fmendonca/bifrostv2/blob/main/run-containers.sh "run-containers.sh") para que os containers de backend e front subam.


## ⚙️ Configuração
 
em Desenvolvimento
  

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
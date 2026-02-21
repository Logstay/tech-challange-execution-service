<h1 align="center">Tech Challenge - Execution Service</h1>

<p align="center">
  Microsserviço de <strong>Catálogo e Estoque</strong> do Tech Challenge - Fase 4 do 12SOAT FIAP.
  Responsável pelo gerenciamento de <strong>Consumíveis</strong>, <strong>Definições de Serviço</strong> e <strong>Controle de Inventário</strong>.
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.24-blue?logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/PostgreSQL-15-blue?logo=postgresql" alt="PostgreSQL">
  <img src="https://img.shields.io/badge/Kafka-3.6-black?logo=apachekafka" alt="Kafka">
  <img src="https://img.shields.io/badge/Kubernetes-1.28-blue?logo=kubernetes" alt="Kubernetes">
  <img src="https://img.shields.io/badge/Architecture-Hexagonal-orange" alt="Hexagonal Architecture">
</p>

---

## 📖 Sobre o Projeto

O **Execution Service** é o microsserviço responsável pelos contextos de **Catálogo** e **Execução** na arquitetura de microserviços do Tech Challenge. Ele gerencia:

- **Consumíveis (Peças)**: CRUD completo de peças e itens de estoque
- **Definições de Serviço**: Catálogo de serviços disponíveis com precificação
- **Controle de Inventário**: Baixa automática de estoque via eventos Kafka

### 🎯 Bounded Contexts

Este serviço implementa dois Bounded Contexts do DDD:

1. **Catalog Context**: Catálogo de serviços e consumíveis
   - ServiceDefinition (Definição de Serviço)
   - Consumable (Consumível/Peça)

2. **Inventory Context**: Controle de estoque
   - Inventory (Estoque)
   - Stock Movements (Movimentações)

### 🔗 Comunicação com Outros Serviços

```
┌─────────────────┐     Kafka      ┌───────────────────┐
│   OS Service    │ ───────────────► │ Execution Service │
│                 │  inventory-    │                   │
│                 │  decreased     │                   │
└─────────────────┘                └───────────────────┘
                                          │
                                          │ PostgreSQL
                                          ▼
                                   ┌─────────────────┐
                                   │  execution_db   │
                                   └─────────────────┘
```

---

## 🏗️ Arquitetura

### Visão Geral

O projeto implementa **Clean Architecture** com **Arquitetura Hexagonal** e **DDD**.

### Estrutura do Projeto

```
📁 tech-challange-execution-service/
├── cmd/                        # Ponto de entrada da aplicação
│   ├── main.go                 # Bootstrap da aplicação
│   └── initialize/             # Inicialização de dependências
├── config/                     # Configurações da aplicação
│   └── local/                  # Docker Compose local
├── internal/                   # Código de negócio (não exportado)
│   ├── adapters/               # Adaptadores (entrada/saída)
│   │   ├── http/               # Controllers REST (Gin)
│   │   ├── kafka/              # Consumer/Producer Kafka
│   │   └── postgres/           # Persistência PostgreSQL
│   ├── application/            # Casos de uso (Services)
│   │   ├── consumable/         # CRUD de consumíveis
│   │   │   ├── commands/       # Create, Update, Delete
│   │   │   └── queries/        # Get, List
│   │   └── servicedefinition/  # CRUD de definições
│   │       ├── commands/       # Create, Update, Delete
│   │       └── queries/        # Get, List
│   ├── domain/                 # Entidades e regras de negócio
│   │   ├── consumable/         # Agregado Consumable
│   │   ├── servicedefinition/  # Agregado ServiceDefinition
│   │   └── inventory/          # Agregado Inventory
│   │       └── events/         # Domain Events
│   ├── interfaces/             # Interfaces de serviços
│   └── ports/                  # Interfaces/Contratos
│       ├── consumable/         # Ports de consumíveis
│       ├── servicedefinition/  # Ports de serviços
│       ├── inventory/          # Ports de inventário
│       └── message-broker/     # Ports de mensageria
├── pkg/                        # Utilitários compartilhados
│   ├── logger/                 # Logging estruturado
│   ├── tracing/                # OpenTelemetry
│   └── null/                   # Tipos nullable
├── docs/                       # Documentação Swagger
├── test/                       # Testes de integração
├── k8s/                        # Manifests Kubernetes
│   └── helm/                   # Helm chart
└── Dockerfile                  # Container da aplicação
```

---

## 🛠️ Stack Tecnológica

| Tecnologia | Versão | Uso |
|------------|--------|-----|
| **Go** | 1.24 | Linguagem principal |
| **PostgreSQL** | 15 | Banco de dados relacional |
| **Kafka** | 3.6 | Mensageria assíncrona |
| **Gin** | 1.10 | Framework HTTP |
| **OpenTelemetry** | - | Tracing distribuído |
| **Docker** | 24+ | Containerização |
| **Kubernetes** | 1.28+ | Orquestração |
| **Helm** | 3+ | Deploy K8s |

---

## 🚀 Instruções de Execução

### Pré-requisitos

- Docker 24+ e Docker Compose
- Go 1.24+ (para desenvolvimento)

### Variáveis de Ambiente

```bash
# PostgreSQL
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=execution_db
POSTGRES_USER=execution
POSTGRES_PASSWORD=execution123

# Kafka
KAFKA_PRODUCER_HOSTS=localhost:9092
KAFKA_CONSUMER_HOSTS=localhost:9092

# Tracing
OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317

# Server
SERVER_PORT=8082
```

### Execução Local

```bash
# Subir infraestrutura (PostgreSQL, Kafka, etc)
cd config/local && docker compose up -d

# Executar aplicação
go run cmd/main.go

# Ou via Makefile
make dev-run
```

### Deploy Kubernetes

```bash
# Via Helm
helm install execution-service ./k8s/helm -n app

# Verificar pods
kubectl get pods -n app -l app=execution-service
```

---

## 📚 API Documentation

### Endpoints de Consumíveis

| Método | Endpoint | Descrição | Auth |
|--------|----------|-----------|------|
| `GET` | `/v1/consumables` | Listar consumíveis (paginado) | ✅ |
| `GET` | `/v1/consumables/{uid}` | Buscar por UID | ✅ |
| `POST` | `/v1/consumables` | Criar consumível | ✅ |
| `PUT` | `/v1/consumables/{uid}` | Atualizar consumível | ✅ |
| `DELETE` | `/v1/consumables/{uid}` | Deletar consumível | ✅ |

### Endpoints de Definições de Serviço

| Método | Endpoint | Descrição | Auth |
|--------|----------|-----------|------|
| `GET` | `/v1/service-definitions` | Listar serviços (paginado) | ✅ |
| `GET` | `/v1/service-definitions/{uid}` | Buscar por UID | ✅ |
| `POST` | `/v1/service-definitions` | Criar serviço | ✅ |
| `PUT` | `/v1/service-definitions/{uid}` | Atualizar serviço | ✅ |
| `DELETE` | `/v1/service-definitions/{uid}` | Deletar serviço | ✅ |

### Health Check

| Método | Endpoint | Descrição | Auth |
|--------|----------|-----------|------|
| `GET` | `/health` | Health check | ❌ |

### Swagger

- **Local**: http://localhost:8082/docs/index.html
- **Regenerar**: `make swag`

---

## 📨 Eventos Kafka

### Eventos Consumidos

| Tópico | Evento | Descrição |
|--------|--------|-----------|
| `inventory-decreased` | DecreaseInventory | Baixa de estoque |
| `inventory-decreased-retry` | RetryDecrease | Retry de baixa |

### Eventos Publicados

| Tópico | Evento | Descrição |
|--------|--------|-----------|
| `inventory-decreased-failed` | DecreaseFailed | Falha na baixa de estoque |

### Payload de Baixa de Estoque

```json
{
  "consumable_uid": "uuid-do-consumivel",
  "quantity": 2,
  "work_order_id": "uuid-da-os",
  "reason": "Uso em ordem de serviço"
}
```

---

## 🧪 Testes e Qualidade

```bash
# Testes unitários
make unit_test

# Testes de integração
make it_test

# Análise de lint
make check

# Cobertura
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

**Meta de Cobertura**: Mínimo 80%

---

## 🗄️ Banco de Dados (PostgreSQL)

### Estrutura de Tabelas

```sql
-- Consumíveis
CREATE TABLE consumable (
    uid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Definições de Serviço
CREATE TABLE service_definition (
    uid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Inventário
CREATE TABLE inventory (
    uid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    consumable_uid UUID REFERENCES consumable(uid),
    quantity INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

---

## 📊 Observabilidade

- **Tracing**: OpenTelemetry + Jaeger/New Relic
- **Logs**: Logrus estruturado
- **Métricas**: Prometheus-ready

---

## 📜 Scripts Makefile

| Comando | Descrição |
|---------|-----------|
| `make up` | Sobe infraestrutura Docker |
| `make down` | Para containers |
| `make dev-run` | Executa aplicação local |
| `make check` | Análise de lint |
| `make swag` | Regenera Swagger |
| `make unit_test` | Executa testes unitários |
| `make it_test` | Executa testes de integração |

---

<div align="center">

**Tech Challenge - Fase 4 - 12SOAT FIAP**

[![Go](https://img.shields.io/badge/Made%20with-Go-blue?logo=go)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/Database-PostgreSQL-blue?logo=postgresql)](https://postgresql.org/)
[![Kafka](https://img.shields.io/badge/Messaging-Kafka-black?logo=apachekafka)](https://kafka.apache.org/)

</div>
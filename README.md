# Mão na Roda - Backend

API RESTful para marketplace de serviços técnicos, conectando profissionais a clientes com geolocalização e sistema de avaliações. Desenvolvido seguindo **Clean Architecture** em **Go**.

## 🚀 Tecnologias

- **Linguagem**: Go 1.25+
- **Framework Web**: Gin Gonic
- **Banco de Dados**: PostgreSQL + PostGIS (via Supabase)
- **Driver DB**: pgx/v5 (Pool de conexões)
- **Autenticação**: Supabase Auth (Validação JWT)

## 🏗 Arquitetura

O projeto segue a Clean Architecture com a seguinte estrutura:

- `cmd/api`: Ponto de entrada (Main e Wiring).
- `internal/domain`: Entidades core e interfaces (Regras de negócio puras).
- `internal/service`: Casos de uso e lógica de aplicação.
- `internal/repository`: Implementação de persistência (SQL/PostGIS).
- `internal/handler`: Controladores HTTP (Gin).
- `pkg/middleware`: Middlewares reutilizáveis (Auth, Logger).
- `pkg/database`: Setup de conexão.

## 🛠 Como Rodar

### Pré-requisitos
- Go 1.25+ instalado.
- Docker (opcional, para rodar banco localmente).
- Conta no Supabase (ou Postgres local com extensão PostGIS habilitada).

### Configuração

1. Clone o repositório.
2. Crie um arquivo `.env` na raiz com a connection string do banco:
   ```env
   DATABASE_URL=postgres://user:pass@host:5432/dbname?sslmode=disable
   PORT=8080
   ```
   *Nota: Certifique-se que o banco tenha a extensão PostGIS habilitada: `CREATE EXTENSION postgis;`*

3. Execute as migrações (manualmente no banco por enquanto, vide schema sugerido no `implementation_plan.md`).

### Executando

```bash
go run cmd/api/main.go
```

## 🧪 Testando Endpoints (via Postman/Insomnia)

### Autenticação
1. **Registrar Usuário**: `POST /register`
   - Body: `{ "name": "João", "email": "joao@email.com", "password": "123", "type": "client" }`
2. **Login**: `POST /login`
   - Body: `{ "email": "joao@email.com", "password": "123" }`
   - *Retorna um Token Mockado para desenvolvimento*

### Profissionais
1. **Cadastrar Perfil Profissional** (Requer Auth Header `Authorization: Bearer <token>`):
   - `POST /profissionais`
   - Header: `X-User-ID: <uuid-do-usuario>` (Mock manual do ID para teste)
   - Body: 
     ```json
     {
       "cpf": "12345678900",
       "bio": "Eletricista experiente",
       "category": "Eletricista",
       "price_hour": 150.0,
       "latitude": -23.550520,
       "longitude": -46.633308
     }
     ```
2. **Buscar Profissionais**: `GET /profissionais/busca?lat=-23.55&long=-46.63&radius=10&category=Eletricista`

### Serviços
1. **Solicitar Serviço**: `POST /servicos`
   - Body: `{ "professional_id": "<uuid>", "description": "Conserto tomada", "scheduled_for": "2023-11-01T14:00:00Z" }`
2. **Avaliar**: `POST /servicos/<id>/avaliar`
   - Body: `{ "rating": 5, "comment": "Excelente!" }`

## 📝 Notas de Desenvolvimento

- **Autenticação**: O middleware de Auth valida apenas a presença do token. Em produção, descomentar a validação de assinatura JWT no `pkg/middleware/auth.go`.
- **Upload**: O upload de imagens é um stub retornando URLs fictícias. Integrar com Supabase Storage no `internal/service/professional_service.go`.

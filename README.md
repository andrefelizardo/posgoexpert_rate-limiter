# posgoexpert_rate-limiter

## Descrição

Rate limiter em Go que controla requisições com base em IP e token (via header "API_KEY").

## Configuração

Configure o arquivo `.env` na raiz com as seguintes variáveis:

- RATE_LIMIT_IP (ex.: 5)
- RATE_LIMIT_TOKEN (ex.: 10)
- BLOCK_TIME_IP (em segundos, ex.: 300)
- BLOCK_TIME_TOKEN (em segundos, ex.: 300)
- REDIS_ADDR (ex.: localhost:6379)

## Executando a Aplicação

1. Construa e inicie os serviços:
   > docker-compose up --build
2. Isso iniciará o Redis e a aplicação.

## Executando os Testes

### Com Docker Compose

1. Execute os testes com:
   > docker-compose run --rm app

### Localmente

1. Rode:
   > go test ./...

# CLOB em Go - Documentação Completa

## Visão Geral

CLOB (Central Limit Order Book) implementado em Go é um sistema completo de negociação de ativos que gerencia um livro de ordens centralizado com matching engine e controle de saldos. Esta implementação é minimalista e utiliza apenas a biblioteca padrão do Go, sem dependências externas.

## Características

- **Central Limit Order Book** completo
- **Matching Engine** em tempo real
- **Gestão de Saldos** com reservas e atualizações
- **Sistema totalmente thread-safe** usando mutexes
- **Implementado apenas com stdlib** do Go
- **Design modular** seguindo princípios Clean Architecture

## Build & Run

### Requisitos

- Go 1.25 ou superior

### Compilação e Execução

```bash
# Baixar dependências (se houver)
go mod tidy

# Compilar e executar diretamente
go run ./cmd/server

# Para compilar um binário
go build -o clob_server ./cmd/server

# Executar o binário compilado
./clob_server
```

Por padrão, o servidor escuta na porta `:3000`.

## Documentação API / Swagger

Para visualizar a documentação Swagger da API:

```bash
# Iniciar o servidor
go run ./cmd/server

# Abrir Swagger UI no navegador
$BROWSER http://localhost:3000/swagger/index.html
```

Você também pode usar as informações abaixo para interagir com a API.

## Endpoints da API

### Gestão de Contas

#### Criar uma Nova Conta

- **Método:** `POST`
- **URL:** `/accounts`
- **Descrição:** Cria uma nova conta de usuário no sistema
- **Exemplo:**
  ```bash
  curl -X POST http://localhost:3000/accounts
  ```
- **Resposta:** ID da conta criada

#### Creditar Saldo em uma Conta

- **Método:** `POST`
- **URL:** `/accounts/{id}/credit`
- **Descrição:** Adiciona saldo de um determinado ativo à conta
- **Corpo:**
  ```json
  {
  	"asset": "BTC",
  	"amount": 10000000 // Valores em centavos/satoshis (int64)
  }
  ```
- **Exemplo:**
  ```bash
  curl -X POST http://localhost:3000/accounts/123/credit -H "Content-Type: application/json" -d '{"asset":"BTC","amount":10000000}'
  ```

#### Obter Saldo de uma Conta

- **Método:** `GET`
- **URL:** `/accounts/{id}`
- **Descrição:** Retorna todos os saldos disponíveis e reservados por ativo
- **Exemplo:**
  ```bash
  curl http://localhost:3000/accounts/123
  ```
- **Resposta:**
  ```json
  {
  	"id": "123",
  	"balances": {
  		"BTC": {
  			"available": 5000000,
  			"reserved": 2000000
  		},
  		"BRL": {
  			"available": 1000000000,
  			"reserved": 0
  		}
  	}
  }
  ```

### Gestão de Ordens

#### Inserir Nova Ordem

- **Método:** `POST`
- **URL:** `/orders`
- **Descrição:** Insere uma nova ordem no livro de ofertas e executa matching imediatamente
- **Corpo:**
  ```json
  {
  	"account_id": "123",
  	"instrument": "BTC/BRL",
  	"side": "BUY", // ou "SELL"
  	"price": 50000000, // 500.000,00 BRL (em centavos)
  	"quantity": 100000000 // 1,0 BTC (em satoshis)
  }
  ```
- **Exemplo:**
  ```bash
  curl -X POST http://localhost:3000/orders -H "Content-Type: application/json" -d '{"account_id":"123","instrument":"BTC/BRL","side":"BUY","price":50000000,"quantity":100000000}'
  ```
- **Resposta:** ID da ordem criada e status de matching

#### Cancelar Ordem

- **Método:** `POST`
- **URL:** `/orders/{id}/cancel`
- **Descrição:** Remove uma ordem do livro e libera os recursos reservados
- **Exemplo:**
  ```bash
  curl -X POST http://localhost:3000/orders/order123/cancel
  ```

#### Consultar Livro de Ofertas

- **Método:** `GET`
- **URL:** `/book/{instrument}`
- **Descrição:** Retorna o estado atual do livro de ofertas para um instrumento
- **Exemplo:**
  ```bash
  curl http://localhost:3000/book/BTC/BRL
  ```
- **Resposta:**
  ```json
  {
  	"instrument": "BTC/BRL",
  	"bids": [
  		{ "price": 50000000, "quantity": 100000000 },
  		{ "price": 49000000, "quantity": 200000000 }
  	],
  	"asks": [
  		{ "price": 51000000, "quantity": 150000000 },
  		{ "price": 52000000, "quantity": 300000000 }
  	]
  }
  ```

## Decisão sobre Preço de Execução

**O preço de execução adotado nesta implementação é o preço da ordem resting (a que já estava no livro).**

### Justificativa:

1. **Transparência para o Market Maker**: Quem coloca uma ordem limit no livro sabe exatamente a que preço ela será executada quando ocorrer um match.

2. **Compatibilidade com mercados tradicionais**: A maioria das bolsas de valores e exchanges de criptomoedas adotam este modelo, onde o preço da ordem que já estava no livro é respeitado.

3. **Redução de incerteza**: Ordens market (sem preço limite) são executadas ao melhor preço disponível no livro, reduzindo surpresas.

4. **Price/Time Priority**: Mantém o princípio de prioridade por preço e tempo (FIFO), garantindo tratamento justo das ordens.

Este modelo garante que, se um usuário coloca uma ordem de compra a 500.000 BRL e alguém posteriormente coloca uma ordem de venda a 490.000 BRL, o match ocorre a 500.000 BRL, beneficiando o vendedor que recebe um preço melhor que seu limite. Isso incentiva os agressores de mercado a continuarem fornecendo liquidez.

## Arquitetura do Sistema

A arquitetura segue os princípios Clean Architecture:

```
/clob_go
├── cmd
│   └── server                 # Ponto de entrada da aplicação
├── internal
    ├── domain                 # Regras de negócio e entidades
    │   ├── account            # Entidades relacionadas a contas e saldos
    │   └── book               # Entidades do livro de ordens e matching
    │   └── order              # Entidades relacionadas a ordens
    ├── application            # Casos de uso da aplicação
    │   ├── account            # Casos de uso para gestão de contas
    │   └── book               # Casos de uso para gestão de books
    │   └── order              # Casos de uso para gestão de ordens
    └── infra                  # Implementações de infraestrutura
        ├── controllers        # Controladores HTTP
        └── repositories       # Implementações de repositórios
```

## Detalhes de Implementação

1. **Valores Monetários**: Todos os valores monetários são armazenados como `int64` para evitar problemas de precisão com ponto flutuante. Por exemplo, 1 BTC é representado como 100.000.000 satoshis e 500.000,00 BRL como 50.000.000 centavos.

2. **Thread Safety**: Todas as operações críticas são protegidas por mutexes para garantir consistência em ambientes concorrentes.

3. **Matching Engine**: O matching ocorre em tempo real quando uma nova ordem é inserida. O algoritmo busca pares compatíveis no livro de ofertas, gerando um ou mais trades quando os preços se cruzam.

4. **Gestão de Saldos**:
   - Ao inserir uma ordem de compra, o valor (preço × quantidade) é reservado no saldo de quote (ex: BRL)
   - Ao inserir uma ordem de venda, a quantidade é reservada no saldo de base (ex: BTC)
   - Quando um match ocorre, os saldos reservados são consumidos e os novos ativos são creditados nas contas

## Exemplos de Fluxo Completo

### 1. Inicializando o Sistema e Criando Contas

```bash
# Iniciar o servidor
go run ./cmd/server

# Criar duas contas
curl -X POST http://localhost:3000/accounts
# Resposta: {"id":"acc1"}

curl -X POST http://localhost:3000/accounts
# Resposta: {"id":"acc2"}

# Creditar BTC na conta 1
curl -X POST http://localhost:3000/accounts/acc1/credit -H "Content-Type: application/json" -d '{"asset":"BTC","amount":100000000}'

# Creditar BRL na conta 2
curl -X POST http://localhost:3000/accounts/acc2/credit -H "Content-Type: application/json" -d '{"asset":"BRL","amount":50000000000}'
```

### 2. Inserindo Ordens e Match

```bash
# Conta 1 coloca ordem de venda (1 BTC por 500.000 BRL)
curl -X POST http://localhost:3000/orders -H "Content-Type: application/json" -d '{"account_id":"acc1","instrument":"BTC/BRL","side":"SELL","price":50000000,"quantity":100000000}'
# Resposta: {"id":"ord1","status":"OPEN"}

# Verificar o livro de ofertas
curl http://localhost:3000/book/BTC/BRL
# Mostra a ordem de venda no livro

# Conta 2 coloca ordem de compra que vai cruzar (1 BTC por 510.000 BRL)
curl -X POST http://localhost:3000/orders -H "Content-Type: application/json" -d '{"account_id":"acc2","instrument":"BTC/BRL","side":"BUY","price":51000000,"quantity":100000000}'
# Resposta: {"id":"ord2","status":"FILLED","trades":[{"instrument":"BTC/BRL","quantity":100000000,"price":50000000,"buyer_id":"acc2","seller_id":"acc1"}]}

# Verificar saldos após o match
curl http://localhost:3000/accounts/acc1
# Mostra 0 BTC e 50.000.000 BRL

curl http://localhost:3000/accounts/acc2
# Mostra 1 BTC e 45.000.000 BRL (original - 50.000.000)
```

### 3. Cancelamento de Ordem

```bash
# Colocar uma nova ordem
curl -X POST http://localhost:3000/orders -H "Content-Type: application/json" -d '{"account_id":"acc1","instrument":"BTC/BRL","side":"SELL","price":60000000,"quantity":50000000}'
# Resposta: {"id":"ord3","status":"OPEN"}

# Cancelar a ordem
curl -X POST http://localhost:3000/orders/ord3/cancel

# Verificar que a ordem não está mais no livro
curl http://localhost:3000/book/BTC/BRL
# A ordem ord3 não deve mais aparecer
```

## Testes

```bash
Executar lint no projeto
make lint

# Executar todos os testes e gerar cobertura dos testes e rodar o lint no projeto
make test-lint

# Executar todos os testes
go test -v ./...

# Executar testes específicos
go test -v ./internal/domain/book/...

# Verificar cobertura
go test -cover ./...

# Gerar relatório de cobertura HTML
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Considerações e Limitações

1. **Implementação em Memória**: Esta implementação mantém todos os dados em memória, sem persistência.

2. **Sem Autenticação/Autorização**: O sistema não implementa mecanismos de autenticação ou autorização.

3. **Ordem Limit Apenas**: Apenas ordens limit são suportadas. Ordens market, stop ou OCO não estão implementadas.

4. **Sem Taxas**: O sistema não considera taxas nas operações.

5. **Matching Simples**: O matching é feito com base apenas no melhor preço disponível, sem considerar estratégias mais complexas.

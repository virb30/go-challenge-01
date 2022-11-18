# Desafio 01

Webserver http, contextos, banco de dados e manipulação de arquivos com Go.
 
## Descrição

Sistemas em go
- client.go
- server.go
  
O client.go realiza uma requisição HTTP no server.go solicitando a cotação do dólar.
 
O server.go consome a API contendo o câmbio de Dólar e Real no endereço: https://economia.awesomeapi.com.br/json/last/USD-BRL e retorna no formato JSON `{"valor": 1.2345}` o resultado para o cliente.
 
O server.go registra no banco de dados SQLite cada cotação recebida, com um timeout máximo para chamar a API de cotação do dólar de 200ms e timeout máximo para conseguir persistir os dados no banco de 10ms.
 
O client.go recebe do server.go apenas o valor atual do câmbio (campo "bid" do JSON - convertido para "valor"). O client.go tem um timeout máximo de 300ms para receber o resultado do server.go.
 
O client.go salva a cotação atual em um arquivo "cotacao.txt" no formato: Dólar: {valor}
 
O endpoint gerado pelo server.go é: /cotacao e a porta 8080.

### Como usar

Clonar o repositório e entrar na pasta

```bash
#clonar o repositório
git clone github.com/virb30/go-challenge-01 .
# entrar no diretório
cd go-challenge-01
```

Instalar módulos necessários e iniciar Server:

```bash
# navegar até a pasta do server
cd server
## instalar dependências
go mod tidy
# iniciar server
go run server.go
```

Client

```bash
# navegar até a pasta do client
cd client
# executar client
go run client.go
```
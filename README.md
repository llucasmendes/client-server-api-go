
# Client-Server API em Go

Este projeto consiste em dois sistemas em Go que se comunicam para obter a cotação do dólar e registrar essa informação em um banco de dados SQLite.

## Estrutura do Projeto

A estrutura do projeto é a seguinte:

```
client-server-api-go/
├── client/
│   └── client.go
├── server/
│   └── server.go
├── go.mod
└── go.sum
```

- `client/`: Contém o código do cliente que faz uma requisição HTTP ao servidor para obter a cotação do dólar e salva o resultado em um arquivo `cotacao.txt`.
- `server/`: Contém o código do servidor que obtém a cotação do dólar de uma API externa, salva a cotação em um banco de dados SQLite e responde ao cliente com o valor da cotação.

## Configuração e Execução

### Pré-requisitos

- Go instalado na máquina. [Download Go](https://golang.org/dl/)
- Dependências Go (SQLite driver)

### Passos para Execução

1. Clone o repositório:
   ```sh
   git clone <URL_DO_REPOSITORIO>
   cd client-server-api-go
   ```

2. Inicialize o módulo Go:
   ```sh
   go mod tidy
   ```

3. Instale o driver SQLite:
   ```sh
   go get -u github.com/mattn/go-sqlite3
   ```

4. Navegue até o diretório `server` e execute o servidor:
   ```sh
   cd server
   go run server.go
   ```

   O servidor estará escutando na porta `8080`.

5. Abra outro terminal, navegue até o diretório `client` e execute o cliente:
   ```sh
   cd ../client
   go run client.go
   ```

### Funcionamento

1. **Servidor (`server.go`)**:
   - Inicia um servidor HTTP na porta `8080`.
   - Define o endpoint `/cotacao` que faz uma solicitação para a API de câmbio de dólar para real.
   - Salva a cotação obtida no banco de dados SQLite com um timeout de 10ms.
   - Retorna a cotação no formato JSON para o cliente.

2. **Cliente (`client.go`)**:
   - Faz uma requisição HTTP para o servidor no endpoint `/cotacao` com um timeout de 300ms.
   - Recebe a cotação do servidor e salva no arquivo `cotacao.txt` no formato `Dólar: {valor}`.

### Erros e Tratamento

- Logs de erro são gerados se os tempos de execução excederem os limites configurados.
- O cliente e o servidor utilizam contextos para definir timeouts e garantir que as operações não ultrapassem os tempos definidos.

### Licença

Este projeto é distribuído sob a licença MIT. Veja o arquivo `LICENSE` para mais detalhes.
```
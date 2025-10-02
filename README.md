# katseye

## Docker

### Configurar variáveis de ambiente

O Docker Compose utiliza o arquivo `.env.docker`. Ajuste os valores conforme necessário (por exemplo, redefina `JWT_SECRET` e as credenciais do MongoDB). O arquivo padrão já inclui a criação de um usuário de aplicação e credenciais raiz para o Mongo.

### Build da imagem

```
docker build -t katseye-api .
```

### Subir a stack completa (API + MongoDB + Redis)

```
docker compose up --build
```

A API ficará disponível em `http://localhost:8080`, enquanto o MongoDB expõe a porta `27017` e o Redis a porta `6379`.

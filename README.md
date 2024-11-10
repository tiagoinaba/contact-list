# Lista de contatos
Esta é a minha solução para o teste prático da DavinTI. Implementei em 4 horas,
utilizando um stack não muito convencional, mas que utilizo sempre por ser
extremamente produtivo.

## Tecnologias
- Go
- HTMX
- SQLite

Escolhi essas tecnologias, pois, além de super produtivas, são leves e rápidas.
A biblioteca [HTMX](https://htmx.org/) me permite criar soluções interativas na
web sem precisar escrever muito JavaScript, podendo focar na funcionalidade. O
SQLite escolhi por sua natureza que dispensa servidores e a facilidade de
exportação dos dados. Por último, escolhi Go para o back-end, pois é a
linguagem que mais utilizo no meu dia a dia.

## Como executar?

### Docker
```
docker build -t davinti .
docker run -it -p 8080:8080/tcp --rm --name davinti-app davinti
```

### Go
```
go mod tidy
go run cmd/main.go
```

Shopping cart
======

## 1. The API
To check full API documentation, please start the service and visit [Swagger API page](http://localhost:8080/swagger/).
Instructions of how to run service locally can be found [here](#3-how-to-run-service-locally)

## 2. Authentication
Service is using Basic Authentication.
Example header: 
```
Authorisation: Basic dXNlcjpwYXNzd29yZAo=
```
where we pass user:password base64 encoded.


## 3. How to run service locally

1. spin up mysql `./mysql_docker.sh`
2. install goose schema migration tool `GO111MODULE=off go get -u github.com/pressly/goose/cmd/goose`
3. run migrations `$(go env GOPATH)/bin/goose -dir=./migrations mysql "shoppingcart:secret@/shoppingcart?parseTime=true" up`
4. `source .env.local` Set up default environment for local
5. Start app `go run ./cmd/shoppingcart/`
Shopping cart
======
The Shopping cart service provides a basic functionality to work with shopping carts, like:
* Create a shopping cart
* Get information about the shopping cart
* Empty the shopping cart
* Add a product to the shopping cart
* Remove a product from the shopping cart

This service doesn't hold information about products, users or orders. 
In order to authorize the user, the service is using Auth service (mocked). 


## 1. The API
To check full API documentation, please start the service locally and visit [Swagger API page](http://localhost:8080/swagger/).
Instructions of how to run the service locally can be found [here](#3-how-to-run-service-locally)

## 2. Authentication
Service is using Basic Authentication. 
In order to verify user credentials Auth service is used. It's mocked, so any credentials will work.

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

## 4. How to run tests

`go test -v github.com/bugimetal/shoppingcart/...`

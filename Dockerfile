# Build 1
FROM golang:1.14-alpine as builder

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o shoppingcart /app/cmd/shoppingcart/
CMD ["/app/shoppingcart"]

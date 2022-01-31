FROM maerskao.azurecr.io/alcl/golang-1.16-alpine:2cd1b4abe716a9e2659f6e96cca066367f5a9754 AS builder

WORKDIR /app

COPY go.mod go.sum Makefile ./
RUN make mod

COPY . .
RUN make build

FROM maerskao.azurecr.io/alcl/alpine:3.14

COPY --from=builder /app/bin/alcl-go-function-template /alcl-go-function-template

ENTRYPOINT [ "/alcl-go-function-template" ]

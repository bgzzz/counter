FROM golang:alpine AS builder
RUN apk add --update --no-cache make

WORKDIR /app

COPY go.mod go.sum Makefile ./
RUN make mod

COPY . .
RUN make build

FROM golang:alpine

COPY --from=builder /app/bin/counter /counter
#in case of running it from the container
COPY --from=builder /app/bin/counterctl /counterctl

ENTRYPOINT [ "/counter" ]

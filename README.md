# counter

This repo contains implementation of the HTTP server handling basic logic
of counter application. In addition to that this repo contain client
 application to interact with the server. 

Counter application is a HTTP server exposing API to increment, decrement 
and get the counter value.

## HOWTOs

### Server 

#### Run server

Help example: 
```bash
➜  counter git:(readme) ✗ ./bin/counter --help
NAME:
   counter-server - counter-server is server side application returning the value of the counter

USAGE:
   counter [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --log-level value  log-level "debug" (more on the supported levels here: https://github.com/sirupsen/logrus/blob/fdf1618bf7436ec3ee65753a6e2999c335e97221/logrus.go#L25) (default: "debug") [$LOG_LEVEL]
   --port value       --port 8080 (default: 8080) [$PORT]
   --help, -h         show help (default: false)
``` 

By default server is running on port 8080, in order to run it with 
default parameters you need to execute the binary:

```bash
➜  counter git:(client) ✗ ./bin/counter
{"level":"info","msg":"Running server on port: 8080","time":"2022-01-31T15:33:15+01:00"}
{"level":"debug","msg":"increment counter with value %d was executed1","time":"2022-01-31T15:33:21+01:00"}
{"level":"debug","msg":"get the counter with value %d was executed1","time":"2022-01-31T15:33:21+01:00"}
{"level":"debug","msg":"increment counter with value %d was executed2","time":"2022-01-31
```

#### Interact with the server 

You can use curl to interact with the server. It supports POST, GET, DELETE 
HTTP methods without payload on the URI api/v1/counter

##### Get counter 
```bash
➜  counter git:(readme) ✗ curl -X GET localhost:8080/api/v1/counter
{"Counter":0}% 
```

##### Increment counter 
```bash
➜  counter git:(readme) ✗ curl -X POST localhost:8080/api/v1/counter
{"Counter":1}% 
```

##### Decrement counter 
```bash
➜  counter git:(readme) ✗ curl -X DELETE localhost:8080/api/v1/counter
{"Counter":0}% 
```

##### Edge cases

There is an edge cases support due to the fact that counter value is uint64
which means server will respond with an error message and non 2XX HTTP code
any time when there will be a client attempt to decrease the counter to less 
then 0 and increase the counter higher then unit64 max. 

Example:
```bash
➜  counter git:(readme) ✗ curl -X DELETE localhost:8080/api/v1/counter
unable to decrement, counter has reached its minimum value
```

### Build the project 

To build the project you need to execute following command

```bash
➜  counter git:(readme) make build
CGO_ENABLED=0 go build -o ./bin/counter ./cmd/counter
CGO_ENABLED=0 go build -o ./bin/counterctl ./cmd/counterctl
```

This execution will create two binaries counter for server side 
application and counterctl for client side CLI.

### Run the tests 

```bash
➜  counter git:(readme) make test
CGO_ENABLED=0 go test -v ./...
?       github.com/bgzzz/counter/cmd/counter    [no test files]
?       github.com/bgzzz/counter/cmd/counterctl [no test files]
=== RUN   TestClientGet
=== RUN   TestClientGet/client_get_table_test_0
=== PAUSE TestClientGet/client_get_table_test_0
=== RUN   TestClientGet/client_get_table_test_1
=== PAUSE TestClientGet/client_get_table_test_1
=== RUN   TestClientGet/client_get_table_test_2
...
```

### All to gather 

```bash
➜  counter git:(readme) ✗ make
CGO_ENABLED=0 go mod download
CGO_ENABLED=0 go install github.com/golangci/golangci-lint/cmd/golangci-lint
CGO_ENABLED=0 go test -v ./...
?       github.com/bgzzz/counter/cmd/counter    [no test files]
?       github.com/bgzzz/counter/cmd/counterctl [no test files]
=== RUN   TestClien
```

### Dockerization 

If you want to run binaries in a docker env use the following command to create an image:
```bash
➜  counter git:(readme) ✗ docker build -t counter .
[+] Building 59.9s (15/15) FINISHED                                                   
 => [internal] load build definition from Dockerfile                             0.0s
 => => transferring dockerfile: 383B                                             0.0s
 => [internal] load .dockerignore            
```

## Remarks

1. Code is structured such that there are two separate apps (/pkg/client, pkg/server) bind with common contract model  
2. Prod ready implementation can be implemented via openAPI contract definition and code-generation of client server primitives. Contract test tooling can be applied afterwards.   
3. Atomic package can be used if we can afford to cycle uint counter
4. Would be good to implement opentelemetry since it is a communication of two services.
5. Would be good to implement timeouts and context processing, also taking into account that timeouts and retries can be handled by infra components like service mesh
6. Potentially different (more convenient router/framework (echo, chi, gorilla)) can be used   
7. Would be good to cover all functions with unit tests.  
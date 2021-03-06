# MEMO

# DB Migration 

local to docker container
docker container to local

0. make db architecture

1.
migrate create -ext sql -dir db/migration -seq init_schema
-seq means migrate version name.

2. 
copy sql and paste to init_schema.up.sql

3.
write drop table sth in a down file

3.
migrate -path db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up


brew install postgresql
brew install tableplus
brew install sqlc // sql to go 

## pq - A pure Go postgres driver for Go's database/sql package
go get github.com/lib/pq
go get github.com/stretchr/testify
                                  /require

## CRUD sqlc

move to root directory
 -> command "sqlc init"
 -> set sqlc.yaml file
 -> command "sqlc generate"


# Unit tests

When writing Unit tests,
Should make sure that they are independent from each other.
Why: hard to maintain if has hundred of tests that depends on each other.

# Transaction

DB Transaction

Why?
1. To provide a reliable and consistent unit of work, even in case of system failure.
2. To provide isolation between programs that access the database concurrently.

A: Atomicity
C: Consistency
I: Isolation
D: Durability

Closure is often used when want to get the result from a callback function.
the callback function itself doesn't know the exact type of the result it should return

## Caution when Database Transaction
concurrency carefully!
The best way that is run it with several concurrent go routines

## Update Account within Transaction
Require careful handing of concurrent transactions to avoid deadlock
database locking 

# TDD | test driven development
Tests first to make our current code breaks

* Test command
go test -v -coverprofile cover.out ./...
go tool cover -html=cover.out

# Deadlock postgresql
https://wiki.postgresql.org/wiki/Lock_Monitoring
Deadlock ocuurs because 2 concurrent transactions both need to wait for each other

can resolve deadlock consider order of transactions

SELECT * FROM accounts WHERE id = $1 LIMIT 1
+ FOR NO KEY UPDATE

## Isolation level
dirty read | read uncommitted 
non-repeatable read |read committed
phantom read | repeatable read
serialization anomaly | serializable

postgresql tx begin first when change isolation level.
postgresql not working level of dirty read | read uncommitted.
postgresql uses a dependencies checking mechanism
to detect potensial read phenomena and stop them by throwing error

# Github Actions | workflows
Golang Unit tests in external postgres service in github using .github/workflows/ci.yml

project folder -> 
    mkdir -p .github/workflows
    touch .github/workflows/ci.yml


# HTTP API

Gin.
the router field(package) is private so that Start function makes to access api package

* _ "github.com/lib/pq"

# Viper
golang packege, dealing with configuration file
ex) ENV, YAML ...

app.env
UPPERCASE=123

# Mock DB

make store to mock using interface.
interface cannot be pointer  

gomock makes fake DB Unit tests easily.
through the interface about DB service, mockgen makes Mock Service for Unit tests

# Test Multiple senarios
using t.Run()

## struct `` options -> nedd validator
binding:"oneof= " -> How to avoid hard-coding -> just 'currency'

import "github.com/go-playground/validator/v10"
var validCurreny validator.Func
  -> binding.Validator.Engine() (server.go)

# User Authentication and Authorization

## Update DB schema from previous version 

add foreign key and unique constraint
ref: 
    > many to one
    < one to many
    - one to one

make unique constraint
  indexes
  (field1, field2) [unique] 

Right way to migrate!
if you want to apply a new schema change,
create a new migration version.

command "migrate -verbose up 1 || down 1" -> choose version when migration.
so that can backup the db version

## handle DB errors for adding users migration

change owner -> user.Username(foreign key) to account test file.
command "make mock" regenerate
  -> added User interface to mock store file

Handle error code to the others
400 403 404
500

500 -> client accesss db constraits error code -> 403 

# bcrypt
make hashed password

# gomcok matcher

hash and salt makes pw test difficult
using matcher interface can make custom matcher

just implemented, unit test be stronger 

# JWT To PASETO

JSON Web Token
"Header.Payload.Signature"

JWT Signing Algorithms

1. Symmetric digital signature algorithm
- the same secret key is used to sign & verify Token
- For local use: internal services, where the secret key can be shared
- HS256, HS384, HS512
  - HS256 = HMAC + SHA256
  - HMAC: Hash-based Message Authentication Code
  - SHA: Secure Hash Algorithm
  - 256/384/512: number of output bits

2. Asymmetric disigtal signature algorithm
- The private key is used to sign Token
- The publi key is used to verify Token
- For public use: internal service signs token, but external service needs to verify it
- RS256, RS384, RS512 || PS256, PS384, PS512 || ES ...
  - RS256 = RSA PKCSv1.5 + SHA256 [PKCS: Public-Key cryptography Standards]
  - PS256 = RSA PSS + SHA256 [PSS: Probabilistic Signature Schema]
  - ES256 = ECDSA + SHA256 [ECDSA: Elliptic CUrve Digital Signature Algorithm]

Problem of JWT
  - weak Algorithms
  - Trivial Forgery(??????)


  PASETO
  Platform-Agnostic SEcurity TOkens

  - Stronger algorithm 
  (only need to select the version of PASETO)
  (Only 2 mosst recenbt PASETO versions are accepted)

  local -> Symmetric
  public -> Asymmetric
  
## Create JWT & PASETO token

interface 2 method

CreateToken
VerifyToken

go get github.com/google/uuid

Implementation of interface
  - Add Method for interface
  - Where the struct required function(method) of the interface
  - Implement method
  - Write func () MethodName() {}

Keyfunc is that receives the parsed but unverified token
verify its header, about siging algorithm matches

Test JWTToken function
both happy and error case

jwt to paseto
go get github.com/o1egl/paseto

PASETO CreaetToken -> paseto.Encrypt
PASETO VerifyToken -> paseto.Decrypt

## Login API with token

Add token config
For test, 
make newTestServer with token maker instead of call NewServer

In loginUser,
sensitive data inside db.User struct -> function Upper to lower

## Authentication middleware and Authorization rules

Authorization in Header (API)
access-token belongs specipic user, should not be able to access other users

Using Gin Middleware (Authorization)
  ! -> context.Abort()
  ok -> action api func

Ways
  1. Extract Authorization header from the request
  2. 

  *Multiple Test

    testCases := []struct {
      name          string
      setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
      checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
      ...
    }{}

    for i := range testCases {

      tc := testCases[i]
      t.Run(tc.name, func(t *testing.T) {
        server := newTestServer(t, nil)
        ...
      })
    }

Gin Router Group
  -> .Group("/path").Use(middleware func)
  Add middleware all of routes in the group

* gomock test rules
buildStubs: func(store *mockdb.MockStore) {
  // error point
  store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account3.ID)).Times(1).Return(account3, nil)
  // after db times to 0
  store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(0)
  store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
},

# Add docker

1. git

git checkout -b "new_branch"
master -> new branch -> update code -> test
-> merge to master!

2. mac(local) and program package version check

3. Dockerfile
make Dockerfile in root
to define the base image to build our app
(based Golang -> needs Golang image)

alpine is small image

    FROM golang:version 
      -> to specify the base image

    WORKDIR /app
      -> to declare the current working directory inside the image

    COPY . .
      -> first, copy the current folder files
      -> second paste to the WORKDIR setting path (place to store)

    RUN go build -o main main.go
      -> build our app to a single binary exec file
      -> -o is setting output file

    EXPOSE 0000
      -> the container listens on the specified network port at runtime
      -> 0000 -> port number
      -> noting that the EXPOSE doesn't actually publish the port

    CMD ["/app/main"]
      -> to define the default command to run when the container starts

  *How to make more smaller?*
  Multistage
  ## Build stage(first stage)
    FROM + AS builder
  ## Run stage(second stage)
    FROM alpine:3.13
    WORKDIR /app
    COPY --from=builder /app/main . # first stage's (--from is option)

  *Make Image*
  docker build -t image:latest .

  
# Docker Network to connect container

  *Run Image*
  docker run --name app -p PORT:PORT image:latest

  config file .env -> development to production config file
  COPY app.env .

  docker run --name app -p PORT:PORT -e GIN_MODE=release image:latest

To allow 2 stand-alone containers to talk to each other by names
 
  *Golang & PostgreSQL container network setting is different (about localhost)
  docker container inspect container_name
  172.17.0.3 / 172.17.0.2

  but when rebuild image, it IP address will change
  So that the better way is do not rebuild
  -> using viper to read the config
  -> override config file 
  -> Add command into docker run : -e DB_SOURCE="postgres://user:pw@Inspect_IP:PORT/app?sslmode=disable"

  *but when rerun the new container, IP address will change
  So that using user-defined network instead!

  -> check networks
  docker network ls
  docker network inspect bridge

  bridge is offering default IP to container
  -> to create own network
  docker network create new_network_name
  docker network connect new_network_name
  -> docker run app
  -> --network created_network_name
  -> IP address changes to container_name (same network's container name)
  docker run --name app --network network_name -p PORT:PORT -e GIN_MODE=release -e DB_SOURCE="postgres://user:pw@container_name:PORT/app?sslmode=disable" image:latest
  
  now using its name instead of the IP address.

  config file: Inspect_IP -> localhost

* Reviews Code in Github
In Pull requests, add comment about some specific code lines

# How to write docker-compose.yml 

- follow file reference Version 3

version: "3.9"
services:
  postgres:
    image: postgres:12-alpine
    environment:
      - POSTGRES_USER=root
      ...

  // service name
  api:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - localhost to name connected the same network (postgres)
      - DB_SOURCE=postgres...

  - docker compose up
    -> Network simplebank_default       Create...
    -> New network includes postgres and api service 
  - docker compose down 
    -> down container and remove images

  - Edit Dockerfile and docker-compose.yml for migration
  specially Dockerfile Add RUN command "apk add curl" for that

  - Add start.sh in root
  chmod +x start.sh

  ### !/bin/sh
  ### will be run by /bin/sh
  ### alpine image
  ### bash is not available

  set -e

  echo "run db migration"
  /app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

  echo "start the app"
  ### takes all parameters passed to the script and run it 
  exec "$@"

  *CMD & ENTRYPOINT (CMD in Dockerfile reference)
  : if CMS id used to provide default arguments for the ENTRYPOINT instruction, 
    both the CMD and ENTRYPOINT instructions should be specified woth the JSON array format

  - docker compose up again (Error handling)
    -> the postgres server was not ready to accept connection with api service yet
    -> Need to wait in docker-compose.yml (depends_on:)

  *depends_on (Control startup and shutdown order in Compose)
  : does not wait for db and redis to be "ready" before starting web only until they have been started.

  - wait-for
    -> COPY wait-for . (Dockerfile)
    -> entrypoint: [ "/app/wait-for", "db:port", "--", "/app/start.sh"] (docker-compose.yml)
       command: [ "/app/main" ]

# AWS 

ECR 
  - a fully-managed docker container registry
  - makes it easy to store, manage and deploy docker container images.

  1. make repository (= docker image store)
  2. image push. But can use github actions to automatically (build tag push, when to merge master branch)
  3. add yml file for github actions
  4. IAM setting for credential AWS ECR
  5. github -> setting -> secret -> actions: add new accecc key and secret key ex name: AWS_ACCESS_KEY_ID (in yml)
  
RDS
  - don't have to care to maintain or scale the DB cluster

  Point
    1. select option
    2. set vpc security group 

Secrets Manager
  - Set SM for RDS 
  - other keys (from .env file`)
  - Input Keys and Values
  
  command for rand token symmetric key
  openssl rand -hex 64 -> 64bytes 
  (get 128 charactors rand string)
  openssl rand -hex 64 | head -c 32
  (get 32 charactors rand string)

  

AWS cli with iAM
  0. download cli
  1. Create new key
  2. command "aws configure"
  3. set key and secret
  4. see the file by command "ls -l ~/.aws"
  5. see the credentials by command "cat ~/.aws/credentials"

  help
  : aws [service_name] help

SM using AWS cli (Store & retrieve production secrets with AWS secrets manager)
  aws secretsmanager get-secret-value --secret-id simple_bank
  {
    "ARN": "arn:aws:secretsmanager:ap-northeast-1:312901933285:secret:simple_bank-AvZaeI",
    "Name": "simple_bank",
    "VersionId": "8eb30622-d8a8-4764-bf76-0f2b95f0d7c8",
    "SecretString": "{\"DB_DRIVER\":\"postgres\",\"SERVER_ADDRESS\":\"0.0.0.0:8080\",\"ACCESS_TOKEN_DURATION\":\"15m\",\"TOKEN_SYMMETRIC_KEY\":\"6e35ec823dcc0664a9e715ca3cfb10ed\"}",
    "VersionStages": [
        "AWSCURRENT"
    ],
    "CreatedDate": "2022-05-19T16:49:02.440000+09:00"
  }

-> Makes SecretString to like app.env file (data structure)
  1. aws secretsmanager get-secret-value --secret-id simple_bank --query SecretString

    "{\"DB_DRIVER\":\"postgres\",\"SERVER_ADDRESS\":\"0.0.0.0:8080\",\"ACCESS_TOKEN_DURATION\":\"15m\",\"TOKEN_SYMMETRIC_KEY\":\"6e35ec823dcc0664a9e715ca3cfb10ed\"}"
  2. 
    brew install jq
    jq ''

    aws secretsmanager get-secret-value --secret-id simple_bank --query SecretString --output text | jq 'to_entries'
    [
      {
        "key": "DB_DRIVER",
        "value": "postgres"
      },
      {
        "key": "SERVER_ADDRESS",
        "value": "0.0.0.0:8080"
      },
      {
        "key": "ACCESS_TOKEN_DURATION",
        "value": "15m"
      },
      {
        "key": "TOKEN_SYMMETRIC_KEY",
        "value": "6e35ec823dcc0664a9e715ca3cfb10ed"
      }
    ]

  3.
    aws secretsmanager get-secret-value --secret-id simple_bank --query SecretString --output text | jq 'to_entries|map(.value)' 
    [
      "postgres",
      "0.0.0.0:8080",
      "15m",
      "6e35ec823dcc0664a9e715ca3cfb10ed"
    ]

  4. 
    aws secretsmanager get-secret-value --secret-id simple_bank --query SecretString --output text | jq 'to_entries|map("\(.key)=\(.value)")' 
      [
        "DB_DRIVER=postgres",
        "SERVER_ADDRESS=0.0.0.0:8080",
        "ACCESS_TOKEN_DURATION=15m",
        "TOKEN_SYMMETRIC_KEY=6e35ec823dcc0664a9e715ca3cfb10ed"
      ]
      -> + |.[] (remove [])
  Final.
    aws secretsmanager get-secret-value --secret-id simple_bank --query SecretString --output text | jq -r 'to_entries|map("\(.key)=\(.value)")|.[]' 
    -> -r remove string
    -> + > app.env (to send)

* git checkout . 
(reset all the changes)
  
  Edit deploy.yml file
  - name: Load secrets and save to app.env
    run: aws secretsmanager get-secret-value --secret-id simple_bank --query SecretString --output text | jq -r 'to_entries|map("\(.key)=\(.value)")|.[]' > app.env 

  no need to install jq in ubuntu (already installed)

* write source app.env to start.sh file
  - need to set app.env file before db migration

# gRPC

  ### sequential numbering (= field tag)
  - 1-15 = until 1byte -> save memory
  - ^16 = starting 2bytes
  - string username = 1;
  - string fullname = 2;

  ### mkdir proto
    setup proto file
    setup service using proto file

  ### make command

  add --proto_path=proto 

  proto: 
    protoc --proto_path=proto --go_out=. --go_opt=paths=source_relative \
      --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
      proto/*.proto

  ### make start server func for grpc
  
  ### Evans is the best gRPC client interactive console

  ### Mac install command
  - brew tap ktr0731/evans
  - brew install evans
  * Connect server from client
    - evans --host localhost --port 9090 -r repl 
  // Default port is 50051

  ### show all of gRPC service func to be able to use from client
  - show service

  ###  call service func to server
  - call CreateUser

  ### Implement UnimplementedSimpleBankServer
    Don't have to validate, it was already proccessed

  ### error 
  status.Errorf(codes.Unimplemented, "message")

  ### gRPC gateway
  : serve gRPC, HTTP request at the same time

  #### Ways 
  https://github.com/grpc-ecosystem/grpc-gateway
  
  - mkdir tools
  - touch tools.go

  - Add code
  ```
    package tools

    import (
        _ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
        _ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
        _ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
        _ "google.golang.org/protobuf/cmd/protoc-gen-go"
    )
  ```

  - go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc

  ### main.go
  - make func runGatewayServer() ... 
  - using mux // mux http to grpc
  
  ```
    mux := http.NewServeMux()
    mux.Handle("/", grpcMux)

    grpcServer := grpc.NewServer()
    pb.RegisterSimpleBankServer(grpcServer, server)

    reflection.Register(grpcServer)
    
    err = http.Serve(listener, mux)
  ```

# Swagger (Swagger Hub & UI)
Tool is easy to make Open API document 

## add proto options

### swagger hub (Charged version)
 
- update make proto command from Makefile
`--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simplebank`
- Sign up swagger hub
- `Import json file` from proto options command
- Can update to add `protoc-gen-openapiv2 options`

### swagger-ui (Free version)

- git clone `https://github.com/swagger-api/swagger-ui.git`
- all of file in `dist` dir cp to doc/swagger
- change url to origin swagger json filename in `swagger-initializer.js`
- update server code
  in runGatewayServer(main.go)
	```
  fs := http.FileServer(http.Dir("./doc/swagger"))
	http.Handle("/swagger/", http.StripPrefix("/swagger/", fs))
  ```

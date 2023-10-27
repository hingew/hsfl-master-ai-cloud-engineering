# User Service

[![User-service](https://github.com/hingew/hsfl-master-ai-cloud-engineering/actions/workflows/user-service.yml/badge.svg?branch=develop)](https://github.com/hingew/hsfl-master-ai-cloud-engineering/actions/workflows/user-service.yml)

## Description

Authentication of the user

## Configuration

Create a private key for signing JWT tokens:

``` 
ssh-keygen -t ecdsa -f user-service.pem -m pem
```

To configure the user service, create a `config.yml` with the following content:

```yml
database:
    host: localhost
    port: 5432
    username: postgres
    password: password
    dbname: postgres
jwt:
    signKey: /path/to/user-service.pem
```


## Up and running

```sh
go build
./user-service --config=config.yml
```

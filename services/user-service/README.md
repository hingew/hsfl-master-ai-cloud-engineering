# User Service

[![User-service](https://github.com/hingew/hsfl-master-ai-cloud-engineering/actions/workflows/user-service.yml/badge.svg?branch=develop)](https://github.com/hingew/hsfl-master-ai-cloud-engineering/actions/workflows/user-service.yml)
[![codecov](https://codecov.io/gh/hingew/hsfl-master-ai-cloud-engineering/graph/badge.svg?token=CDPMA4XLME&flag=user-service)](https://codecov.io/gh/hingew/hsfl-master-ai-cloud-engineering)

## Description

Authentication of the user

## Configuration

Create a private key for signing JWT tokens:

``` 
ssh-keygen -t ecdsa -f user-service.pem -m pem
```

Put the contents of the `user-service.pem` into the `AUTH_SIGN_KEY` environment variable, or paste it into the 
`.env` file at the project root.

## Up and running

```sh
go build

export USE_TESTDATA=$USE_TESTDATA
export POSTGRES_HOST=$POSTGRES_HOST
export POSTGRES_PORT=$POSTGRES_PORT
export POSTGRES_USERNAME=$POSTGRES_USERNAME
export POSTGRES_PASSWORD=$POSTGRES_PASSWORD
export POSTGRES_DBNAME=$POSTGRES_DBNAME
export AUTH_SIGN_KEY=$AUTH_SIGN_KEY

./user-service
```

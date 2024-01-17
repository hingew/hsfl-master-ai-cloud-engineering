# User Service

[![User-service](https://github.com/hingew/hsfl-master-ai-cloud-engineering/actions/workflows/user-service.yml/badge.svg?branch=develop)](https://github.com/hingew/hsfl-master-ai-cloud-engineering/actions/workflows/user-service.yml)
[![codecov](https://codecov.io/gh/hingew/hsfl-master-ai-cloud-engineering/graph/badge.svg?token=CDPMA4XLME&flag=user-service)](https://codecov.io/gh/hingew/hsfl-master-ai-cloud-engineering)

## Description

Authentication of the user

## Configuration

Create a private key for signing JWT tokens:

``` 
openssl genrsa -out key.pem 2048
openssl rsa -in key.pem -outform PEM -pubout -out public.pem
```

Update the content of your `.env` file in the project root.
```
AUTH_PRIVATE_KEY="<your-private-key>"
AUTH_PUBLIC_KEY="<your-public-key>"
```

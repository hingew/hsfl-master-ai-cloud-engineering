# PDF-Creation Service

[![Creation-service](https://github.com/hingew/hsfl-master-ai-cloud-engineering/actions/workflows/creation-service.yml/badge.svg?branch=develop)](https://github.com/hingew/hsfl-master-ai-cloud-engineering/actions/workflows/creation-service.yml)
[![codecov](https://codecov.io/gh/hingew/hsfl-master-ai-cloud-engineering/graph/badge.svg?token=CDPMA4XLME&flag=creation-service)](https://codecov.io/gh/hingew/hsfl-master-ai-cloud-engineering)

## Description

The PDF Creation Service requests the template from the PDF Template Service and additionally receives dynamic fields from the client to incorporate them into the generated PDF. The PDF Creation Service has the advantage of operating without additional data storage, making it easily scalable. This design facilitates effortless scaling to meet increased demands

## Configuration
Configure the following environment variables:

```
PORT=<application port, default is 3000>
TEMPLATING_GRPC_ENDPOINT=<endpoint of the grpc templating server>
```



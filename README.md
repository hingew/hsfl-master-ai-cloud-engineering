# PDF-Designer

[![codecov](https://codecov.io/gh/hingew/hsfl-master-ai-cloud-engineering/graph/badge.svg?token=CDPMA4XLME)](https://codecov.io/gh/hingew/hsfl-master-ai-cloud-engineering)

## Description

The PDF Designer project allows users to generate PDFs based on a custom template. These templates can be dynamically populated with data, such as invoice generation or listing data in a table.

## Authors

Hauke Ingwersen\
hauke.ingwersen@stud.hs-flensburg.de\
Hochschule Flensburg

Robert Pfeiffer\
robert.pfeiffer@stud.hs-flensburg.de\
Hochschule Flensburg

Jannes Nebendahl\
jannes.nebendahl@stud.hs-flensburg.de\
Hochschule Flensburg

## Installation

Please copy the .env.example file to .env and adjust the AUTH_SIGN_KEY=.

`cp .env.example .env`

Create a private key for signing JWT tokens:

``` 
openssl genrsa -out key.pem 2048
openssl rsa -in key.pem -outform PEM -pubout -out public.pem
```



Copy the contents of the generated user-service.pem file and paste it into the AUTH_SIGN_KEY environment variable or directly into the .env file located at the project root.

Update the content of your `.env` file in the project root.
```
AUTH_PRIVATE_KEY="<your-private-key>"
AUTH_PUBLIC_KEY="<your-public-key>"
```

To build the docker container you have to run the following command:

`docker-compose up -d --build`

### Installation with MacOS ARM

Add the following line under _web_ in the _docker-compose.yml_ file only if MacOS with ARM is being used:
`platform: linux/amd64`

## Start Application

You can host the application in the following 3 environments. To start them you can use the script files in the script directory.

### Docker

Run the script `script\start_docker`

### Kubernetes with Minicube

Run the script `script\start_minicube`

### Kubernetes on own Cluster

Run the script `script\start_own_cluster`

## Current Architecture

![currentArchitecture](currentArchitecture.png)

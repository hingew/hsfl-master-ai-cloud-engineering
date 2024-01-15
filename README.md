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

```cp .env.example .env```

To create a private key for signing JWT tokens, use the following command:

```ssh-keygen -t ecdsa -f user-service.pem -m pem```

Copy the contents of the generated user-service.pem file and paste it into the AUTH_SIGN_KEY environment variable or directly into the .env file located at the project root.

Step-by-step instructions for starting all services with Docker.

```docker-compose up```

### Installation with MacOS ARM
Add the following line under *web* in the *docker-compose.yml* file only if MacOS with ARM is being used:
```platform: linux/amd64```



## Current Architecture
![currentArchitecture](currentArchitecture.png)


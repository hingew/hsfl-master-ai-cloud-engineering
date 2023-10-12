module github.com/hingew/hsfl-master-ai-cloud-engineering/user-service

go 1.21.1

require (
	github.com/hingew/hsfl-master-ai-cloud-engineering/lib v0.0.0
	github.com/lib/pq v1.10.9
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	golang.org/x/crypto v0.14.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/hingew/hsfl-master-ai-cloud-engineering/lib v0.0.0 => ../../lib/

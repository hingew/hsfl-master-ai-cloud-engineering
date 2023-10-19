module github.com/hingew/hsfl-master-ai-cloud-engineering/creation-service

go 1.21.1

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/hingew/hsfl-master-ai-cloud-engineering/lib v0.0.0
	github.com/lib/pq v1.10.9
	golang.org/x/crypto v0.14.0
)

require github.com/jung-kurt/gofpdf v1.16.2 // indirect

replace github.com/hingew/hsfl-master-ai-cloud-engineering/lib v0.0.0 => ../../lib/

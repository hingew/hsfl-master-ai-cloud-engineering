module github.com/hingew/hsfl-master-ai-cloud-engineering/creation-service

go 1.21.1

require github.com/hingew/hsfl-master-ai-cloud-engineering/lib v0.0.0

require (
	github.com/jung-kurt/gofpdf v1.16.2
	github.com/stretchr/testify v1.8.4
	go.uber.org/mock v0.3.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
)

replace github.com/hingew/hsfl-master-ai-cloud-engineering/lib v0.0.0 => ../../lib/

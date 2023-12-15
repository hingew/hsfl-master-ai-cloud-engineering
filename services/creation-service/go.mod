module github.com/hingew/hsfl-master-ai-cloud-engineering/creation-service

go 1.21.1

require github.com/hingew/hsfl-master-ai-cloud-engineering/lib v0.0.0

require (
	github.com/jung-kurt/gofpdf v1.16.2
	github.com/stretchr/testify v1.8.4
	go.uber.org/mock v0.3.0
	google.golang.org/grpc v1.57.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.15.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230525234030-28d5490b6b19 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/hingew/hsfl-master-ai-cloud-engineering/lib v0.0.0 => ../../lib/

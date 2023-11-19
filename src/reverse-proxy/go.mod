module github.com/akatranlp/hsfl-master-ai-cloud-engineering/reverse-proxy

go 1.21.3

require (
	github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib v0.0.0
	github.com/caarlos0/env/v10 v10.0.0
	github.com/joho/godotenv v1.5.1
	github.com/stretchr/testify v1.8.4
	go.uber.org/mock v0.3.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
)

replace github.com/akatranlp/hsfl-master-ai-cloud-engineering/lib => ../../lib

default: fmt lint build generate_docs

dev:
	goreleaser build --id $(shell go env GOOS) --single-target --snapshot --clean

snapshot:
	goreleaser release --snapshot --clean

build: dev

lint:
	golangci-lint run

generate_docs:
	go tool tfplugindocs

fmt:
	gofmt -s -w -e .

test:
	go test -v -cover -timeout=120s -parallel=10 ./...

docker_start:
	(cd fixtures; docker compose up -d)
	sleep 20

docker_stop:
	(cd fixtures; docker compose stop)

testacc:
	ICINGA2_API_PASSWORD="icingaweb" ICINGA2_API_URL="https://127.0.0.1:5665/v1" ICINGA2_API_USER=icingaweb ICINGA2_INSECURE_SKIP_TLS_VERIFY=true TF_ACC=1 go test -v -cover -timeout 120m ./...

.PHONY: dev snapshot build lint generate_docs fmt test docker_start docker_stop testacc

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

docker_start:
	(cd fixtures; docker compose -p icinga2-provider up -d)
	sleep 20

docker_stop:
	(cd fixtures; docker compose -p icinga2-provider stop)

testacc:
	ICINGA2_API_PASSWORD="icingaweb" ICINGA2_API_URL="https://127.0.0.1:5665/v1" ICINGA2_API_USER=icingaweb ICINGA2_INSECURE_SKIP_TLS_VERIFY=true ICINGA2_TRIES=3 ICINGA2_RETRY_DELAY=10 TF_ACC=1 go test -v -cover -timeout 120m ./...

acceptance: docker_start testacc

.PHONY: acceptance dev snapshot build lint generate_docs fmt docker_start docker_stop testacc

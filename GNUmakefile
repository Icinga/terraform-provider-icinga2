TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=icinga2

default: build

docker_start:
	docker run -d --name icinga2 -p 8080:80 -p 8443:443 -p 5665:5665 -it jordan/icinga2:2.15.1
	sleep 20

docker_get_root_password:
	$(eval password:=$(shell docker exec icinga2 bash -c 'grep password /etc/icinga2/conf.d/api-users.conf' | awk -F'"' '{ print $$2}' ))
	echo $(password)

docker_clean:
	docker stop icinga2
	docker rm icinga2

build: fmtcheck
	go install

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	$(eval password:=$(shell docker exec icinga2 bash -c 'grep password /etc/icinga2/conf.d/api-users.conf' | awk -F'"' '{ print $$2}' ))
	ICINGA2_API_PASSWORD="$(password)" ICINGA2_API_URL="https://127.0.0.1:5665/v1" ICINGA2_API_USER=root ICINGA2_INSECURE_SKIP_TLS_VERIFY=true TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"


test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

.PHONY: build test testacc vet fmt fmtcheck errcheck test-compile


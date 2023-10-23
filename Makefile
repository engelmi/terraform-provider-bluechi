
build:
	go build -o terraform-provider-bluechi

install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/bluechi/bluechi/1.0.0/linux_amd64
	mv terraform-provider-bluechi plugins/github.com/bluechi/bluechi/0.0.0/linux_amd64/

uninstall:
	rm ~/.terraform.d/plugins/registry.terraform.io/bluechi/bluechi/1.0.0/linux_amd64/terraform-provider-bluechi

test:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m


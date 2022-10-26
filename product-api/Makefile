check_install:
	which swagger || GO111MODULE=off go install -u github.com/go-swagger/go-swagger/cmd/swagger

install_swagger:
	go install -u github.com/go-swagger/go-swagger/cmd/swagger

swagger: check_install
	swagger generate spec -o ./swagger.yaml --scan-models
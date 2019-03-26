module github.com/tetrateio/training/samples/modernbank/tools/trafficGen

require (
	github.com/dustinkirkland/golang-petname v0.0.0-20170921220637-d3c2ba80e75e
	github.com/go-openapi/analysis v0.19.0 // indirect
	github.com/go-openapi/errors v0.19.0 // indirect
	github.com/go-openapi/jsonpointer v0.19.0 // indirect
	github.com/go-openapi/jsonreference v0.19.0 // indirect
	github.com/go-openapi/loads v0.19.0 // indirect
	github.com/go-openapi/runtime v0.19.0 // indirect
	github.com/go-openapi/spec v0.19.0 // indirect
	github.com/go-openapi/strfmt v0.19.0 // indirect
	github.com/go-openapi/swag v0.19.0
	github.com/go-openapi/validate v0.19.0 // indirect
	github.com/kr/pty v1.1.4 // indirect
	github.com/tetrateio/training/samples/modernbank/microservices/account v0.0.0
	github.com/tetrateio/training/samples/modernbank/microservices/transaction v0.0.0-20190325224548-63e051025289
	github.com/tetrateio/training/samples/modernbank/microservices/transaction-log v0.0.0
	github.com/tetrateio/training/samples/modernbank/microservices/user v0.0.0
	github.com/tidwall/pretty v0.0.0-20190325153808-1166b9ac2b65 // indirect
	golang.org/x/crypto v0.0.0-20190325154230-a5d413f7728c // indirect
	golang.org/x/net v0.0.0-20190324223953-e3b2ff56ed87 // indirect
	golang.org/x/sys v0.0.0-20190322080309-f49334f85ddc // indirect
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4
)

replace github.com/tetrateio/training/samples/modernbank/microservices/user => ../../microservices/user

replace github.com/tetrateio/training/samples/modernbank/microservices/account => ../../microservices/account

replace github.com/tetrateio/training/samples/modernbank/microservices/transaction-log => ../../microservices/transaction-log

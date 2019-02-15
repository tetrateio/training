module github.com/tetrateio/training/samples/modernbank/microservices/transaction

require (
	github.com/go-openapi/errors v0.18.0
	github.com/go-openapi/loads v0.18.0
	github.com/go-openapi/runtime v0.18.0
	github.com/go-openapi/spec v0.18.0
	github.com/go-openapi/strfmt v0.18.0
	github.com/go-openapi/swag v0.18.0
	github.com/go-openapi/validate v0.18.0
	github.com/jessevdk/go-flags v1.4.0
	github.com/tetrateio/training/samples/modernbank/microservices/account v0.0.0
	github.com/tetrateio/training/samples/modernbank/microservices/transaction-log v0.0.0
	golang.org/x/net v0.0.0-20190213061140-3a22650c66bd
)

replace github.com/tetrateio/training/samples/modernbank/microservices/account => ../account

replace github.com/tetrateio/training/samples/modernbank/microservices/transaction-log => ../transaction-log

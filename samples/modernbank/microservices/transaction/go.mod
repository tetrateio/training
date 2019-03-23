module github.com/tetrateio/training/samples/modernbank/microservices/transaction

require (
	github.com/go-openapi/errors v0.18.0
	github.com/go-openapi/loads v0.18.0
	github.com/go-openapi/runtime v0.18.0
	github.com/go-openapi/spec v0.18.0
	github.com/go-openapi/strfmt v0.18.0
	github.com/go-openapi/swag v0.18.0
	github.com/go-openapi/validate v0.18.0
	github.com/spf13/pflag v1.0.3
	github.com/tetrateio/training/samples/modernbank/microservices/account v0.0.0
	github.com/tetrateio/training/samples/modernbank/microservices/transaction-log v0.0.0
	golang.org/x/net v0.0.0-20190322120337-addf6b3196f6
)

replace github.com/tetrateio/training/samples/modernbank/microservices/account => ../account

replace github.com/tetrateio/training/samples/modernbank/microservices/transaction-log => ../transaction-log

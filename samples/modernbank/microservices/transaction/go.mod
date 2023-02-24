module github.com/tetrateio/training/samples/modernbank/microservices/transaction

require (
	github.com/go-openapi/errors v0.19.0
	github.com/go-openapi/loads v0.19.0
	github.com/go-openapi/runtime v0.19.0
	github.com/go-openapi/spec v0.19.0
	github.com/go-openapi/strfmt v0.19.0
	github.com/go-openapi/swag v0.19.0
	github.com/go-openapi/validate v0.19.0
	github.com/spf13/pflag v1.0.3
	github.com/tetrateio/training/samples/modernbank/microservices/account v0.0.0
	github.com/tetrateio/training/samples/modernbank/microservices/transaction-log v0.0.0
	golang.org/x/net v0.7.0
)

require (
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/asaskevich/govalidator v0.0.0-20180720115003-f9ffefc3facf // indirect
	github.com/docker/go-units v0.3.3 // indirect
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8 // indirect
	github.com/go-openapi/analysis v0.19.0 // indirect
	github.com/go-openapi/jsonpointer v0.19.0 // indirect
	github.com/go-openapi/jsonreference v0.19.0 // indirect
	github.com/mailru/easyjson v0.0.0-20190312143242-1de009706dbe // indirect
	github.com/mitchellh/mapstructure v1.1.2 // indirect
	golang.org/x/text v0.7.0 // indirect
	gopkg.in/yaml.v2 v2.2.2 // indirect
)

replace github.com/tetrateio/training/samples/modernbank/microservices/account => ../account

replace github.com/tetrateio/training/samples/modernbank/microservices/transaction-log => ../transaction-log

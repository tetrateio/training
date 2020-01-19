module github.com/tetrateio/training/samples/modernbank/tools/trafficGen

replace github.com/tetrateio/training/samples/modernbank/microservices/user => ../../microservices/user

replace github.com/tetrateio/training/samples/modernbank/microservices/account => ../../microservices/account

replace github.com/tetrateio/training/samples/modernbank/microservices/transaction-log => ../../microservices/transaction-log

go 1.13

require (
	github.com/dustinkirkland/golang-petname v0.0.0-20191129215211-8e5a1ed0cff0
	github.com/go-openapi/swag v0.19.6
	github.com/tetrateio/training/samples/modernbank/microservices/transaction v0.0.0-20200119150454-ca57c573ddde
	github.com/tetrateio/training/samples/modernbank/microservices/transaction-log v0.0.0
	github.com/tetrateio/training/samples/modernbank/microservices/user v0.0.0-00010101000000-000000000000
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
)

// This file is safe to edit. Once it exists it will not be overwritten

package serve

import (
	"crypto/tls"
	"log"
	"net/http"

	flag "github.com/spf13/pflag"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	accountsClient "github.com/tetrateio/training/samples/modernbank/microservices/account/pkg/client"
	accountsClientResources "github.com/tetrateio/training/samples/modernbank/microservices/account/pkg/client/accounts"
	translogClient "github.com/tetrateio/training/samples/modernbank/microservices/transaction-log/pkg/client"
	translogClientResources "github.com/tetrateio/training/samples/modernbank/microservices/transaction-log/pkg/client/transactions"
	translogModel "github.com/tetrateio/training/samples/modernbank/microservices/transaction-log/pkg/model"
	"github.com/tetrateio/training/samples/modernbank/microservices/transaction/pkg/model"
	"github.com/tetrateio/training/samples/modernbank/microservices/transaction/pkg/serve/restapi"
	"github.com/tetrateio/training/samples/modernbank/microservices/transaction/pkg/serve/restapi/health"
	"github.com/tetrateio/training/samples/modernbank/microservices/transaction/pkg/serve/restapi/transactions"
)

//go:generate swagger generate server --target ../../../transaction --name Transaction --spec ../../../../scripts/flat/transaction.yaml --api-package restapi --model-package pkg/model --server-package pkg/serve
var version *string = flag.String("version", "v1", "the version of service to run. Should match version label used in Istio.")

var (
	translog *translogClientResources.Client
	accounts *accountsClientResources.Client
)

func configureFlags(api *restapi.TransactionAPI) {
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		swag.CommandLineOptionsGroup{Options: version},
	}
}

func configureAPI(api *restapi.TransactionAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError
	api.Logger = log.Printf

	translog = translogClient.Default.Transactions
	accounts = accountsClient.Default.Accounts

	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()

	// TODO: handle partial failure caused by one of these requests failing
	// Ultimately, synchronous communication is the wrong paradigm here but this is a demo app!
	api.TransactionsCreateTransactionHandler = transactions.CreateTransactionHandlerFunc(func(params transactions.CreateTransactionParams) middleware.Responder {
		// Can't send negative monies!
		if *params.Body.Amount < 0 {
			api.Logger("Receveived transaction for negative amount")
			return transactions.NewCreateTransactionBadRequest().WithVersion(*version)
		}

		// Move the monies
		// TODO: Verify both accounts exist before moving the money around
		sendingParams := accountsClientResources.NewChangeBalanceParams().WithNumber(*params.Body.Sender).WithDelta(*params.Body.Amount * -1)
		propagateTracingHeadersAccounts(params, sendingParams)
		if _, err := accounts.ChangeBalance(sendingParams); err != nil {
			api.Logger("Error updating sender balance: %v", err)
			return transactions.NewCreateTransactionInternalServerError().WithVersion(*version)
		}
		receivingParams := accountsClientResources.NewChangeBalanceParams().WithNumber(*params.Body.Receiver).WithDelta(*params.Body.Amount)
		propagateTracingHeadersAccounts(params, receivingParams)
		if _, err := accounts.ChangeBalance(receivingParams); err != nil {
			api.Logger("Error updating receiver balance: %v", err)
			return transactions.NewCreateTransactionInternalServerError().WithVersion(*version)
		}

		// Add to transaction-log
		translogParams := translogClientResources.NewCreateTransactionParams().WithBody(transToTransLogNewTransaction(params.Body))
		propagateTracingHeadersTransactionLog(params, translogParams)
		res, err := translog.CreateTransaction(translogParams)
		if err != nil {
			log.Printf("failed to create transaction in log from %v to %v for %v: %v", *params.Body.Sender, *params.Body.Receiver, *params.Body.Amount, err)
			return transactions.NewCreateTransactionInternalServerError().WithVersion(*version)
		}
		return transactions.NewCreateTransactionCreated().WithPayload(transLogToTransTransaction(res.Payload)).WithVersion(*version)
	})

	api.HealthHealthCheckHandler = health.HealthCheckHandlerFunc(func(_ health.HealthCheckParams) middleware.Responder {
		return health.NewHealthCheckOK().WithVersion(*version)
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

func propagateTracingHeadersAccounts(inParams transactions.CreateTransactionParams, outParams *accountsClientResources.ChangeBalanceParams) {
	outParams.SetB3(inParams.B3)
	outParams.SetXB3Flags(inParams.XB3Flags)
	outParams.SetXB3Parentspanid(inParams.XB3Parentspanid)
	outParams.SetXB3Sampled(inParams.XB3Sampled)
	outParams.SetXB3SpanID(inParams.XB3SpanID)
	outParams.SetXB3Traceid(inParams.XB3Traceid)
	outParams.SetXRequestID(inParams.XRequestID)
}

// Something, something, generics...
func propagateTracingHeadersTransactionLog(inParams transactions.CreateTransactionParams, outParams *translogClientResources.CreateTransactionParams) {
	outParams.SetB3(inParams.B3)
	outParams.SetXB3Flags(inParams.XB3Flags)
	outParams.SetXB3Parentspanid(inParams.XB3Parentspanid)
	outParams.SetXB3Sampled(inParams.XB3Sampled)
	outParams.SetXB3SpanID(inParams.XB3SpanID)
	outParams.SetXB3Traceid(inParams.XB3Traceid)
	outParams.SetXRequestID(inParams.XRequestID)
}

func transToTransLogNewTransaction(trans *model.Newtransaction) *translogModel.Newtransaction {
	return &translogModel.Newtransaction{
		Amount:   trans.Amount,
		Sender:   trans.Sender,
		Receiver: trans.Receiver,
	}
}

func transLogToTransTransaction(translog *translogModel.Transaction) *model.Transaction {
	return &model.Transaction{
		ID:        translog.ID,
		Amount:    translog.Amount,
		Sender:    translog.Sender,
		Receiver:  translog.Receiver,
		Timestamp: translog.Timestamp,
	}
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}

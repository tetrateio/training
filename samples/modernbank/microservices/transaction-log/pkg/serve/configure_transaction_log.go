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

	"github.com/tetrateio/training/samples/modernbank/microservices/transaction-log/pkg/serve/restapi"
	"github.com/tetrateio/training/samples/modernbank/microservices/transaction-log/pkg/serve/restapi/health"
	"github.com/tetrateio/training/samples/modernbank/microservices/transaction-log/pkg/serve/restapi/transactions"
	"github.com/tetrateio/training/samples/modernbank/microservices/transaction-log/pkg/store"
	"github.com/tetrateio/training/samples/modernbank/microservices/transaction-log/pkg/store/mongodb"
)

//go:generate swagger generate server --target ../../../transaction-log --name TransactionLog --spec ../../../../scripts/flat/transaction-log.yaml --api-package restapi --model-package pkg/model --server-package pkg/serve
var version *string = flag.String("version", "v1", "the version of service to run. Should match version label used in Istio.")

var transactionStore store.Interface

func configureFlags(api *restapi.TransactionLogAPI) {
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		swag.CommandLineOptionsGroup{Options: version},
	}
}

func configureAPI(api *restapi.TransactionLogAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError
	api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()

	transactionStore = mongodb.NewMongoDB()

	api.TransactionsCreateTransactionHandler = transactions.CreateTransactionHandlerFunc(func(params transactions.CreateTransactionParams) middleware.Responder {
		res, err := transactionStore.Create(params.Body)
		if err != nil {
			api.Logger("Unable to create transaction: %v", err)
			return transactions.NewCreateTransactionInternalServerError().WithVersion(*version)
		}
		api.Logger("Created transaction %q", *res.ID)
		return transactions.NewCreateTransactionCreated().WithPayload(res).WithVersion(*version)
	})
	api.TransactionsGetTransactionReceivedHandler = transactions.GetTransactionReceivedHandlerFunc(func(params transactions.GetTransactionReceivedParams) middleware.Responder {
		res, err := transactionStore.GetReceived(params.Receiver, params.Transaction)
		if err != nil {
			if _, ok := err.(*store.NotFound); ok {
				return transactions.NewGetTransactionReceivedNotFound().WithVersion(*version)
			}
			api.Logger("Unable to retrieve transaction %v received by %v: %v", params.Transaction, params.Receiver, err)
			return transactions.NewGetTransactionReceivedInternalServerError().WithVersion(*version)
		}
		api.Logger("Retrieved transaction %v received by %v", params.Transaction, params.Receiver)
		return transactions.NewGetTransactionReceivedOK().WithPayload(res).WithVersion(*version)
	})
	api.TransactionsGetTransactionSentHandler = transactions.GetTransactionSentHandlerFunc(func(params transactions.GetTransactionSentParams) middleware.Responder {
		res, err := transactionStore.GetSent(params.Sender, params.Transaction)
		if err != nil {
			if _, ok := err.(*store.NotFound); ok {
				return transactions.NewGetTransactionSentNotFound().WithVersion(*version)
			}
			api.Logger("Unable to retrieve transaction %v sent by %v: %v", params.Transaction, params.Sender, err)
			return transactions.NewGetTransactionSentInternalServerError().WithVersion(*version)
		}
		api.Logger("Retrieved transaction %v sent by %v", params.Transaction, params.Sender)
		return transactions.NewGetTransactionSentOK().WithPayload(res).WithVersion(*version)
	})
	api.TransactionsListTransactionsReceivedHandler = transactions.ListTransactionsReceivedHandlerFunc(func(params transactions.ListTransactionsReceivedParams) middleware.Responder {
		res, err := transactionStore.ListReceived(params.Receiver)
		if err != nil {
			if _, ok := err.(*store.NotFound); ok {
				return transactions.NewListTransactionsReceivedNotFound().WithVersion(*version)
			}
			api.Logger("Unable to list transactions received by %v: %v", params.Receiver, err)
			return transactions.NewListTransactionsReceivedInternalServerError().WithVersion(*version)
		}
		api.Logger("Retrieved all transactions received by %v", params.Receiver)
		return transactions.NewListTransactionsReceivedOK().WithPayload(res).WithVersion(*version)
	})
	api.TransactionsListTransactionsSentHandler = transactions.ListTransactionsSentHandlerFunc(func(params transactions.ListTransactionsSentParams) middleware.Responder {
		res, err := transactionStore.ListSent(params.Sender)
		if err != nil {
			if _, ok := err.(*store.NotFound); ok {
				return transactions.NewListTransactionsSentNotFound().WithVersion(*version)
			}
			api.Logger("Unable to list transactions sent by %v: %v", params.Sender, err)
			return transactions.NewListTransactionsSentInternalServerError().WithVersion(*version)
		}
		api.Logger("Retrieved all transactions sent by %v", params.Sender)
		return transactions.NewListTransactionsSentOK().WithPayload(res).WithVersion(*version)
	})
	api.HealthHealthCheckHandler = health.HealthCheckHandlerFunc(func(_ health.HealthCheckParams) middleware.Responder {
		return health.NewHealthCheckOK().WithVersion(*version)
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
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

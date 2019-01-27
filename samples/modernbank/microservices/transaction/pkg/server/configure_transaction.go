// This file is safe to edit. Once it exists it will not be overwritten

package server

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/tetrateio/training/samples/modernbank/microservices/transaction/pkg/server/restapi"
	"github.com/tetrateio/training/samples/modernbank/microservices/transaction/pkg/server/restapi/transactions"
)

//go:generate swagger generate server --target ../../../transaction --name Transaction --spec ../../../../swagger/transaction.yaml --api-package restapi --model-package pkg/model --server-package pkg/server

func configureFlags(api *restapi.TransactionAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *restapi.TransactionAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.TransactionsCreateTransactionHandler = transactions.CreateTransactionHandlerFunc(func(params transactions.CreateTransactionParams) middleware.Responder {
		return middleware.NotImplemented("operation transactions.CreateTransaction has not yet been implemented")
	})
	api.TransactionsGetTransactionReceivedHandler = transactions.GetTransactionReceivedHandlerFunc(func(params transactions.GetTransactionReceivedParams) middleware.Responder {
		return middleware.NotImplemented("operation transactions.GetTransactionReceived has not yet been implemented")
	})
	api.TransactionsGetTransactionSentHandler = transactions.GetTransactionSentHandlerFunc(func(params transactions.GetTransactionSentParams) middleware.Responder {
		return middleware.NotImplemented("operation transactions.GetTransactionSent has not yet been implemented")
	})
	api.TransactionsListTransactionsReceivedHandler = transactions.ListTransactionsReceivedHandlerFunc(func(params transactions.ListTransactionsReceivedParams) middleware.Responder {
		return middleware.NotImplemented("operation transactions.ListTransactionsReceived has not yet been implemented")
	})
	api.TransactionsListTransactionsSentHandler = transactions.ListTransactionsSentHandlerFunc(func(params transactions.ListTransactionsSentParams) middleware.Responder {
		return middleware.NotImplemented("operation transactions.ListTransactionsSent has not yet been implemented")
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

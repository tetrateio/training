// This file is safe to edit. Once it exists it will not be overwritten

package server

import (
	"crypto/tls"
	"log"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	accountsClient "github.com/tetrateio/training/samples/modernbank/microservices/account/pkg/client"
	accountsClientResources "github.com/tetrateio/training/samples/modernbank/microservices/account/pkg/client/accounts"
	translogClient "github.com/tetrateio/training/samples/modernbank/microservices/transaction-log/pkg/client"
	translogClientResources "github.com/tetrateio/training/samples/modernbank/microservices/transaction-log/pkg/client/transactions"
	translogModel "github.com/tetrateio/training/samples/modernbank/microservices/transaction-log/pkg/model"
	"github.com/tetrateio/training/samples/modernbank/microservices/transaction/pkg/model"
	"github.com/tetrateio/training/samples/modernbank/microservices/transaction/pkg/server/restapi"
	"github.com/tetrateio/training/samples/modernbank/microservices/transaction/pkg/server/restapi/transactions"
)

//go:generate swagger generate server --target ../../../transaction --name Transaction --spec ../../../../scripts/flat/transaction.yaml --api-package restapi --model-package pkg/model --server-package pkg/server

var (
	translog *translogClientResources.Client
	accounts *accountsClientResources.Client
)

func configureFlags(api *restapi.TransactionAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *restapi.TransactionAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	translog = translogClient.Default.Transactions
	accounts = accountsClient.Default.Accounts

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()

	// TODO: handle partial failure caused by one of these requests failing
	// Ultimately, synchronous communication is the wrong method here but this is a demo app!
	api.TransactionsCreateTransactionHandler = transactions.CreateTransactionHandlerFunc(func(params transactions.CreateTransactionParams) middleware.Responder {
		// Can't send negative monies!
		if *params.Body.Amount < 0 {
			return transactions.NewCreateTransactionBadRequest()
		}

		// Move the monies
		// TODO: Verify both accounts exist before moving the money around
		sendingParams := accountsClientResources.NewChangeBalanceParams().WithNumber(*params.Body.Sender).WithDelta(*params.Body.Amount * -1)
		if _, err := accounts.ChangeBalance(sendingParams); err != nil {
			return transactions.NewCreateTransactionInternalServerError()
		}
		receivingParams := accountsClientResources.NewChangeBalanceParams().WithNumber(*params.Body.Receiver).WithDelta(*params.Body.Amount)
		if _, err := accounts.ChangeBalance(receivingParams); err != nil {
			return transactions.NewCreateTransactionInternalServerError()
		}

		// Add to transaction-log
		translogParams := translogClientResources.NewCreateTransactionParams().WithBody(transToTransLogNewTransaction(params.Body))
		res, err := translog.CreateTransaction(translogParams)
		if err != nil {
			log.Printf("failed to create transaction in log from %v to %v for %v: %v", *params.Body.Sender, *params.Body.Receiver, *params.Body.Amount, err)
			return transactions.NewCreateTransactionInternalServerError()
		}
		return transactions.NewCreateTransactionCreated().WithPayload(transLogToTransTransaction(res.Payload))
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
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
		ID:       translog.ID,
		Amount:   translog.Amount,
		Sender:   translog.Sender,
		Receiver: translog.Receiver,
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

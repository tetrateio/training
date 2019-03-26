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

	"github.com/tetrateio/training/samples/modernbank/microservices/account/pkg/serve/restapi"
	"github.com/tetrateio/training/samples/modernbank/microservices/account/pkg/serve/restapi/accounts"
	"github.com/tetrateio/training/samples/modernbank/microservices/account/pkg/serve/restapi/health"
	"github.com/tetrateio/training/samples/modernbank/microservices/account/pkg/store"
	"github.com/tetrateio/training/samples/modernbank/microservices/account/pkg/store/mongodb"
)

var accountStore store.Interface

//go:generate swagger generate server --target ../../../account --name Account --spec ../../../../swagger/account.yaml --api-package restapi --model-package pkg/model --server-package pkg/serve
var version *string = flag.String("version", "v1", "the version of service to run. Should match version label used in Istio.")

func configureFlags(api *restapi.AccountAPI) {
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		swag.CommandLineOptionsGroup{Options: version},
	}
}

func configureAPI(api *restapi.AccountAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError
	api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()

	accountStore = mongodb.NewMongoDB()

	api.AccountsCreateAccountHandler = accounts.CreateAccountHandlerFunc(func(params accounts.CreateAccountParams) middleware.Responder {
		res, err := accountStore.Create(params.HTTPRequest.Context(), params.Owner, params.Type)
		if err != nil {
			api.Logger("Error adding new account for %q to store: %v", params.Owner, err)
			return accounts.NewCreateAccountInternalServerError().WithVersion(*version)
		}
		api.Logger("Successfully created account for %q", params.Owner)
		return accounts.NewCreateAccountCreated().WithPayload(res).WithVersion(*version)
	})
	api.AccountsListAccountsHandler = accounts.ListAccountsHandlerFunc(func(params accounts.ListAccountsParams) middleware.Responder {
		res, err := accountStore.List(params.HTTPRequest.Context(), params.Owner)
		if err != nil {
			api.Logger("Error listing accounts for %q: %v", params.Owner, err)
			if _, ok := err.(*store.NotFound); ok {
				return accounts.NewListAccountsNotFound().WithVersion(*version)
			}
			return accounts.NewListAccountsInternalServerError().WithVersion(*version)
		}
		return accounts.NewListAccountsOK().WithPayload(res).WithVersion(*version)
	})
	api.AccountsChangeBalanceHandler = accounts.ChangeBalanceHandlerFunc(func(params accounts.ChangeBalanceParams) middleware.Responder {
		if err := accountStore.UpdateBalance(params.HTTPRequest.Context(), params.Number, params.Delta); err != nil {
			if _, ok := err.(*store.NotFound); ok {
				api.Logger("Account not found - %v", params.Number)
				return accounts.NewChangeBalanceNotFound().WithVersion(*version)
			}
			api.Logger("Error updating balance on account %v by amount %v: %V", params.Number, params.Delta, err)
			return accounts.NewChangeBalanceInternalServerError().WithVersion(*version)
		}
		api.Logger("Successfully updated balance on account %v by amount %v", params.Number, params.Delta)
		return accounts.NewChangeBalanceOK().WithVersion(*version)
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

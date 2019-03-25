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
	accountModel "github.com/tetrateio/training/samples/modernbank/microservices/account/pkg/model"

	"github.com/tetrateio/training/samples/modernbank/microservices/user/pkg/model"
	"github.com/tetrateio/training/samples/modernbank/microservices/user/pkg/serve/restapi"
	"github.com/tetrateio/training/samples/modernbank/microservices/user/pkg/serve/restapi/accounts"
	"github.com/tetrateio/training/samples/modernbank/microservices/user/pkg/serve/restapi/health"
	"github.com/tetrateio/training/samples/modernbank/microservices/user/pkg/serve/restapi/users"
	"github.com/tetrateio/training/samples/modernbank/microservices/user/pkg/store"
	"github.com/tetrateio/training/samples/modernbank/microservices/user/pkg/store/mongodb"
)

var userStore store.Interface

//go:generate swagger generate server --target ../../../user --name User --spec ../../../../swagger/user.yaml --api-package restapi --model-package pkg/model --server-package pkg/serve

var version *string = flag.String("version", "v1", "the version of service to run. Should match version label used in Istio.")

var (
	accountService *accountsClientResources.Client
)

func configureFlags(api *restapi.UserAPI) {
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		swag.CommandLineOptionsGroup{Options: version},
	}
}

func configureAPI(api *restapi.UserAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError
	api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()

	accountService = accountsClient.Default.Accounts
	userStore = mongodb.NewMongoDB()

	api.UsersCreateUserHandler = users.CreateUserHandlerFunc(func(params users.CreateUserParams) middleware.Responder {
		res, err := userStore.Create(params.Body)
		if err != nil {
			api.Logger("Error creating user %q: %v", *params.Body.Username, err)
			if _, ok := err.(*store.Conflict); ok {
				return users.NewCreateUserConflict().WithVersion(*version)
			}
			return users.NewCreateUserInternalServerError().WithVersion(*version)
		}
		api.Logger("Created user %q", *params.Body.Username)
		for i := 0; i < 10; i++ {
			newAccountParams := accountsClientResources.NewCreateAccountParams().WithOwner(*params.Body.Username)
			if i%2 == 0 {
				newAccountParams.SetType("saving")
			} else {
				newAccountParams.SetType("cash")
			}
			if _, err := accountService.CreateAccount(newAccountParams); err != nil {
				api.Logger("Error creating account for user %q: %v", *params.Body.Username, err)
			}
		}
		return users.NewCreateUserCreated().WithPayload(res).WithVersion(*version)
	})
	api.UsersGetUserByUserNameHandler = users.GetUserByUserNameHandlerFunc(func(params users.GetUserByUserNameParams) middleware.Responder {
		res, err := userStore.Get(params.Username)
		if err != nil {
			api.Logger("Error retrieving user %q: %v", params.Username, err)
			if _, ok := err.(*store.NotFound); ok {
				return users.NewGetUserByUserNameNotFound().WithVersion(*version)
			}
			return users.NewGetUserByUserNameInternalServerError().WithVersion(*version)
		}
		api.Logger("Retrieved user %q", params.Username)
		return users.NewGetUserByUserNameOK().WithPayload(res).WithVersion(*version)
	})
	api.AccountsListAccountsHandler = accounts.ListAccountsHandlerFunc(func(params accounts.ListAccountsParams) middleware.Responder {
		listParams := accountsClientResources.NewListAccountsParams().WithOwner(params.Username)
		res, err := accountService.ListAccounts(listParams)
		if err != nil {
			api.Logger("Error retrieving accounts for user %q: %v", params.Username, err)
			return accounts.NewListAccountsInternalServerError().WithVersion(*version)
		}
		return accounts.NewListAccountsOK().WithPayload(convertAccountList(res.Payload)).WithVersion(*version)
	})
	api.HealthHealthCheckHandler = health.HealthCheckHandlerFunc(func(_ health.HealthCheckParams) middleware.Responder {
		return health.NewHealthCheckOK().WithVersion(*version)
	})

	api.ServerShutdown = func() {}
	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// This is not nice...
func convertAccountList(in []*accountModel.Account) []*model.Account {
	res := []*model.Account{}
	for _, a := range in {
		res = append(res, &model.Account{
			Balance: a.Balance,
			Number:  a.Number,
			Owner:   a.Owner,
			Type:    a.Type,
		})
	}
	return res
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

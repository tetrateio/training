// This file is safe to edit. Once it exists it will not be overwritten

package server

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/tetrateio/training/samples/modernbank/microservices/user/pkg/server/restapi"
	"github.com/tetrateio/training/samples/modernbank/microservices/user/pkg/server/restapi/users"
	"github.com/tetrateio/training/samples/modernbank/microservices/user/pkg/store"
)

var userStore store.Interface

//go:generate swagger generate server --target ../../../user --name User --spec ../../../../swagger/user.yaml --api-package restapi --model-package pkg/model --server-package pkg/server

func configureFlags(api *restapi.UserAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *restapi.UserAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()

	userStore = store.NewInMemory()

	api.UsersCreateUserHandler = users.CreateUserHandlerFunc(func(params users.CreateUserParams) middleware.Responder {
		res, err := userStore.Create(params.Body)
		if err != nil {
			if _, ok := err.(*store.Conflict); ok {
				return users.NewCreateUserConflict()
			}
			return users.NewCreateUserInternalServerError()
		}
		return users.NewCreateUserCreated().WithPayload(res)
	})
	api.UsersDeleteUserHandler = users.DeleteUserHandlerFunc(func(params users.DeleteUserParams) middleware.Responder {
		if err := userStore.Delete(params.Username); err != nil {
			if _, ok := err.(*store.NotFound); ok {
				return users.NewDeleteUserNotFound()
			}
			return users.NewDeleteUserInternalServerError()
		}
		return users.NewDeleteUserOK()
	})
	api.UsersGetUserByUserNameHandler = users.GetUserByUserNameHandlerFunc(func(params users.GetUserByUserNameParams) middleware.Responder {
		res, err := userStore.Get(params.Username)
		if err != nil {
			if _, ok := err.(*store.NotFound); ok {
				return users.NewGetUserByUserNameNotFound()
			}
			return users.NewGetUserByUserNameInternalServerError()
		}
		return users.NewGetUserByUserNameOK().WithPayload(res)
	})
	api.UsersUpdateUserHandler = users.UpdateUserHandlerFunc(func(params users.UpdateUserParams) middleware.Responder {
		res, err := userStore.Update(params.Username, params.Body)
		if err != nil {
			if _, ok := err.(*store.NotFound); ok {
				return users.NewUpdateUserNotFound()
			}
			return users.NewUpdateUserInternalServerError()
		}
		return users.NewUpdateUserOK().WithPayload(res)
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

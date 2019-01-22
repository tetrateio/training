// This file is safe to edit. Once it exists it will not be overwritten

package rest

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/tetrateio/training/samples/modernbank/user/pkg/rest/api"
	"github.com/tetrateio/training/samples/modernbank/user/pkg/rest/api/user"
	"github.com/tetrateio/training/samples/modernbank/user/pkg/service"
)

//go:generate swagger generate server --target ../../../user --name User --spec ../../../swagger/user.yaml --api-package api --model-package pkg/model --server-package pkg/rest

var store service.Store

func configureFlags(api *api.UserAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *api.UserAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	store = service.NewInMemoryStore()

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.UserCreateUserHandler = user.CreateUserHandlerFunc(func(params user.CreateUserParams) middleware.Responder {
		if err := store.Create(params.Body); err != nil {
			if _, ok := err.(*service.Conflict); ok {
				return user.NewCreateUserConflict()
			}
			return user.NewCreateUserInternalServerError()
		}
		return user.NewCreateUserCreated()
	})

	api.UserGetUserByUserNameHandler = user.GetUserByUserNameHandlerFunc(func(params user.GetUserByUserNameParams) middleware.Responder {
		res, err := store.Get(params.Username)
		if err != nil {
			if _, ok := err.(*service.NotFound); ok {
				return user.NewGetUserByUserNameNotFound()
			}
			return user.NewGetUserByUserNameInternalServerError()
		}
		return user.NewGetUserByUserNameOK().WithPayload(res)
	})

	api.UserUpdateUserHandler = user.UpdateUserHandlerFunc(func(params user.UpdateUserParams) middleware.Responder {
		if err := store.Update(params.Username, params.Body); err != nil {
			if _, ok := err.(*service.NotFound); ok {
				return user.NewUpdateUserNotFound()
			}
			return user.NewUpdateUserInternalServerError()
		}
		return user.NewUpdateUserOK()
	})

	api.UserDeleteUserHandler = user.DeleteUserHandlerFunc(func(params user.DeleteUserParams) middleware.Responder {
		if err := store.Delete(params.Username); err != nil {
			if _, ok := err.(*service.NotFound); ok {
				return user.NewDeleteUserNotFound()
			}
			return user.NewDeleteUserInternalServerError()
		}
		return user.NewDeleteUserOK()
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

// Code generated by go-swagger; DO NOT EDIT.

package accounts

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
)

// NewListAccountsParams creates a new ListAccountsParams object
// with the default values initialized.
func NewListAccountsParams() *ListAccountsParams {
	var ()
	return &ListAccountsParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewListAccountsParamsWithTimeout creates a new ListAccountsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewListAccountsParamsWithTimeout(timeout time.Duration) *ListAccountsParams {
	var ()
	return &ListAccountsParams{

		timeout: timeout,
	}
}

// NewListAccountsParamsWithContext creates a new ListAccountsParams object
// with the default values initialized, and the ability to set a context for a request
func NewListAccountsParamsWithContext(ctx context.Context) *ListAccountsParams {
	var ()
	return &ListAccountsParams{

		Context: ctx,
	}
}

// NewListAccountsParamsWithHTTPClient creates a new ListAccountsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewListAccountsParamsWithHTTPClient(client *http.Client) *ListAccountsParams {
	var ()
	return &ListAccountsParams{
		HTTPClient: client,
	}
}

/*ListAccountsParams contains all the parameters to send to the API endpoint
for the list accounts operation typically these are written to a http.Request
*/
type ListAccountsParams struct {

	/*Owner
	  Owner of the accounts

	*/
	Owner string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the list accounts params
func (o *ListAccountsParams) WithTimeout(timeout time.Duration) *ListAccountsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the list accounts params
func (o *ListAccountsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the list accounts params
func (o *ListAccountsParams) WithContext(ctx context.Context) *ListAccountsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the list accounts params
func (o *ListAccountsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the list accounts params
func (o *ListAccountsParams) WithHTTPClient(client *http.Client) *ListAccountsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the list accounts params
func (o *ListAccountsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithOwner adds the owner to the list accounts params
func (o *ListAccountsParams) WithOwner(owner string) *ListAccountsParams {
	o.SetOwner(owner)
	return o
}

// SetOwner adds the owner to the list accounts params
func (o *ListAccountsParams) SetOwner(owner string) {
	o.Owner = owner
}

// WriteToRequest writes these params to a swagger request
func (o *ListAccountsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// query param owner
	qrOwner := o.Owner
	qOwner := qrOwner
	if qOwner != "" {
		if err := r.SetQueryParam("owner", qOwner); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

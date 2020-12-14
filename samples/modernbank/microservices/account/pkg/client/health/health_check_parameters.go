// Code generated by go-swagger; DO NOT EDIT.

package health

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

// NewHealthCheckParams creates a new HealthCheckParams object
// with the default values initialized.
func NewHealthCheckParams() *HealthCheckParams {
	var ()
	return &HealthCheckParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewHealthCheckParamsWithTimeout creates a new HealthCheckParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewHealthCheckParamsWithTimeout(timeout time.Duration) *HealthCheckParams {
	var ()
	return &HealthCheckParams{

		timeout: timeout,
	}
}

// NewHealthCheckParamsWithContext creates a new HealthCheckParams object
// with the default values initialized, and the ability to set a context for a request
func NewHealthCheckParamsWithContext(ctx context.Context) *HealthCheckParams {
	var ()
	return &HealthCheckParams{

		Context: ctx,
	}
}

// NewHealthCheckParamsWithHTTPClient creates a new HealthCheckParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewHealthCheckParamsWithHTTPClient(client *http.Client) *HealthCheckParams {
	var ()
	return &HealthCheckParams{
		HTTPClient: client,
	}
}

/*HealthCheckParams contains all the parameters to send to the API endpoint
for the health check operation typically these are written to a http.Request
*/
type HealthCheckParams struct {

	/*B3*/
	B3 *string
	/*XB3Flags*/
	XB3Flags *string
	/*XB3Parentspanid*/
	XB3Parentspanid *string
	/*XB3Sampled*/
	XB3Sampled *string
	/*XB3SpanID*/
	XB3SpanID *string
	/*XB3Traceid*/
	XB3Traceid *string
	/*XRequestID*/
	XRequestID *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the health check params
func (o *HealthCheckParams) WithTimeout(timeout time.Duration) *HealthCheckParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the health check params
func (o *HealthCheckParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the health check params
func (o *HealthCheckParams) WithContext(ctx context.Context) *HealthCheckParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the health check params
func (o *HealthCheckParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the health check params
func (o *HealthCheckParams) WithHTTPClient(client *http.Client) *HealthCheckParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the health check params
func (o *HealthCheckParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithB3 adds the b3 to the health check params
func (o *HealthCheckParams) WithB3(b3 *string) *HealthCheckParams {
	o.SetB3(b3)
	return o
}

// SetB3 adds the b3 to the health check params
func (o *HealthCheckParams) SetB3(b3 *string) {
	o.B3 = b3
}

// WithXB3Flags adds the xB3Flags to the health check params
func (o *HealthCheckParams) WithXB3Flags(xB3Flags *string) *HealthCheckParams {
	o.SetXB3Flags(xB3Flags)
	return o
}

// SetXB3Flags adds the xB3Flags to the health check params
func (o *HealthCheckParams) SetXB3Flags(xB3Flags *string) {
	o.XB3Flags = xB3Flags
}

// WithXB3Parentspanid adds the xB3Parentspanid to the health check params
func (o *HealthCheckParams) WithXB3Parentspanid(xB3Parentspanid *string) *HealthCheckParams {
	o.SetXB3Parentspanid(xB3Parentspanid)
	return o
}

// SetXB3Parentspanid adds the xB3Parentspanid to the health check params
func (o *HealthCheckParams) SetXB3Parentspanid(xB3Parentspanid *string) {
	o.XB3Parentspanid = xB3Parentspanid
}

// WithXB3Sampled adds the xB3Sampled to the health check params
func (o *HealthCheckParams) WithXB3Sampled(xB3Sampled *string) *HealthCheckParams {
	o.SetXB3Sampled(xB3Sampled)
	return o
}

// SetXB3Sampled adds the xB3Sampled to the health check params
func (o *HealthCheckParams) SetXB3Sampled(xB3Sampled *string) {
	o.XB3Sampled = xB3Sampled
}

// WithXB3SpanID adds the xB3SpanID to the health check params
func (o *HealthCheckParams) WithXB3SpanID(xB3SpanID *string) *HealthCheckParams {
	o.SetXB3SpanID(xB3SpanID)
	return o
}

// SetXB3SpanID adds the xB3SpanId to the health check params
func (o *HealthCheckParams) SetXB3SpanID(xB3SpanID *string) {
	o.XB3SpanID = xB3SpanID
}

// WithXB3Traceid adds the xB3Traceid to the health check params
func (o *HealthCheckParams) WithXB3Traceid(xB3Traceid *string) *HealthCheckParams {
	o.SetXB3Traceid(xB3Traceid)
	return o
}

// SetXB3Traceid adds the xB3Traceid to the health check params
func (o *HealthCheckParams) SetXB3Traceid(xB3Traceid *string) {
	o.XB3Traceid = xB3Traceid
}

// WithXRequestID adds the xRequestID to the health check params
func (o *HealthCheckParams) WithXRequestID(xRequestID *string) *HealthCheckParams {
	o.SetXRequestID(xRequestID)
	return o
}

// SetXRequestID adds the xRequestId to the health check params
func (o *HealthCheckParams) SetXRequestID(xRequestID *string) {
	o.XRequestID = xRequestID
}

// WriteToRequest writes these params to a swagger request
func (o *HealthCheckParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.B3 != nil {

		// header param b3
		if err := r.SetHeaderParam("b3", *o.B3); err != nil {
			return err
		}

	}

	if o.XB3Flags != nil {

		// header param x-b3-flags
		if err := r.SetHeaderParam("x-b3-flags", *o.XB3Flags); err != nil {
			return err
		}

	}

	if o.XB3Parentspanid != nil {

		// header param x-b3-parentspanid
		if err := r.SetHeaderParam("x-b3-parentspanid", *o.XB3Parentspanid); err != nil {
			return err
		}

	}

	if o.XB3Sampled != nil {

		// header param x-b3-sampled
		if err := r.SetHeaderParam("x-b3-sampled", *o.XB3Sampled); err != nil {
			return err
		}

	}

	if o.XB3SpanID != nil {

		// header param x-b3-spanId
		if err := r.SetHeaderParam("x-b3-spanId", *o.XB3SpanID); err != nil {
			return err
		}

	}

	if o.XB3Traceid != nil {

		// header param x-b3-traceid
		if err := r.SetHeaderParam("x-b3-traceid", *o.XB3Traceid); err != nil {
			return err
		}

	}

	if o.XRequestID != nil {

		// header param x-request-id
		if err := r.SetHeaderParam("x-request-id", *o.XRequestID); err != nil {
			return err
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
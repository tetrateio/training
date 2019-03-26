// Code generated by go-swagger; DO NOT EDIT.

package users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"

	model "github.com/tetrateio/training/samples/modernbank/microservices/user/pkg/model"
)

// NewCreateUserParams creates a new CreateUserParams object
// no default values defined in spec.
func NewCreateUserParams() CreateUserParams {

	return CreateUserParams{}
}

// CreateUserParams contains all the bound params for the create user operation
// typically these are obtained from a http.Request
//
// swagger:parameters createUser
type CreateUserParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  In: header
	*/
	B3 *string
	/*Created user object
	  Required: true
	  In: body
	*/
	Body *model.User
	/*
	  In: header
	*/
	XB3Flags *string
	/*
	  In: header
	*/
	XB3Parentspanid *string
	/*
	  In: header
	*/
	XB3Sampled *string
	/*
	  In: header
	*/
	XB3SpanID *string
	/*
	  In: header
	*/
	XB3Traceid *string
	/*
	  In: header
	*/
	XRequestID *string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewCreateUserParams() beforehand.
func (o *CreateUserParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if err := o.bindB3(r.Header[http.CanonicalHeaderKey("b3")], true, route.Formats); err != nil {
		res = append(res, err)
	}

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body model.User
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("body", "body"))
			} else {
				res = append(res, errors.NewParseError("body", "body", "", err))
			}
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.Body = &body
			}
		}
	} else {
		res = append(res, errors.Required("body", "body"))
	}
	if err := o.bindXB3Flags(r.Header[http.CanonicalHeaderKey("x-b3-flags")], true, route.Formats); err != nil {
		res = append(res, err)
	}

	if err := o.bindXB3Parentspanid(r.Header[http.CanonicalHeaderKey("x-b3-parentspanid")], true, route.Formats); err != nil {
		res = append(res, err)
	}

	if err := o.bindXB3Sampled(r.Header[http.CanonicalHeaderKey("x-b3-sampled")], true, route.Formats); err != nil {
		res = append(res, err)
	}

	if err := o.bindXB3SpanID(r.Header[http.CanonicalHeaderKey("x-b3-spanId")], true, route.Formats); err != nil {
		res = append(res, err)
	}

	if err := o.bindXB3Traceid(r.Header[http.CanonicalHeaderKey("x-b3-traceid")], true, route.Formats); err != nil {
		res = append(res, err)
	}

	if err := o.bindXRequestID(r.Header[http.CanonicalHeaderKey("x-request-id")], true, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindB3 binds and validates parameter B3 from header.
func (o *CreateUserParams) bindB3(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.B3 = &raw

	return nil
}

// bindXB3Flags binds and validates parameter XB3Flags from header.
func (o *CreateUserParams) bindXB3Flags(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.XB3Flags = &raw

	return nil
}

// bindXB3Parentspanid binds and validates parameter XB3Parentspanid from header.
func (o *CreateUserParams) bindXB3Parentspanid(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.XB3Parentspanid = &raw

	return nil
}

// bindXB3Sampled binds and validates parameter XB3Sampled from header.
func (o *CreateUserParams) bindXB3Sampled(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.XB3Sampled = &raw

	return nil
}

// bindXB3SpanID binds and validates parameter XB3SpanID from header.
func (o *CreateUserParams) bindXB3SpanID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.XB3SpanID = &raw

	return nil
}

// bindXB3Traceid binds and validates parameter XB3Traceid from header.
func (o *CreateUserParams) bindXB3Traceid(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.XB3Traceid = &raw

	return nil
}

// bindXRequestID binds and validates parameter XRequestID from header.
func (o *CreateUserParams) bindXRequestID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.XRequestID = &raw

	return nil
}

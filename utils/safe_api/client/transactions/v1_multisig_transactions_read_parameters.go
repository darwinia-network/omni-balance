// Code generated by go-swagger; DO NOT EDIT.

package transactions

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewV1MultisigTransactionsReadParams creates a new V1MultisigTransactionsReadParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewV1MultisigTransactionsReadParams() *V1MultisigTransactionsReadParams {
	return &V1MultisigTransactionsReadParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewV1MultisigTransactionsReadParamsWithTimeout creates a new V1MultisigTransactionsReadParams object
// with the ability to set a timeout on a request.
func NewV1MultisigTransactionsReadParamsWithTimeout(timeout time.Duration) *V1MultisigTransactionsReadParams {
	return &V1MultisigTransactionsReadParams{
		timeout: timeout,
	}
}

// NewV1MultisigTransactionsReadParamsWithContext creates a new V1MultisigTransactionsReadParams object
// with the ability to set a context for a request.
func NewV1MultisigTransactionsReadParamsWithContext(ctx context.Context) *V1MultisigTransactionsReadParams {
	return &V1MultisigTransactionsReadParams{
		Context: ctx,
	}
}

// NewV1MultisigTransactionsReadParamsWithHTTPClient creates a new V1MultisigTransactionsReadParams object
// with the ability to set a custom HTTPClient for a request.
func NewV1MultisigTransactionsReadParamsWithHTTPClient(client *http.Client) *V1MultisigTransactionsReadParams {
	return &V1MultisigTransactionsReadParams{
		HTTPClient: client,
	}
}

/*
V1MultisigTransactionsReadParams contains all the parameters to send to the API endpoint

	for the v1 multisig transactions read operation.

	Typically these are written to a http.Request.
*/
type V1MultisigTransactionsReadParams struct {

	// SafeTxHash.
	SafeTxHash string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the v1 multisig transactions read params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *V1MultisigTransactionsReadParams) WithDefaults() *V1MultisigTransactionsReadParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the v1 multisig transactions read params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *V1MultisigTransactionsReadParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the v1 multisig transactions read params
func (o *V1MultisigTransactionsReadParams) WithTimeout(timeout time.Duration) *V1MultisigTransactionsReadParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the v1 multisig transactions read params
func (o *V1MultisigTransactionsReadParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the v1 multisig transactions read params
func (o *V1MultisigTransactionsReadParams) WithContext(ctx context.Context) *V1MultisigTransactionsReadParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the v1 multisig transactions read params
func (o *V1MultisigTransactionsReadParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the v1 multisig transactions read params
func (o *V1MultisigTransactionsReadParams) WithHTTPClient(client *http.Client) *V1MultisigTransactionsReadParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the v1 multisig transactions read params
func (o *V1MultisigTransactionsReadParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithSafeTxHash adds the safeTxHash to the v1 multisig transactions read params
func (o *V1MultisigTransactionsReadParams) WithSafeTxHash(safeTxHash string) *V1MultisigTransactionsReadParams {
	o.SetSafeTxHash(safeTxHash)
	return o
}

// SetSafeTxHash adds the safeTxHash to the v1 multisig transactions read params
func (o *V1MultisigTransactionsReadParams) SetSafeTxHash(safeTxHash string) {
	o.SafeTxHash = safeTxHash
}

// WriteToRequest writes these params to a swagger request
func (o *V1MultisigTransactionsReadParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param safe_tx_hash
	if err := r.SetPathParam("safe_tx_hash", o.SafeTxHash); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

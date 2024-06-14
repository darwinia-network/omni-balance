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

	models "omni-balance/utils/safe_api"
)

// NewV1SafesMultisigTransactionsEstimationsCreateParams creates a new V1SafesMultisigTransactionsEstimationsCreateParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewV1SafesMultisigTransactionsEstimationsCreateParams() *V1SafesMultisigTransactionsEstimationsCreateParams {
	return &V1SafesMultisigTransactionsEstimationsCreateParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewV1SafesMultisigTransactionsEstimationsCreateParamsWithTimeout creates a new V1SafesMultisigTransactionsEstimationsCreateParams object
// with the ability to set a timeout on a request.
func NewV1SafesMultisigTransactionsEstimationsCreateParamsWithTimeout(timeout time.Duration) *V1SafesMultisigTransactionsEstimationsCreateParams {
	return &V1SafesMultisigTransactionsEstimationsCreateParams{
		timeout: timeout,
	}
}

// NewV1SafesMultisigTransactionsEstimationsCreateParamsWithContext creates a new V1SafesMultisigTransactionsEstimationsCreateParams object
// with the ability to set a context for a request.
func NewV1SafesMultisigTransactionsEstimationsCreateParamsWithContext(ctx context.Context) *V1SafesMultisigTransactionsEstimationsCreateParams {
	return &V1SafesMultisigTransactionsEstimationsCreateParams{
		Context: ctx,
	}
}

// NewV1SafesMultisigTransactionsEstimationsCreateParamsWithHTTPClient creates a new V1SafesMultisigTransactionsEstimationsCreateParams object
// with the ability to set a custom HTTPClient for a request.
func NewV1SafesMultisigTransactionsEstimationsCreateParamsWithHTTPClient(client *http.Client) *V1SafesMultisigTransactionsEstimationsCreateParams {
	return &V1SafesMultisigTransactionsEstimationsCreateParams{
		HTTPClient: client,
	}
}

/*
V1SafesMultisigTransactionsEstimationsCreateParams contains all the parameters to send to the API endpoint

	for the v1 safes multisig transactions estimations create operation.

	Typically these are written to a http.Request.
*/
type V1SafesMultisigTransactionsEstimationsCreateParams struct {

	// Address.
	Address string

	// Data.
	Data *models.SafeMultisigTransactionEstimate

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the v1 safes multisig transactions estimations create params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *V1SafesMultisigTransactionsEstimationsCreateParams) WithDefaults() *V1SafesMultisigTransactionsEstimationsCreateParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the v1 safes multisig transactions estimations create params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *V1SafesMultisigTransactionsEstimationsCreateParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the v1 safes multisig transactions estimations create params
func (o *V1SafesMultisigTransactionsEstimationsCreateParams) WithTimeout(timeout time.Duration) *V1SafesMultisigTransactionsEstimationsCreateParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the v1 safes multisig transactions estimations create params
func (o *V1SafesMultisigTransactionsEstimationsCreateParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the v1 safes multisig transactions estimations create params
func (o *V1SafesMultisigTransactionsEstimationsCreateParams) WithContext(ctx context.Context) *V1SafesMultisigTransactionsEstimationsCreateParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the v1 safes multisig transactions estimations create params
func (o *V1SafesMultisigTransactionsEstimationsCreateParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the v1 safes multisig transactions estimations create params
func (o *V1SafesMultisigTransactionsEstimationsCreateParams) WithHTTPClient(client *http.Client) *V1SafesMultisigTransactionsEstimationsCreateParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the v1 safes multisig transactions estimations create params
func (o *V1SafesMultisigTransactionsEstimationsCreateParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithAddress adds the address to the v1 safes multisig transactions estimations create params
func (o *V1SafesMultisigTransactionsEstimationsCreateParams) WithAddress(address string) *V1SafesMultisigTransactionsEstimationsCreateParams {
	o.SetAddress(address)
	return o
}

// SetAddress adds the address to the v1 safes multisig transactions estimations create params
func (o *V1SafesMultisigTransactionsEstimationsCreateParams) SetAddress(address string) {
	o.Address = address
}

// WithData adds the data to the v1 safes multisig transactions estimations create params
func (o *V1SafesMultisigTransactionsEstimationsCreateParams) WithData(data *models.SafeMultisigTransactionEstimate) *V1SafesMultisigTransactionsEstimationsCreateParams {
	o.SetData(data)
	return o
}

// SetData adds the data to the v1 safes multisig transactions estimations create params
func (o *V1SafesMultisigTransactionsEstimationsCreateParams) SetData(data *models.SafeMultisigTransactionEstimate) {
	o.Data = data
}

// WriteToRequest writes these params to a swagger request
func (o *V1SafesMultisigTransactionsEstimationsCreateParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param address
	if err := r.SetPathParam("address", o.Address); err != nil {
		return err
	}
	if o.Data != nil {
		if err := r.SetBodyParam(o.Data); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

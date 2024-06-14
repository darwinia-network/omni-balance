// Code generated by go-swagger; DO NOT EDIT.

package safes

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

// NewV1SafesDelegatesDeleteParams creates a new V1SafesDelegatesDeleteParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewV1SafesDelegatesDeleteParams() *V1SafesDelegatesDeleteParams {
	return &V1SafesDelegatesDeleteParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewV1SafesDelegatesDeleteParamsWithTimeout creates a new V1SafesDelegatesDeleteParams object
// with the ability to set a timeout on a request.
func NewV1SafesDelegatesDeleteParamsWithTimeout(timeout time.Duration) *V1SafesDelegatesDeleteParams {
	return &V1SafesDelegatesDeleteParams{
		timeout: timeout,
	}
}

// NewV1SafesDelegatesDeleteParamsWithContext creates a new V1SafesDelegatesDeleteParams object
// with the ability to set a context for a request.
func NewV1SafesDelegatesDeleteParamsWithContext(ctx context.Context) *V1SafesDelegatesDeleteParams {
	return &V1SafesDelegatesDeleteParams{
		Context: ctx,
	}
}

// NewV1SafesDelegatesDeleteParamsWithHTTPClient creates a new V1SafesDelegatesDeleteParams object
// with the ability to set a custom HTTPClient for a request.
func NewV1SafesDelegatesDeleteParamsWithHTTPClient(client *http.Client) *V1SafesDelegatesDeleteParams {
	return &V1SafesDelegatesDeleteParams{
		HTTPClient: client,
	}
}

/*
V1SafesDelegatesDeleteParams contains all the parameters to send to the API endpoint

	for the v1 safes delegates delete operation.

	Typically these are written to a http.Request.
*/
type V1SafesDelegatesDeleteParams struct {

	// Address.
	Address string

	// Data.
	Data *models.SafeDelegateDelete

	// DelegateAddress.
	DelegateAddress string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the v1 safes delegates delete params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *V1SafesDelegatesDeleteParams) WithDefaults() *V1SafesDelegatesDeleteParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the v1 safes delegates delete params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *V1SafesDelegatesDeleteParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the v1 safes delegates delete params
func (o *V1SafesDelegatesDeleteParams) WithTimeout(timeout time.Duration) *V1SafesDelegatesDeleteParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the v1 safes delegates delete params
func (o *V1SafesDelegatesDeleteParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the v1 safes delegates delete params
func (o *V1SafesDelegatesDeleteParams) WithContext(ctx context.Context) *V1SafesDelegatesDeleteParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the v1 safes delegates delete params
func (o *V1SafesDelegatesDeleteParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the v1 safes delegates delete params
func (o *V1SafesDelegatesDeleteParams) WithHTTPClient(client *http.Client) *V1SafesDelegatesDeleteParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the v1 safes delegates delete params
func (o *V1SafesDelegatesDeleteParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithAddress adds the address to the v1 safes delegates delete params
func (o *V1SafesDelegatesDeleteParams) WithAddress(address string) *V1SafesDelegatesDeleteParams {
	o.SetAddress(address)
	return o
}

// SetAddress adds the address to the v1 safes delegates delete params
func (o *V1SafesDelegatesDeleteParams) SetAddress(address string) {
	o.Address = address
}

// WithData adds the data to the v1 safes delegates delete params
func (o *V1SafesDelegatesDeleteParams) WithData(data *models.SafeDelegateDelete) *V1SafesDelegatesDeleteParams {
	o.SetData(data)
	return o
}

// SetData adds the data to the v1 safes delegates delete params
func (o *V1SafesDelegatesDeleteParams) SetData(data *models.SafeDelegateDelete) {
	o.Data = data
}

// WithDelegateAddress adds the delegateAddress to the v1 safes delegates delete params
func (o *V1SafesDelegatesDeleteParams) WithDelegateAddress(delegateAddress string) *V1SafesDelegatesDeleteParams {
	o.SetDelegateAddress(delegateAddress)
	return o
}

// SetDelegateAddress adds the delegateAddress to the v1 safes delegates delete params
func (o *V1SafesDelegatesDeleteParams) SetDelegateAddress(delegateAddress string) {
	o.DelegateAddress = delegateAddress
}

// WriteToRequest writes these params to a swagger request
func (o *V1SafesDelegatesDeleteParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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

	// path param delegate_address
	if err := r.SetPathParam("delegate_address", o.DelegateAddress); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

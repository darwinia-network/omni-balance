// Code generated by go-swagger; DO NOT EDIT.

package data_decoder

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

// NewV1DataDecoderCreateParams creates a new V1DataDecoderCreateParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewV1DataDecoderCreateParams() *V1DataDecoderCreateParams {
	return &V1DataDecoderCreateParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewV1DataDecoderCreateParamsWithTimeout creates a new V1DataDecoderCreateParams object
// with the ability to set a timeout on a request.
func NewV1DataDecoderCreateParamsWithTimeout(timeout time.Duration) *V1DataDecoderCreateParams {
	return &V1DataDecoderCreateParams{
		timeout: timeout,
	}
}

// NewV1DataDecoderCreateParamsWithContext creates a new V1DataDecoderCreateParams object
// with the ability to set a context for a request.
func NewV1DataDecoderCreateParamsWithContext(ctx context.Context) *V1DataDecoderCreateParams {
	return &V1DataDecoderCreateParams{
		Context: ctx,
	}
}

// NewV1DataDecoderCreateParamsWithHTTPClient creates a new V1DataDecoderCreateParams object
// with the ability to set a custom HTTPClient for a request.
func NewV1DataDecoderCreateParamsWithHTTPClient(client *http.Client) *V1DataDecoderCreateParams {
	return &V1DataDecoderCreateParams{
		HTTPClient: client,
	}
}

/*
V1DataDecoderCreateParams contains all the parameters to send to the API endpoint

	for the v1 data decoder create operation.

	Typically these are written to a http.Request.
*/
type V1DataDecoderCreateParams struct {

	// Data.
	Data *models.DataDecoder

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the v1 data decoder create params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *V1DataDecoderCreateParams) WithDefaults() *V1DataDecoderCreateParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the v1 data decoder create params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *V1DataDecoderCreateParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the v1 data decoder create params
func (o *V1DataDecoderCreateParams) WithTimeout(timeout time.Duration) *V1DataDecoderCreateParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the v1 data decoder create params
func (o *V1DataDecoderCreateParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the v1 data decoder create params
func (o *V1DataDecoderCreateParams) WithContext(ctx context.Context) *V1DataDecoderCreateParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the v1 data decoder create params
func (o *V1DataDecoderCreateParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the v1 data decoder create params
func (o *V1DataDecoderCreateParams) WithHTTPClient(client *http.Client) *V1DataDecoderCreateParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the v1 data decoder create params
func (o *V1DataDecoderCreateParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithData adds the data to the v1 data decoder create params
func (o *V1DataDecoderCreateParams) WithData(data *models.DataDecoder) *V1DataDecoderCreateParams {
	o.SetData(data)
	return o
}

// SetData adds the data to the v1 data decoder create params
func (o *V1DataDecoderCreateParams) SetData(data *models.DataDecoder) {
	o.Data = data
}

// WriteToRequest writes these params to a swagger request
func (o *V1DataDecoderCreateParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
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

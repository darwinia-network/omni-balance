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
	"github.com/go-openapi/swag"
)

// NewV1SafesTransfersListParams creates a new V1SafesTransfersListParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewV1SafesTransfersListParams() *V1SafesTransfersListParams {
	return &V1SafesTransfersListParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewV1SafesTransfersListParamsWithTimeout creates a new V1SafesTransfersListParams object
// with the ability to set a timeout on a request.
func NewV1SafesTransfersListParamsWithTimeout(timeout time.Duration) *V1SafesTransfersListParams {
	return &V1SafesTransfersListParams{
		timeout: timeout,
	}
}

// NewV1SafesTransfersListParamsWithContext creates a new V1SafesTransfersListParams object
// with the ability to set a context for a request.
func NewV1SafesTransfersListParamsWithContext(ctx context.Context) *V1SafesTransfersListParams {
	return &V1SafesTransfersListParams{
		Context: ctx,
	}
}

// NewV1SafesTransfersListParamsWithHTTPClient creates a new V1SafesTransfersListParams object
// with the ability to set a custom HTTPClient for a request.
func NewV1SafesTransfersListParamsWithHTTPClient(client *http.Client) *V1SafesTransfersListParams {
	return &V1SafesTransfersListParams{
		HTTPClient: client,
	}
}

/*
V1SafesTransfersListParams contains all the parameters to send to the API endpoint

	for the v1 safes transfers list operation.

	Typically these are written to a http.Request.
*/
type V1SafesTransfersListParams struct {

	/* From.

	   _from
	*/
	From *string

	// Address.
	Address string

	/* BlockNumber.

	   block_number
	*/
	BlockNumber *string

	/* BlockNumberGt.

	   block_number__gt
	*/
	BlockNumberGt *string

	/* BlockNumberLt.

	   block_number__lt
	*/
	BlockNumberLt *string

	/* Erc20.

	   erc20
	*/
	Erc20 *string

	/* Erc721.

	   erc721
	*/
	Erc721 *string

	/* Ether.

	   ether
	*/
	Ether *string

	/* ExecutionDateGt.

	   execution_date__gt
	*/
	ExecutionDateGt *string

	/* ExecutionDateGte.

	   execution_date__gte
	*/
	ExecutionDateGte *string

	/* ExecutionDateLt.

	   execution_date__lt
	*/
	ExecutionDateLt *string

	/* ExecutionDateLte.

	   execution_date__lte
	*/
	ExecutionDateLte *string

	/* Limit.

	   Number of results to return per page.
	*/
	Limit *int64

	/* Offset.

	   The initial index from which to return the results.
	*/
	Offset *int64

	/* To.

	   to
	*/
	To *string

	/* TokenAddress.

	   token_address
	*/
	TokenAddress *string

	/* TransactionHash.

	   transaction_hash
	*/
	TransactionHash *string

	/* Value.

	   value
	*/
	Value *string

	/* ValueGt.

	   value__gt
	*/
	ValueGt *string

	/* ValueLt.

	   value__lt
	*/
	ValueLt *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the v1 safes transfers list params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *V1SafesTransfersListParams) WithDefaults() *V1SafesTransfersListParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the v1 safes transfers list params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *V1SafesTransfersListParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) WithTimeout(timeout time.Duration) *V1SafesTransfersListParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) WithContext(ctx context.Context) *V1SafesTransfersListParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) WithHTTPClient(client *http.Client) *V1SafesTransfersListParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithFrom adds the from to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) WithFrom(from *string) *V1SafesTransfersListParams {
	o.SetFrom(from)
	return o
}

// SetFrom adds the from to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) SetFrom(from *string) {
	o.From = from
}

// WithAddress adds the address to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) WithAddress(address string) *V1SafesTransfersListParams {
	o.SetAddress(address)
	return o
}

// SetAddress adds the address to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) SetAddress(address string) {
	o.Address = address
}

// WithBlockNumber adds the blockNumber to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) WithBlockNumber(blockNumber *string) *V1SafesTransfersListParams {
	o.SetBlockNumber(blockNumber)
	return o
}

// SetBlockNumber adds the blockNumber to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) SetBlockNumber(blockNumber *string) {
	o.BlockNumber = blockNumber
}

// WithBlockNumberGt adds the blockNumberGt to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) WithBlockNumberGt(blockNumberGt *string) *V1SafesTransfersListParams {
	o.SetBlockNumberGt(blockNumberGt)
	return o
}

// SetBlockNumberGt adds the blockNumberGt to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) SetBlockNumberGt(blockNumberGt *string) {
	o.BlockNumberGt = blockNumberGt
}

// WithBlockNumberLt adds the blockNumberLt to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) WithBlockNumberLt(blockNumberLt *string) *V1SafesTransfersListParams {
	o.SetBlockNumberLt(blockNumberLt)
	return o
}

// SetBlockNumberLt adds the blockNumberLt to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) SetBlockNumberLt(blockNumberLt *string) {
	o.BlockNumberLt = blockNumberLt
}

// WithErc20 adds the erc20 to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) WithErc20(erc20 *string) *V1SafesTransfersListParams {
	o.SetErc20(erc20)
	return o
}

// SetErc20 adds the erc20 to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) SetErc20(erc20 *string) {
	o.Erc20 = erc20
}

// WithErc721 adds the erc721 to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) WithErc721(erc721 *string) *V1SafesTransfersListParams {
	o.SetErc721(erc721)
	return o
}

// SetErc721 adds the erc721 to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) SetErc721(erc721 *string) {
	o.Erc721 = erc721
}

// WithEther adds the ether to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) WithEther(ether *string) *V1SafesTransfersListParams {
	o.SetEther(ether)
	return o
}

// SetEther adds the ether to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) SetEther(ether *string) {
	o.Ether = ether
}

// WithExecutionDateGt adds the executionDateGt to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) WithExecutionDateGt(executionDateGt *string) *V1SafesTransfersListParams {
	o.SetExecutionDateGt(executionDateGt)
	return o
}

// SetExecutionDateGt adds the executionDateGt to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) SetExecutionDateGt(executionDateGt *string) {
	o.ExecutionDateGt = executionDateGt
}

// WithExecutionDateGte adds the executionDateGte to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) WithExecutionDateGte(executionDateGte *string) *V1SafesTransfersListParams {
	o.SetExecutionDateGte(executionDateGte)
	return o
}

// SetExecutionDateGte adds the executionDateGte to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) SetExecutionDateGte(executionDateGte *string) {
	o.ExecutionDateGte = executionDateGte
}

// WithExecutionDateLt adds the executionDateLt to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) WithExecutionDateLt(executionDateLt *string) *V1SafesTransfersListParams {
	o.SetExecutionDateLt(executionDateLt)
	return o
}

// SetExecutionDateLt adds the executionDateLt to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) SetExecutionDateLt(executionDateLt *string) {
	o.ExecutionDateLt = executionDateLt
}

// WithExecutionDateLte adds the executionDateLte to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) WithExecutionDateLte(executionDateLte *string) *V1SafesTransfersListParams {
	o.SetExecutionDateLte(executionDateLte)
	return o
}

// SetExecutionDateLte adds the executionDateLte to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) SetExecutionDateLte(executionDateLte *string) {
	o.ExecutionDateLte = executionDateLte
}

// WithLimit adds the limit to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) WithLimit(limit *int64) *V1SafesTransfersListParams {
	o.SetLimit(limit)
	return o
}

// SetLimit adds the limit to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) SetLimit(limit *int64) {
	o.Limit = limit
}

// WithOffset adds the offset to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) WithOffset(offset *int64) *V1SafesTransfersListParams {
	o.SetOffset(offset)
	return o
}

// SetOffset adds the offset to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) SetOffset(offset *int64) {
	o.Offset = offset
}

// WithTo adds the to to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) WithTo(to *string) *V1SafesTransfersListParams {
	o.SetTo(to)
	return o
}

// SetTo adds the to to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) SetTo(to *string) {
	o.To = to
}

// WithTokenAddress adds the tokenAddress to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) WithTokenAddress(tokenAddress *string) *V1SafesTransfersListParams {
	o.SetTokenAddress(tokenAddress)
	return o
}

// SetTokenAddress adds the tokenAddress to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) SetTokenAddress(tokenAddress *string) {
	o.TokenAddress = tokenAddress
}

// WithTransactionHash adds the transactionHash to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) WithTransactionHash(transactionHash *string) *V1SafesTransfersListParams {
	o.SetTransactionHash(transactionHash)
	return o
}

// SetTransactionHash adds the transactionHash to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) SetTransactionHash(transactionHash *string) {
	o.TransactionHash = transactionHash
}

// WithValue adds the value to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) WithValue(value *string) *V1SafesTransfersListParams {
	o.SetValue(value)
	return o
}

// SetValue adds the value to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) SetValue(value *string) {
	o.Value = value
}

// WithValueGt adds the valueGt to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) WithValueGt(valueGt *string) *V1SafesTransfersListParams {
	o.SetValueGt(valueGt)
	return o
}

// SetValueGt adds the valueGt to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) SetValueGt(valueGt *string) {
	o.ValueGt = valueGt
}

// WithValueLt adds the valueLt to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) WithValueLt(valueLt *string) *V1SafesTransfersListParams {
	o.SetValueLt(valueLt)
	return o
}

// SetValueLt adds the valueLt to the v1 safes transfers list params
func (o *V1SafesTransfersListParams) SetValueLt(valueLt *string) {
	o.ValueLt = valueLt
}

// WriteToRequest writes these params to a swagger request
func (o *V1SafesTransfersListParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.From != nil {

		// query param _from
		var qrFrom string

		if o.From != nil {
			qrFrom = *o.From
		}
		qFrom := qrFrom
		if qFrom != "" {

			if err := r.SetQueryParam("_from", qFrom); err != nil {
				return err
			}
		}
	}

	// path param address
	if err := r.SetPathParam("address", o.Address); err != nil {
		return err
	}

	if o.BlockNumber != nil {

		// query param block_number
		var qrBlockNumber string

		if o.BlockNumber != nil {
			qrBlockNumber = *o.BlockNumber
		}
		qBlockNumber := qrBlockNumber
		if qBlockNumber != "" {

			if err := r.SetQueryParam("block_number", qBlockNumber); err != nil {
				return err
			}
		}
	}

	if o.BlockNumberGt != nil {

		// query param block_number__gt
		var qrBlockNumberGt string

		if o.BlockNumberGt != nil {
			qrBlockNumberGt = *o.BlockNumberGt
		}
		qBlockNumberGt := qrBlockNumberGt
		if qBlockNumberGt != "" {

			if err := r.SetQueryParam("block_number__gt", qBlockNumberGt); err != nil {
				return err
			}
		}
	}

	if o.BlockNumberLt != nil {

		// query param block_number__lt
		var qrBlockNumberLt string

		if o.BlockNumberLt != nil {
			qrBlockNumberLt = *o.BlockNumberLt
		}
		qBlockNumberLt := qrBlockNumberLt
		if qBlockNumberLt != "" {

			if err := r.SetQueryParam("block_number__lt", qBlockNumberLt); err != nil {
				return err
			}
		}
	}

	if o.Erc20 != nil {

		// query param erc20
		var qrErc20 string

		if o.Erc20 != nil {
			qrErc20 = *o.Erc20
		}
		qErc20 := qrErc20
		if qErc20 != "" {

			if err := r.SetQueryParam("erc20", qErc20); err != nil {
				return err
			}
		}
	}

	if o.Erc721 != nil {

		// query param erc721
		var qrErc721 string

		if o.Erc721 != nil {
			qrErc721 = *o.Erc721
		}
		qErc721 := qrErc721
		if qErc721 != "" {

			if err := r.SetQueryParam("erc721", qErc721); err != nil {
				return err
			}
		}
	}

	if o.Ether != nil {

		// query param ether
		var qrEther string

		if o.Ether != nil {
			qrEther = *o.Ether
		}
		qEther := qrEther
		if qEther != "" {

			if err := r.SetQueryParam("ether", qEther); err != nil {
				return err
			}
		}
	}

	if o.ExecutionDateGt != nil {

		// query param execution_date__gt
		var qrExecutionDateGt string

		if o.ExecutionDateGt != nil {
			qrExecutionDateGt = *o.ExecutionDateGt
		}
		qExecutionDateGt := qrExecutionDateGt
		if qExecutionDateGt != "" {

			if err := r.SetQueryParam("execution_date__gt", qExecutionDateGt); err != nil {
				return err
			}
		}
	}

	if o.ExecutionDateGte != nil {

		// query param execution_date__gte
		var qrExecutionDateGte string

		if o.ExecutionDateGte != nil {
			qrExecutionDateGte = *o.ExecutionDateGte
		}
		qExecutionDateGte := qrExecutionDateGte
		if qExecutionDateGte != "" {

			if err := r.SetQueryParam("execution_date__gte", qExecutionDateGte); err != nil {
				return err
			}
		}
	}

	if o.ExecutionDateLt != nil {

		// query param execution_date__lt
		var qrExecutionDateLt string

		if o.ExecutionDateLt != nil {
			qrExecutionDateLt = *o.ExecutionDateLt
		}
		qExecutionDateLt := qrExecutionDateLt
		if qExecutionDateLt != "" {

			if err := r.SetQueryParam("execution_date__lt", qExecutionDateLt); err != nil {
				return err
			}
		}
	}

	if o.ExecutionDateLte != nil {

		// query param execution_date__lte
		var qrExecutionDateLte string

		if o.ExecutionDateLte != nil {
			qrExecutionDateLte = *o.ExecutionDateLte
		}
		qExecutionDateLte := qrExecutionDateLte
		if qExecutionDateLte != "" {

			if err := r.SetQueryParam("execution_date__lte", qExecutionDateLte); err != nil {
				return err
			}
		}
	}

	if o.Limit != nil {

		// query param limit
		var qrLimit int64

		if o.Limit != nil {
			qrLimit = *o.Limit
		}
		qLimit := swag.FormatInt64(qrLimit)
		if qLimit != "" {

			if err := r.SetQueryParam("limit", qLimit); err != nil {
				return err
			}
		}
	}

	if o.Offset != nil {

		// query param offset
		var qrOffset int64

		if o.Offset != nil {
			qrOffset = *o.Offset
		}
		qOffset := swag.FormatInt64(qrOffset)
		if qOffset != "" {

			if err := r.SetQueryParam("offset", qOffset); err != nil {
				return err
			}
		}
	}

	if o.To != nil {

		// query param to
		var qrTo string

		if o.To != nil {
			qrTo = *o.To
		}
		qTo := qrTo
		if qTo != "" {

			if err := r.SetQueryParam("to", qTo); err != nil {
				return err
			}
		}
	}

	if o.TokenAddress != nil {

		// query param token_address
		var qrTokenAddress string

		if o.TokenAddress != nil {
			qrTokenAddress = *o.TokenAddress
		}
		qTokenAddress := qrTokenAddress
		if qTokenAddress != "" {

			if err := r.SetQueryParam("token_address", qTokenAddress); err != nil {
				return err
			}
		}
	}

	if o.TransactionHash != nil {

		// query param transaction_hash
		var qrTransactionHash string

		if o.TransactionHash != nil {
			qrTransactionHash = *o.TransactionHash
		}
		qTransactionHash := qrTransactionHash
		if qTransactionHash != "" {

			if err := r.SetQueryParam("transaction_hash", qTransactionHash); err != nil {
				return err
			}
		}
	}

	if o.Value != nil {

		// query param value
		var qrValue string

		if o.Value != nil {
			qrValue = *o.Value
		}
		qValue := qrValue
		if qValue != "" {

			if err := r.SetQueryParam("value", qValue); err != nil {
				return err
			}
		}
	}

	if o.ValueGt != nil {

		// query param value__gt
		var qrValueGt string

		if o.ValueGt != nil {
			qrValueGt = *o.ValueGt
		}
		qValueGt := qrValueGt
		if qValueGt != "" {

			if err := r.SetQueryParam("value__gt", qValueGt); err != nil {
				return err
			}
		}
	}

	if o.ValueLt != nil {

		// query param value__lt
		var qrValueLt string

		if o.ValueLt != nil {
			qrValueLt = *o.ValueLt
		}
		qValueLt := qrValueLt
		if qValueLt != "" {

			if err := r.SetQueryParam("value__lt", qValueLt); err != nil {
				return err
			}
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

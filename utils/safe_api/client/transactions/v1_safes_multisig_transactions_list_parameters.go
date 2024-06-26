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

// NewV1SafesMultisigTransactionsListParams creates a new V1SafesMultisigTransactionsListParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewV1SafesMultisigTransactionsListParams() *V1SafesMultisigTransactionsListParams {
	return &V1SafesMultisigTransactionsListParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewV1SafesMultisigTransactionsListParamsWithTimeout creates a new V1SafesMultisigTransactionsListParams object
// with the ability to set a timeout on a request.
func NewV1SafesMultisigTransactionsListParamsWithTimeout(timeout time.Duration) *V1SafesMultisigTransactionsListParams {
	return &V1SafesMultisigTransactionsListParams{
		timeout: timeout,
	}
}

// NewV1SafesMultisigTransactionsListParamsWithContext creates a new V1SafesMultisigTransactionsListParams object
// with the ability to set a context for a request.
func NewV1SafesMultisigTransactionsListParamsWithContext(ctx context.Context) *V1SafesMultisigTransactionsListParams {
	return &V1SafesMultisigTransactionsListParams{
		Context: ctx,
	}
}

// NewV1SafesMultisigTransactionsListParamsWithHTTPClient creates a new V1SafesMultisigTransactionsListParams object
// with the ability to set a custom HTTPClient for a request.
func NewV1SafesMultisigTransactionsListParamsWithHTTPClient(client *http.Client) *V1SafesMultisigTransactionsListParams {
	return &V1SafesMultisigTransactionsListParams{
		HTTPClient: client,
	}
}

/*
V1SafesMultisigTransactionsListParams contains all the parameters to send to the API endpoint

	for the v1 safes multisig transactions list operation.

	Typically these are written to a http.Request.
*/
type V1SafesMultisigTransactionsListParams struct {

	// Address.
	Address string

	/* Executed.

	   executed
	*/
	Executed *string

	/* ExecutionDateGte.

	   execution_date__gte
	*/
	ExecutionDateGte *string

	/* ExecutionDateLte.

	   execution_date__lte
	*/
	ExecutionDateLte *string

	/* Failed.

	   failed
	*/
	Failed *string

	/* HasConfirmations.

	   has_confirmations
	*/
	HasConfirmations *string

	/* Limit.

	   Number of results to return per page.
	*/
	Limit *int64

	/* ModifiedGt.

	   modified__gt
	*/
	ModifiedGt *string

	/* ModifiedGte.

	   modified__gte
	*/
	ModifiedGte *string

	/* ModifiedLt.

	   modified__lt
	*/
	ModifiedLt *string

	/* ModifiedLte.

	   modified__lte
	*/
	ModifiedLte *string

	/* Nonce.

	   nonce
	*/
	Nonce *string

	/* NonceGt.

	   nonce__gt
	*/
	NonceGt *string

	/* NonceGte.

	   nonce__gte
	*/
	NonceGte *string

	/* NonceLt.

	   nonce__lt
	*/
	NonceLt *string

	/* NonceLte.

	   nonce__lte
	*/
	NonceLte *string

	/* Offset.

	   The initial index from which to return the results.
	*/
	Offset *int64

	/* Ordering.

	   Which field to use when ordering the results.
	*/
	Ordering *string

	/* SafeTxHash.

	   safe_tx_hash
	*/
	SafeTxHash *string

	/* SubmissionDateGte.

	   submission_date__gte
	*/
	SubmissionDateGte *string

	/* SubmissionDateLte.

	   submission_date__lte
	*/
	SubmissionDateLte *string

	/* To.

	   to
	*/
	To *string

	/* TransactionHash.

	   transaction_hash
	*/
	TransactionHash *string

	/* Trusted.

	   trusted
	*/
	Trusted *string

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

// WithDefaults hydrates default values in the v1 safes multisig transactions list params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *V1SafesMultisigTransactionsListParams) WithDefaults() *V1SafesMultisigTransactionsListParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the v1 safes multisig transactions list params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *V1SafesMultisigTransactionsListParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) WithTimeout(timeout time.Duration) *V1SafesMultisigTransactionsListParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) WithContext(ctx context.Context) *V1SafesMultisigTransactionsListParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) WithHTTPClient(client *http.Client) *V1SafesMultisigTransactionsListParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithAddress adds the address to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) WithAddress(address string) *V1SafesMultisigTransactionsListParams {
	o.SetAddress(address)
	return o
}

// SetAddress adds the address to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) SetAddress(address string) {
	o.Address = address
}

// WithExecuted adds the executed to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) WithExecuted(executed *string) *V1SafesMultisigTransactionsListParams {
	o.SetExecuted(executed)
	return o
}

// SetExecuted adds the executed to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) SetExecuted(executed *string) {
	o.Executed = executed
}

// WithExecutionDateGte adds the executionDateGte to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) WithExecutionDateGte(executionDateGte *string) *V1SafesMultisigTransactionsListParams {
	o.SetExecutionDateGte(executionDateGte)
	return o
}

// SetExecutionDateGte adds the executionDateGte to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) SetExecutionDateGte(executionDateGte *string) {
	o.ExecutionDateGte = executionDateGte
}

// WithExecutionDateLte adds the executionDateLte to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) WithExecutionDateLte(executionDateLte *string) *V1SafesMultisigTransactionsListParams {
	o.SetExecutionDateLte(executionDateLte)
	return o
}

// SetExecutionDateLte adds the executionDateLte to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) SetExecutionDateLte(executionDateLte *string) {
	o.ExecutionDateLte = executionDateLte
}

// WithFailed adds the failed to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) WithFailed(failed *string) *V1SafesMultisigTransactionsListParams {
	o.SetFailed(failed)
	return o
}

// SetFailed adds the failed to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) SetFailed(failed *string) {
	o.Failed = failed
}

// WithHasConfirmations adds the hasConfirmations to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) WithHasConfirmations(hasConfirmations *string) *V1SafesMultisigTransactionsListParams {
	o.SetHasConfirmations(hasConfirmations)
	return o
}

// SetHasConfirmations adds the hasConfirmations to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) SetHasConfirmations(hasConfirmations *string) {
	o.HasConfirmations = hasConfirmations
}

// WithLimit adds the limit to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) WithLimit(limit *int64) *V1SafesMultisigTransactionsListParams {
	o.SetLimit(limit)
	return o
}

// SetLimit adds the limit to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) SetLimit(limit *int64) {
	o.Limit = limit
}

// WithModifiedGt adds the modifiedGt to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) WithModifiedGt(modifiedGt *string) *V1SafesMultisigTransactionsListParams {
	o.SetModifiedGt(modifiedGt)
	return o
}

// SetModifiedGt adds the modifiedGt to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) SetModifiedGt(modifiedGt *string) {
	o.ModifiedGt = modifiedGt
}

// WithModifiedGte adds the modifiedGte to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) WithModifiedGte(modifiedGte *string) *V1SafesMultisigTransactionsListParams {
	o.SetModifiedGte(modifiedGte)
	return o
}

// SetModifiedGte adds the modifiedGte to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) SetModifiedGte(modifiedGte *string) {
	o.ModifiedGte = modifiedGte
}

// WithModifiedLt adds the modifiedLt to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) WithModifiedLt(modifiedLt *string) *V1SafesMultisigTransactionsListParams {
	o.SetModifiedLt(modifiedLt)
	return o
}

// SetModifiedLt adds the modifiedLt to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) SetModifiedLt(modifiedLt *string) {
	o.ModifiedLt = modifiedLt
}

// WithModifiedLte adds the modifiedLte to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) WithModifiedLte(modifiedLte *string) *V1SafesMultisigTransactionsListParams {
	o.SetModifiedLte(modifiedLte)
	return o
}

// SetModifiedLte adds the modifiedLte to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) SetModifiedLte(modifiedLte *string) {
	o.ModifiedLte = modifiedLte
}

// WithNonce adds the nonce to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) WithNonce(nonce *string) *V1SafesMultisigTransactionsListParams {
	o.SetNonce(nonce)
	return o
}

// SetNonce adds the nonce to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) SetNonce(nonce *string) {
	o.Nonce = nonce
}

// WithNonceGt adds the nonceGt to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) WithNonceGt(nonceGt *string) *V1SafesMultisigTransactionsListParams {
	o.SetNonceGt(nonceGt)
	return o
}

// SetNonceGt adds the nonceGt to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) SetNonceGt(nonceGt *string) {
	o.NonceGt = nonceGt
}

// WithNonceGte adds the nonceGte to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) WithNonceGte(nonceGte *string) *V1SafesMultisigTransactionsListParams {
	o.SetNonceGte(nonceGte)
	return o
}

// SetNonceGte adds the nonceGte to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) SetNonceGte(nonceGte *string) {
	o.NonceGte = nonceGte
}

// WithNonceLt adds the nonceLt to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) WithNonceLt(nonceLt *string) *V1SafesMultisigTransactionsListParams {
	o.SetNonceLt(nonceLt)
	return o
}

// SetNonceLt adds the nonceLt to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) SetNonceLt(nonceLt *string) {
	o.NonceLt = nonceLt
}

// WithNonceLte adds the nonceLte to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) WithNonceLte(nonceLte *string) *V1SafesMultisigTransactionsListParams {
	o.SetNonceLte(nonceLte)
	return o
}

// SetNonceLte adds the nonceLte to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) SetNonceLte(nonceLte *string) {
	o.NonceLte = nonceLte
}

// WithOffset adds the offset to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) WithOffset(offset *int64) *V1SafesMultisigTransactionsListParams {
	o.SetOffset(offset)
	return o
}

// SetOffset adds the offset to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) SetOffset(offset *int64) {
	o.Offset = offset
}

// WithOrdering adds the ordering to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) WithOrdering(ordering *string) *V1SafesMultisigTransactionsListParams {
	o.SetOrdering(ordering)
	return o
}

// SetOrdering adds the ordering to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) SetOrdering(ordering *string) {
	o.Ordering = ordering
}

// WithSafeTxHash adds the safeTxHash to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) WithSafeTxHash(safeTxHash *string) *V1SafesMultisigTransactionsListParams {
	o.SetSafeTxHash(safeTxHash)
	return o
}

// SetSafeTxHash adds the safeTxHash to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) SetSafeTxHash(safeTxHash *string) {
	o.SafeTxHash = safeTxHash
}

// WithSubmissionDateGte adds the submissionDateGte to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) WithSubmissionDateGte(submissionDateGte *string) *V1SafesMultisigTransactionsListParams {
	o.SetSubmissionDateGte(submissionDateGte)
	return o
}

// SetSubmissionDateGte adds the submissionDateGte to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) SetSubmissionDateGte(submissionDateGte *string) {
	o.SubmissionDateGte = submissionDateGte
}

// WithSubmissionDateLte adds the submissionDateLte to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) WithSubmissionDateLte(submissionDateLte *string) *V1SafesMultisigTransactionsListParams {
	o.SetSubmissionDateLte(submissionDateLte)
	return o
}

// SetSubmissionDateLte adds the submissionDateLte to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) SetSubmissionDateLte(submissionDateLte *string) {
	o.SubmissionDateLte = submissionDateLte
}

// WithTo adds the to to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) WithTo(to *string) *V1SafesMultisigTransactionsListParams {
	o.SetTo(to)
	return o
}

// SetTo adds the to to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) SetTo(to *string) {
	o.To = to
}

// WithTransactionHash adds the transactionHash to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) WithTransactionHash(transactionHash *string) *V1SafesMultisigTransactionsListParams {
	o.SetTransactionHash(transactionHash)
	return o
}

// SetTransactionHash adds the transactionHash to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) SetTransactionHash(transactionHash *string) {
	o.TransactionHash = transactionHash
}

// WithTrusted adds the trusted to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) WithTrusted(trusted *string) *V1SafesMultisigTransactionsListParams {
	o.SetTrusted(trusted)
	return o
}

// SetTrusted adds the trusted to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) SetTrusted(trusted *string) {
	o.Trusted = trusted
}

// WithValue adds the value to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) WithValue(value *string) *V1SafesMultisigTransactionsListParams {
	o.SetValue(value)
	return o
}

// SetValue adds the value to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) SetValue(value *string) {
	o.Value = value
}

// WithValueGt adds the valueGt to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) WithValueGt(valueGt *string) *V1SafesMultisigTransactionsListParams {
	o.SetValueGt(valueGt)
	return o
}

// SetValueGt adds the valueGt to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) SetValueGt(valueGt *string) {
	o.ValueGt = valueGt
}

// WithValueLt adds the valueLt to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) WithValueLt(valueLt *string) *V1SafesMultisigTransactionsListParams {
	o.SetValueLt(valueLt)
	return o
}

// SetValueLt adds the valueLt to the v1 safes multisig transactions list params
func (o *V1SafesMultisigTransactionsListParams) SetValueLt(valueLt *string) {
	o.ValueLt = valueLt
}

// WriteToRequest writes these params to a swagger request
func (o *V1SafesMultisigTransactionsListParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param address
	if err := r.SetPathParam("address", o.Address); err != nil {
		return err
	}

	if o.Executed != nil {

		// query param executed
		var qrExecuted string

		if o.Executed != nil {
			qrExecuted = *o.Executed
		}
		qExecuted := qrExecuted
		if qExecuted != "" {

			if err := r.SetQueryParam("executed", qExecuted); err != nil {
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

	if o.Failed != nil {

		// query param failed
		var qrFailed string

		if o.Failed != nil {
			qrFailed = *o.Failed
		}
		qFailed := qrFailed
		if qFailed != "" {

			if err := r.SetQueryParam("failed", qFailed); err != nil {
				return err
			}
		}
	}

	if o.HasConfirmations != nil {

		// query param has_confirmations
		var qrHasConfirmations string

		if o.HasConfirmations != nil {
			qrHasConfirmations = *o.HasConfirmations
		}
		qHasConfirmations := qrHasConfirmations
		if qHasConfirmations != "" {

			if err := r.SetQueryParam("has_confirmations", qHasConfirmations); err != nil {
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

	if o.ModifiedGt != nil {

		// query param modified__gt
		var qrModifiedGt string

		if o.ModifiedGt != nil {
			qrModifiedGt = *o.ModifiedGt
		}
		qModifiedGt := qrModifiedGt
		if qModifiedGt != "" {

			if err := r.SetQueryParam("modified__gt", qModifiedGt); err != nil {
				return err
			}
		}
	}

	if o.ModifiedGte != nil {

		// query param modified__gte
		var qrModifiedGte string

		if o.ModifiedGte != nil {
			qrModifiedGte = *o.ModifiedGte
		}
		qModifiedGte := qrModifiedGte
		if qModifiedGte != "" {

			if err := r.SetQueryParam("modified__gte", qModifiedGte); err != nil {
				return err
			}
		}
	}

	if o.ModifiedLt != nil {

		// query param modified__lt
		var qrModifiedLt string

		if o.ModifiedLt != nil {
			qrModifiedLt = *o.ModifiedLt
		}
		qModifiedLt := qrModifiedLt
		if qModifiedLt != "" {

			if err := r.SetQueryParam("modified__lt", qModifiedLt); err != nil {
				return err
			}
		}
	}

	if o.ModifiedLte != nil {

		// query param modified__lte
		var qrModifiedLte string

		if o.ModifiedLte != nil {
			qrModifiedLte = *o.ModifiedLte
		}
		qModifiedLte := qrModifiedLte
		if qModifiedLte != "" {

			if err := r.SetQueryParam("modified__lte", qModifiedLte); err != nil {
				return err
			}
		}
	}

	if o.Nonce != nil {

		// query param nonce
		var qrNonce string

		if o.Nonce != nil {
			qrNonce = *o.Nonce
		}
		qNonce := qrNonce
		if qNonce != "" {

			if err := r.SetQueryParam("nonce", qNonce); err != nil {
				return err
			}
		}
	}

	if o.NonceGt != nil {

		// query param nonce__gt
		var qrNonceGt string

		if o.NonceGt != nil {
			qrNonceGt = *o.NonceGt
		}
		qNonceGt := qrNonceGt
		if qNonceGt != "" {

			if err := r.SetQueryParam("nonce__gt", qNonceGt); err != nil {
				return err
			}
		}
	}

	if o.NonceGte != nil {

		// query param nonce__gte
		var qrNonceGte string

		if o.NonceGte != nil {
			qrNonceGte = *o.NonceGte
		}
		qNonceGte := qrNonceGte
		if qNonceGte != "" {

			if err := r.SetQueryParam("nonce__gte", qNonceGte); err != nil {
				return err
			}
		}
	}

	if o.NonceLt != nil {

		// query param nonce__lt
		var qrNonceLt string

		if o.NonceLt != nil {
			qrNonceLt = *o.NonceLt
		}
		qNonceLt := qrNonceLt
		if qNonceLt != "" {

			if err := r.SetQueryParam("nonce__lt", qNonceLt); err != nil {
				return err
			}
		}
	}

	if o.NonceLte != nil {

		// query param nonce__lte
		var qrNonceLte string

		if o.NonceLte != nil {
			qrNonceLte = *o.NonceLte
		}
		qNonceLte := qrNonceLte
		if qNonceLte != "" {

			if err := r.SetQueryParam("nonce__lte", qNonceLte); err != nil {
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

	if o.Ordering != nil {

		// query param ordering
		var qrOrdering string

		if o.Ordering != nil {
			qrOrdering = *o.Ordering
		}
		qOrdering := qrOrdering
		if qOrdering != "" {

			if err := r.SetQueryParam("ordering", qOrdering); err != nil {
				return err
			}
		}
	}

	if o.SafeTxHash != nil {

		// query param safe_tx_hash
		var qrSafeTxHash string

		if o.SafeTxHash != nil {
			qrSafeTxHash = *o.SafeTxHash
		}
		qSafeTxHash := qrSafeTxHash
		if qSafeTxHash != "" {

			if err := r.SetQueryParam("safe_tx_hash", qSafeTxHash); err != nil {
				return err
			}
		}
	}

	if o.SubmissionDateGte != nil {

		// query param submission_date__gte
		var qrSubmissionDateGte string

		if o.SubmissionDateGte != nil {
			qrSubmissionDateGte = *o.SubmissionDateGte
		}
		qSubmissionDateGte := qrSubmissionDateGte
		if qSubmissionDateGte != "" {

			if err := r.SetQueryParam("submission_date__gte", qSubmissionDateGte); err != nil {
				return err
			}
		}
	}

	if o.SubmissionDateLte != nil {

		// query param submission_date__lte
		var qrSubmissionDateLte string

		if o.SubmissionDateLte != nil {
			qrSubmissionDateLte = *o.SubmissionDateLte
		}
		qSubmissionDateLte := qrSubmissionDateLte
		if qSubmissionDateLte != "" {

			if err := r.SetQueryParam("submission_date__lte", qSubmissionDateLte); err != nil {
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

	if o.Trusted != nil {

		// query param trusted
		var qrTrusted string

		if o.Trusted != nil {
			qrTrusted = *o.Trusted
		}
		qTrusted := qrTrusted
		if qTrusted != "" {

			if err := r.SetQueryParam("trusted", qTrusted); err != nil {
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

// Code generated by go-swagger; DO NOT EDIT.

package transactions

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	models "omni-balance/utils/safe_api"
)

// V1SafesMultisigTransactionsListReader is a Reader for the V1SafesMultisigTransactionsList structure.
type V1SafesMultisigTransactionsListReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *V1SafesMultisigTransactionsListReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewV1SafesMultisigTransactionsListOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewV1SafesMultisigTransactionsListBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 422:
		result := NewV1SafesMultisigTransactionsListUnprocessableEntity()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[GET /v1/safes/{address}/multisig-transactions/] v1_safes_multisig-transactions_list", response, response.Code())
	}
}

// NewV1SafesMultisigTransactionsListOK creates a V1SafesMultisigTransactionsListOK with default headers values
func NewV1SafesMultisigTransactionsListOK() *V1SafesMultisigTransactionsListOK {
	return &V1SafesMultisigTransactionsListOK{}
}

/*
V1SafesMultisigTransactionsListOK describes a response with status code 200, with default header values.

V1SafesMultisigTransactionsListOK v1 safes multisig transactions list o k
*/
type V1SafesMultisigTransactionsListOK struct {
	Payload *V1SafesMultisigTransactionsListOKBody
}

// IsSuccess returns true when this v1 safes multisig transactions list o k response has a 2xx status code
func (o *V1SafesMultisigTransactionsListOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this v1 safes multisig transactions list o k response has a 3xx status code
func (o *V1SafesMultisigTransactionsListOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this v1 safes multisig transactions list o k response has a 4xx status code
func (o *V1SafesMultisigTransactionsListOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this v1 safes multisig transactions list o k response has a 5xx status code
func (o *V1SafesMultisigTransactionsListOK) IsServerError() bool {
	return false
}

// IsCode returns true when this v1 safes multisig transactions list o k response a status code equal to that given
func (o *V1SafesMultisigTransactionsListOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the v1 safes multisig transactions list o k response
func (o *V1SafesMultisigTransactionsListOK) Code() int {
	return 200
}

func (o *V1SafesMultisigTransactionsListOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /v1/safes/{address}/multisig-transactions/][%d] v1SafesMultisigTransactionsListOK %s", 200, payload)
}

func (o *V1SafesMultisigTransactionsListOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /v1/safes/{address}/multisig-transactions/][%d] v1SafesMultisigTransactionsListOK %s", 200, payload)
}

func (o *V1SafesMultisigTransactionsListOK) GetPayload() *V1SafesMultisigTransactionsListOKBody {
	return o.Payload
}

func (o *V1SafesMultisigTransactionsListOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(V1SafesMultisigTransactionsListOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewV1SafesMultisigTransactionsListBadRequest creates a V1SafesMultisigTransactionsListBadRequest with default headers values
func NewV1SafesMultisigTransactionsListBadRequest() *V1SafesMultisigTransactionsListBadRequest {
	return &V1SafesMultisigTransactionsListBadRequest{}
}

/*
V1SafesMultisigTransactionsListBadRequest describes a response with status code 400, with default header values.

Invalid data
*/
type V1SafesMultisigTransactionsListBadRequest struct {
}

// IsSuccess returns true when this v1 safes multisig transactions list bad request response has a 2xx status code
func (o *V1SafesMultisigTransactionsListBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this v1 safes multisig transactions list bad request response has a 3xx status code
func (o *V1SafesMultisigTransactionsListBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this v1 safes multisig transactions list bad request response has a 4xx status code
func (o *V1SafesMultisigTransactionsListBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this v1 safes multisig transactions list bad request response has a 5xx status code
func (o *V1SafesMultisigTransactionsListBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this v1 safes multisig transactions list bad request response a status code equal to that given
func (o *V1SafesMultisigTransactionsListBadRequest) IsCode(code int) bool {
	return code == 400
}

// Code gets the status code for the v1 safes multisig transactions list bad request response
func (o *V1SafesMultisigTransactionsListBadRequest) Code() int {
	return 400
}

func (o *V1SafesMultisigTransactionsListBadRequest) Error() string {
	return fmt.Sprintf("[GET /v1/safes/{address}/multisig-transactions/][%d] v1SafesMultisigTransactionsListBadRequest", 400)
}

func (o *V1SafesMultisigTransactionsListBadRequest) String() string {
	return fmt.Sprintf("[GET /v1/safes/{address}/multisig-transactions/][%d] v1SafesMultisigTransactionsListBadRequest", 400)
}

func (o *V1SafesMultisigTransactionsListBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewV1SafesMultisigTransactionsListUnprocessableEntity creates a V1SafesMultisigTransactionsListUnprocessableEntity with default headers values
func NewV1SafesMultisigTransactionsListUnprocessableEntity() *V1SafesMultisigTransactionsListUnprocessableEntity {
	return &V1SafesMultisigTransactionsListUnprocessableEntity{}
}

/*
V1SafesMultisigTransactionsListUnprocessableEntity describes a response with status code 422, with default header values.

Invalid ethereum address
*/
type V1SafesMultisigTransactionsListUnprocessableEntity struct {
}

// IsSuccess returns true when this v1 safes multisig transactions list unprocessable entity response has a 2xx status code
func (o *V1SafesMultisigTransactionsListUnprocessableEntity) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this v1 safes multisig transactions list unprocessable entity response has a 3xx status code
func (o *V1SafesMultisigTransactionsListUnprocessableEntity) IsRedirect() bool {
	return false
}

// IsClientError returns true when this v1 safes multisig transactions list unprocessable entity response has a 4xx status code
func (o *V1SafesMultisigTransactionsListUnprocessableEntity) IsClientError() bool {
	return true
}

// IsServerError returns true when this v1 safes multisig transactions list unprocessable entity response has a 5xx status code
func (o *V1SafesMultisigTransactionsListUnprocessableEntity) IsServerError() bool {
	return false
}

// IsCode returns true when this v1 safes multisig transactions list unprocessable entity response a status code equal to that given
func (o *V1SafesMultisigTransactionsListUnprocessableEntity) IsCode(code int) bool {
	return code == 422
}

// Code gets the status code for the v1 safes multisig transactions list unprocessable entity response
func (o *V1SafesMultisigTransactionsListUnprocessableEntity) Code() int {
	return 422
}

func (o *V1SafesMultisigTransactionsListUnprocessableEntity) Error() string {
	return fmt.Sprintf("[GET /v1/safes/{address}/multisig-transactions/][%d] v1SafesMultisigTransactionsListUnprocessableEntity", 422)
}

func (o *V1SafesMultisigTransactionsListUnprocessableEntity) String() string {
	return fmt.Sprintf("[GET /v1/safes/{address}/multisig-transactions/][%d] v1SafesMultisigTransactionsListUnprocessableEntity", 422)
}

func (o *V1SafesMultisigTransactionsListUnprocessableEntity) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

/*
V1SafesMultisigTransactionsListOKBody v1 safes multisig transactions list o k body
swagger:model V1SafesMultisigTransactionsListOKBody
*/
type V1SafesMultisigTransactionsListOKBody struct {

	// count
	// Required: true
	Count *int64 `json:"count"`

	// next
	// Format: uri
	Next *strfmt.URI `json:"next,omitempty"`

	// previous
	// Format: uri
	Previous *strfmt.URI `json:"previous,omitempty"`

	// results
	// Required: true
	Results []*models.SafeMultisigTransactionResponse `json:"results"`
}

// Validate validates this v1 safes multisig transactions list o k body
func (o *V1SafesMultisigTransactionsListOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateCount(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateNext(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validatePrevious(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateResults(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *V1SafesMultisigTransactionsListOKBody) validateCount(formats strfmt.Registry) error {

	if err := validate.Required("v1SafesMultisigTransactionsListOK"+"."+"count", "body", o.Count); err != nil {
		return err
	}

	return nil
}

func (o *V1SafesMultisigTransactionsListOKBody) validateNext(formats strfmt.Registry) error {
	if swag.IsZero(o.Next) { // not required
		return nil
	}

	if err := validate.FormatOf("v1SafesMultisigTransactionsListOK"+"."+"next", "body", "uri", o.Next.String(), formats); err != nil {
		return err
	}

	return nil
}

func (o *V1SafesMultisigTransactionsListOKBody) validatePrevious(formats strfmt.Registry) error {
	if swag.IsZero(o.Previous) { // not required
		return nil
	}

	if err := validate.FormatOf("v1SafesMultisigTransactionsListOK"+"."+"previous", "body", "uri", o.Previous.String(), formats); err != nil {
		return err
	}

	return nil
}

func (o *V1SafesMultisigTransactionsListOKBody) validateResults(formats strfmt.Registry) error {

	if err := validate.Required("v1SafesMultisigTransactionsListOK"+"."+"results", "body", o.Results); err != nil {
		return err
	}

	for i := 0; i < len(o.Results); i++ {
		if swag.IsZero(o.Results[i]) { // not required
			continue
		}

		if o.Results[i] != nil {
			if err := o.Results[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("v1SafesMultisigTransactionsListOK" + "." + "results" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("v1SafesMultisigTransactionsListOK" + "." + "results" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this v1 safes multisig transactions list o k body based on the context it is used
func (o *V1SafesMultisigTransactionsListOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateResults(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *V1SafesMultisigTransactionsListOKBody) contextValidateResults(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.Results); i++ {

		if o.Results[i] != nil {

			if swag.IsZero(o.Results[i]) { // not required
				return nil
			}

			if err := o.Results[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("v1SafesMultisigTransactionsListOK" + "." + "results" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("v1SafesMultisigTransactionsListOK" + "." + "results" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *V1SafesMultisigTransactionsListOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *V1SafesMultisigTransactionsListOKBody) UnmarshalBinary(b []byte) error {
	var res V1SafesMultisigTransactionsListOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

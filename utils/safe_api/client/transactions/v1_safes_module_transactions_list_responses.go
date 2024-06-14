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

// V1SafesModuleTransactionsListReader is a Reader for the V1SafesModuleTransactionsList structure.
type V1SafesModuleTransactionsListReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *V1SafesModuleTransactionsListReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewV1SafesModuleTransactionsListOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewV1SafesModuleTransactionsListBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 422:
		result := NewV1SafesModuleTransactionsListUnprocessableEntity()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[GET /v1/safes/{address}/module-transactions/] v1_safes_module-transactions_list", response, response.Code())
	}
}

// NewV1SafesModuleTransactionsListOK creates a V1SafesModuleTransactionsListOK with default headers values
func NewV1SafesModuleTransactionsListOK() *V1SafesModuleTransactionsListOK {
	return &V1SafesModuleTransactionsListOK{}
}

/*
V1SafesModuleTransactionsListOK describes a response with status code 200, with default header values.

V1SafesModuleTransactionsListOK v1 safes module transactions list o k
*/
type V1SafesModuleTransactionsListOK struct {
	Payload *V1SafesModuleTransactionsListOKBody
}

// IsSuccess returns true when this v1 safes module transactions list o k response has a 2xx status code
func (o *V1SafesModuleTransactionsListOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this v1 safes module transactions list o k response has a 3xx status code
func (o *V1SafesModuleTransactionsListOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this v1 safes module transactions list o k response has a 4xx status code
func (o *V1SafesModuleTransactionsListOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this v1 safes module transactions list o k response has a 5xx status code
func (o *V1SafesModuleTransactionsListOK) IsServerError() bool {
	return false
}

// IsCode returns true when this v1 safes module transactions list o k response a status code equal to that given
func (o *V1SafesModuleTransactionsListOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the v1 safes module transactions list o k response
func (o *V1SafesModuleTransactionsListOK) Code() int {
	return 200
}

func (o *V1SafesModuleTransactionsListOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /v1/safes/{address}/module-transactions/][%d] v1SafesModuleTransactionsListOK %s", 200, payload)
}

func (o *V1SafesModuleTransactionsListOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /v1/safes/{address}/module-transactions/][%d] v1SafesModuleTransactionsListOK %s", 200, payload)
}

func (o *V1SafesModuleTransactionsListOK) GetPayload() *V1SafesModuleTransactionsListOKBody {
	return o.Payload
}

func (o *V1SafesModuleTransactionsListOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(V1SafesModuleTransactionsListOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewV1SafesModuleTransactionsListBadRequest creates a V1SafesModuleTransactionsListBadRequest with default headers values
func NewV1SafesModuleTransactionsListBadRequest() *V1SafesModuleTransactionsListBadRequest {
	return &V1SafesModuleTransactionsListBadRequest{}
}

/*
V1SafesModuleTransactionsListBadRequest describes a response with status code 400, with default header values.

Invalid data
*/
type V1SafesModuleTransactionsListBadRequest struct {
}

// IsSuccess returns true when this v1 safes module transactions list bad request response has a 2xx status code
func (o *V1SafesModuleTransactionsListBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this v1 safes module transactions list bad request response has a 3xx status code
func (o *V1SafesModuleTransactionsListBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this v1 safes module transactions list bad request response has a 4xx status code
func (o *V1SafesModuleTransactionsListBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this v1 safes module transactions list bad request response has a 5xx status code
func (o *V1SafesModuleTransactionsListBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this v1 safes module transactions list bad request response a status code equal to that given
func (o *V1SafesModuleTransactionsListBadRequest) IsCode(code int) bool {
	return code == 400
}

// Code gets the status code for the v1 safes module transactions list bad request response
func (o *V1SafesModuleTransactionsListBadRequest) Code() int {
	return 400
}

func (o *V1SafesModuleTransactionsListBadRequest) Error() string {
	return fmt.Sprintf("[GET /v1/safes/{address}/module-transactions/][%d] v1SafesModuleTransactionsListBadRequest", 400)
}

func (o *V1SafesModuleTransactionsListBadRequest) String() string {
	return fmt.Sprintf("[GET /v1/safes/{address}/module-transactions/][%d] v1SafesModuleTransactionsListBadRequest", 400)
}

func (o *V1SafesModuleTransactionsListBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewV1SafesModuleTransactionsListUnprocessableEntity creates a V1SafesModuleTransactionsListUnprocessableEntity with default headers values
func NewV1SafesModuleTransactionsListUnprocessableEntity() *V1SafesModuleTransactionsListUnprocessableEntity {
	return &V1SafesModuleTransactionsListUnprocessableEntity{}
}

/*
V1SafesModuleTransactionsListUnprocessableEntity describes a response with status code 422, with default header values.

Invalid ethereum address
*/
type V1SafesModuleTransactionsListUnprocessableEntity struct {
}

// IsSuccess returns true when this v1 safes module transactions list unprocessable entity response has a 2xx status code
func (o *V1SafesModuleTransactionsListUnprocessableEntity) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this v1 safes module transactions list unprocessable entity response has a 3xx status code
func (o *V1SafesModuleTransactionsListUnprocessableEntity) IsRedirect() bool {
	return false
}

// IsClientError returns true when this v1 safes module transactions list unprocessable entity response has a 4xx status code
func (o *V1SafesModuleTransactionsListUnprocessableEntity) IsClientError() bool {
	return true
}

// IsServerError returns true when this v1 safes module transactions list unprocessable entity response has a 5xx status code
func (o *V1SafesModuleTransactionsListUnprocessableEntity) IsServerError() bool {
	return false
}

// IsCode returns true when this v1 safes module transactions list unprocessable entity response a status code equal to that given
func (o *V1SafesModuleTransactionsListUnprocessableEntity) IsCode(code int) bool {
	return code == 422
}

// Code gets the status code for the v1 safes module transactions list unprocessable entity response
func (o *V1SafesModuleTransactionsListUnprocessableEntity) Code() int {
	return 422
}

func (o *V1SafesModuleTransactionsListUnprocessableEntity) Error() string {
	return fmt.Sprintf("[GET /v1/safes/{address}/module-transactions/][%d] v1SafesModuleTransactionsListUnprocessableEntity", 422)
}

func (o *V1SafesModuleTransactionsListUnprocessableEntity) String() string {
	return fmt.Sprintf("[GET /v1/safes/{address}/module-transactions/][%d] v1SafesModuleTransactionsListUnprocessableEntity", 422)
}

func (o *V1SafesModuleTransactionsListUnprocessableEntity) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

/*
V1SafesModuleTransactionsListOKBody v1 safes module transactions list o k body
swagger:model V1SafesModuleTransactionsListOKBody
*/
type V1SafesModuleTransactionsListOKBody struct {

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
	Results []*models.SafeModuleTransactionResponse `json:"results"`
}

// Validate validates this v1 safes module transactions list o k body
func (o *V1SafesModuleTransactionsListOKBody) Validate(formats strfmt.Registry) error {
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

func (o *V1SafesModuleTransactionsListOKBody) validateCount(formats strfmt.Registry) error {

	if err := validate.Required("v1SafesModuleTransactionsListOK"+"."+"count", "body", o.Count); err != nil {
		return err
	}

	return nil
}

func (o *V1SafesModuleTransactionsListOKBody) validateNext(formats strfmt.Registry) error {
	if swag.IsZero(o.Next) { // not required
		return nil
	}

	if err := validate.FormatOf("v1SafesModuleTransactionsListOK"+"."+"next", "body", "uri", o.Next.String(), formats); err != nil {
		return err
	}

	return nil
}

func (o *V1SafesModuleTransactionsListOKBody) validatePrevious(formats strfmt.Registry) error {
	if swag.IsZero(o.Previous) { // not required
		return nil
	}

	if err := validate.FormatOf("v1SafesModuleTransactionsListOK"+"."+"previous", "body", "uri", o.Previous.String(), formats); err != nil {
		return err
	}

	return nil
}

func (o *V1SafesModuleTransactionsListOKBody) validateResults(formats strfmt.Registry) error {

	if err := validate.Required("v1SafesModuleTransactionsListOK"+"."+"results", "body", o.Results); err != nil {
		return err
	}

	for i := 0; i < len(o.Results); i++ {
		if swag.IsZero(o.Results[i]) { // not required
			continue
		}

		if o.Results[i] != nil {
			if err := o.Results[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("v1SafesModuleTransactionsListOK" + "." + "results" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("v1SafesModuleTransactionsListOK" + "." + "results" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this v1 safes module transactions list o k body based on the context it is used
func (o *V1SafesModuleTransactionsListOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateResults(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *V1SafesModuleTransactionsListOKBody) contextValidateResults(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.Results); i++ {

		if o.Results[i] != nil {

			if swag.IsZero(o.Results[i]) { // not required
				return nil
			}

			if err := o.Results[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("v1SafesModuleTransactionsListOK" + "." + "results" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("v1SafesModuleTransactionsListOK" + "." + "results" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *V1SafesModuleTransactionsListOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *V1SafesModuleTransactionsListOKBody) UnmarshalBinary(b []byte) error {
	var res V1SafesModuleTransactionsListOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

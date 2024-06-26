// Code generated by go-swagger; DO NOT EDIT.

package contracts

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

// V1ContractsListReader is a Reader for the V1ContractsList structure.
type V1ContractsListReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *V1ContractsListReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewV1ContractsListOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("[GET /v1/contracts/] v1_contracts_list", response, response.Code())
	}
}

// NewV1ContractsListOK creates a V1ContractsListOK with default headers values
func NewV1ContractsListOK() *V1ContractsListOK {
	return &V1ContractsListOK{}
}

/*
V1ContractsListOK describes a response with status code 200, with default header values.

V1ContractsListOK v1 contracts list o k
*/
type V1ContractsListOK struct {
	Payload *V1ContractsListOKBody
}

// IsSuccess returns true when this v1 contracts list o k response has a 2xx status code
func (o *V1ContractsListOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this v1 contracts list o k response has a 3xx status code
func (o *V1ContractsListOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this v1 contracts list o k response has a 4xx status code
func (o *V1ContractsListOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this v1 contracts list o k response has a 5xx status code
func (o *V1ContractsListOK) IsServerError() bool {
	return false
}

// IsCode returns true when this v1 contracts list o k response a status code equal to that given
func (o *V1ContractsListOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the v1 contracts list o k response
func (o *V1ContractsListOK) Code() int {
	return 200
}

func (o *V1ContractsListOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /v1/contracts/][%d] v1ContractsListOK %s", 200, payload)
}

func (o *V1ContractsListOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /v1/contracts/][%d] v1ContractsListOK %s", 200, payload)
}

func (o *V1ContractsListOK) GetPayload() *V1ContractsListOKBody {
	return o.Payload
}

func (o *V1ContractsListOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(V1ContractsListOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*
V1ContractsListOKBody v1 contracts list o k body
swagger:model V1ContractsListOKBody
*/
type V1ContractsListOKBody struct {

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
	Results []*models.Contract `json:"results"`
}

// Validate validates this v1 contracts list o k body
func (o *V1ContractsListOKBody) Validate(formats strfmt.Registry) error {
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

func (o *V1ContractsListOKBody) validateCount(formats strfmt.Registry) error {

	if err := validate.Required("v1ContractsListOK"+"."+"count", "body", o.Count); err != nil {
		return err
	}

	return nil
}

func (o *V1ContractsListOKBody) validateNext(formats strfmt.Registry) error {
	if swag.IsZero(o.Next) { // not required
		return nil
	}

	if err := validate.FormatOf("v1ContractsListOK"+"."+"next", "body", "uri", o.Next.String(), formats); err != nil {
		return err
	}

	return nil
}

func (o *V1ContractsListOKBody) validatePrevious(formats strfmt.Registry) error {
	if swag.IsZero(o.Previous) { // not required
		return nil
	}

	if err := validate.FormatOf("v1ContractsListOK"+"."+"previous", "body", "uri", o.Previous.String(), formats); err != nil {
		return err
	}

	return nil
}

func (o *V1ContractsListOKBody) validateResults(formats strfmt.Registry) error {

	if err := validate.Required("v1ContractsListOK"+"."+"results", "body", o.Results); err != nil {
		return err
	}

	for i := 0; i < len(o.Results); i++ {
		if swag.IsZero(o.Results[i]) { // not required
			continue
		}

		if o.Results[i] != nil {
			if err := o.Results[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("v1ContractsListOK" + "." + "results" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("v1ContractsListOK" + "." + "results" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this v1 contracts list o k body based on the context it is used
func (o *V1ContractsListOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateResults(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *V1ContractsListOKBody) contextValidateResults(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.Results); i++ {

		if o.Results[i] != nil {

			if swag.IsZero(o.Results[i]) { // not required
				return nil
			}

			if err := o.Results[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("v1ContractsListOK" + "." + "results" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("v1ContractsListOK" + "." + "results" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *V1ContractsListOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *V1ContractsListOKBody) UnmarshalBinary(b []byte) error {
	var res V1ContractsListOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

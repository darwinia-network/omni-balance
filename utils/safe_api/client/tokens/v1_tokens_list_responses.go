// Code generated by go-swagger; DO NOT EDIT.

package tokens

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

// V1TokensListReader is a Reader for the V1TokensList structure.
type V1TokensListReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *V1TokensListReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewV1TokensListOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("[GET /v1/tokens/] v1_tokens_list", response, response.Code())
	}
}

// NewV1TokensListOK creates a V1TokensListOK with default headers values
func NewV1TokensListOK() *V1TokensListOK {
	return &V1TokensListOK{}
}

/*
V1TokensListOK describes a response with status code 200, with default header values.

V1TokensListOK v1 tokens list o k
*/
type V1TokensListOK struct {
	Payload *V1TokensListOKBody
}

// IsSuccess returns true when this v1 tokens list o k response has a 2xx status code
func (o *V1TokensListOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this v1 tokens list o k response has a 3xx status code
func (o *V1TokensListOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this v1 tokens list o k response has a 4xx status code
func (o *V1TokensListOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this v1 tokens list o k response has a 5xx status code
func (o *V1TokensListOK) IsServerError() bool {
	return false
}

// IsCode returns true when this v1 tokens list o k response a status code equal to that given
func (o *V1TokensListOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the v1 tokens list o k response
func (o *V1TokensListOK) Code() int {
	return 200
}

func (o *V1TokensListOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /v1/tokens/][%d] v1TokensListOK %s", 200, payload)
}

func (o *V1TokensListOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /v1/tokens/][%d] v1TokensListOK %s", 200, payload)
}

func (o *V1TokensListOK) GetPayload() *V1TokensListOKBody {
	return o.Payload
}

func (o *V1TokensListOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(V1TokensListOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*
V1TokensListOKBody v1 tokens list o k body
swagger:model V1TokensListOKBody
*/
type V1TokensListOKBody struct {

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
	Results []*models.TokenInfoResponse `json:"results"`
}

// Validate validates this v1 tokens list o k body
func (o *V1TokensListOKBody) Validate(formats strfmt.Registry) error {
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

func (o *V1TokensListOKBody) validateCount(formats strfmt.Registry) error {

	if err := validate.Required("v1TokensListOK"+"."+"count", "body", o.Count); err != nil {
		return err
	}

	return nil
}

func (o *V1TokensListOKBody) validateNext(formats strfmt.Registry) error {
	if swag.IsZero(o.Next) { // not required
		return nil
	}

	if err := validate.FormatOf("v1TokensListOK"+"."+"next", "body", "uri", o.Next.String(), formats); err != nil {
		return err
	}

	return nil
}

func (o *V1TokensListOKBody) validatePrevious(formats strfmt.Registry) error {
	if swag.IsZero(o.Previous) { // not required
		return nil
	}

	if err := validate.FormatOf("v1TokensListOK"+"."+"previous", "body", "uri", o.Previous.String(), formats); err != nil {
		return err
	}

	return nil
}

func (o *V1TokensListOKBody) validateResults(formats strfmt.Registry) error {

	if err := validate.Required("v1TokensListOK"+"."+"results", "body", o.Results); err != nil {
		return err
	}

	for i := 0; i < len(o.Results); i++ {
		if swag.IsZero(o.Results[i]) { // not required
			continue
		}

		if o.Results[i] != nil {
			if err := o.Results[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("v1TokensListOK" + "." + "results" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("v1TokensListOK" + "." + "results" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this v1 tokens list o k body based on the context it is used
func (o *V1TokensListOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateResults(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *V1TokensListOKBody) contextValidateResults(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.Results); i++ {

		if o.Results[i] != nil {

			if swag.IsZero(o.Results[i]) { // not required
				return nil
			}

			if err := o.Results[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("v1TokensListOK" + "." + "results" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("v1TokensListOK" + "." + "results" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *V1TokensListOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *V1TokensListOKBody) UnmarshalBinary(b []byte) error {
	var res V1TokensListOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

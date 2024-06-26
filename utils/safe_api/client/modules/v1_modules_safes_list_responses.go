// Code generated by go-swagger; DO NOT EDIT.

package modules

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	models "omni-balance/utils/safe_api"
)

// V1ModulesSafesListReader is a Reader for the V1ModulesSafesList structure.
type V1ModulesSafesListReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *V1ModulesSafesListReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewV1ModulesSafesListOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 422:
		result := NewV1ModulesSafesListUnprocessableEntity()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[GET /v1/modules/{address}/safes/] v1_modules_safes_list", response, response.Code())
	}
}

// NewV1ModulesSafesListOK creates a V1ModulesSafesListOK with default headers values
func NewV1ModulesSafesListOK() *V1ModulesSafesListOK {
	return &V1ModulesSafesListOK{}
}

/*
V1ModulesSafesListOK describes a response with status code 200, with default header values.

V1ModulesSafesListOK v1 modules safes list o k
*/
type V1ModulesSafesListOK struct {
	Payload *models.ModulesResponse
}

// IsSuccess returns true when this v1 modules safes list o k response has a 2xx status code
func (o *V1ModulesSafesListOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this v1 modules safes list o k response has a 3xx status code
func (o *V1ModulesSafesListOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this v1 modules safes list o k response has a 4xx status code
func (o *V1ModulesSafesListOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this v1 modules safes list o k response has a 5xx status code
func (o *V1ModulesSafesListOK) IsServerError() bool {
	return false
}

// IsCode returns true when this v1 modules safes list o k response a status code equal to that given
func (o *V1ModulesSafesListOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the v1 modules safes list o k response
func (o *V1ModulesSafesListOK) Code() int {
	return 200
}

func (o *V1ModulesSafesListOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /v1/modules/{address}/safes/][%d] v1ModulesSafesListOK %s", 200, payload)
}

func (o *V1ModulesSafesListOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /v1/modules/{address}/safes/][%d] v1ModulesSafesListOK %s", 200, payload)
}

func (o *V1ModulesSafesListOK) GetPayload() *models.ModulesResponse {
	return o.Payload
}

func (o *V1ModulesSafesListOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ModulesResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewV1ModulesSafesListUnprocessableEntity creates a V1ModulesSafesListUnprocessableEntity with default headers values
func NewV1ModulesSafesListUnprocessableEntity() *V1ModulesSafesListUnprocessableEntity {
	return &V1ModulesSafesListUnprocessableEntity{}
}

/*
V1ModulesSafesListUnprocessableEntity describes a response with status code 422, with default header values.

Module address checksum not valid
*/
type V1ModulesSafesListUnprocessableEntity struct {
}

// IsSuccess returns true when this v1 modules safes list unprocessable entity response has a 2xx status code
func (o *V1ModulesSafesListUnprocessableEntity) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this v1 modules safes list unprocessable entity response has a 3xx status code
func (o *V1ModulesSafesListUnprocessableEntity) IsRedirect() bool {
	return false
}

// IsClientError returns true when this v1 modules safes list unprocessable entity response has a 4xx status code
func (o *V1ModulesSafesListUnprocessableEntity) IsClientError() bool {
	return true
}

// IsServerError returns true when this v1 modules safes list unprocessable entity response has a 5xx status code
func (o *V1ModulesSafesListUnprocessableEntity) IsServerError() bool {
	return false
}

// IsCode returns true when this v1 modules safes list unprocessable entity response a status code equal to that given
func (o *V1ModulesSafesListUnprocessableEntity) IsCode(code int) bool {
	return code == 422
}

// Code gets the status code for the v1 modules safes list unprocessable entity response
func (o *V1ModulesSafesListUnprocessableEntity) Code() int {
	return 422
}

func (o *V1ModulesSafesListUnprocessableEntity) Error() string {
	return fmt.Sprintf("[GET /v1/modules/{address}/safes/][%d] v1ModulesSafesListUnprocessableEntity", 422)
}

func (o *V1ModulesSafesListUnprocessableEntity) String() string {
	return fmt.Sprintf("[GET /v1/modules/{address}/safes/][%d] v1ModulesSafesListUnprocessableEntity", 422)
}

func (o *V1ModulesSafesListUnprocessableEntity) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

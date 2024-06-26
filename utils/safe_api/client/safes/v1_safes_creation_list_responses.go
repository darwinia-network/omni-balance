// Code generated by go-swagger; DO NOT EDIT.

package safes

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

// V1SafesCreationListReader is a Reader for the V1SafesCreationList structure.
type V1SafesCreationListReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *V1SafesCreationListReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewV1SafesCreationListOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewV1SafesCreationListNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 422:
		result := NewV1SafesCreationListUnprocessableEntity()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 503:
		result := NewV1SafesCreationListServiceUnavailable()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[GET /v1/safes/{address}/creation/] v1_safes_creation_list", response, response.Code())
	}
}

// NewV1SafesCreationListOK creates a V1SafesCreationListOK with default headers values
func NewV1SafesCreationListOK() *V1SafesCreationListOK {
	return &V1SafesCreationListOK{}
}

/*
V1SafesCreationListOK describes a response with status code 200, with default header values.

V1SafesCreationListOK v1 safes creation list o k
*/
type V1SafesCreationListOK struct {
	Payload *models.SafeCreationInfoResponse
}

// IsSuccess returns true when this v1 safes creation list o k response has a 2xx status code
func (o *V1SafesCreationListOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this v1 safes creation list o k response has a 3xx status code
func (o *V1SafesCreationListOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this v1 safes creation list o k response has a 4xx status code
func (o *V1SafesCreationListOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this v1 safes creation list o k response has a 5xx status code
func (o *V1SafesCreationListOK) IsServerError() bool {
	return false
}

// IsCode returns true when this v1 safes creation list o k response a status code equal to that given
func (o *V1SafesCreationListOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the v1 safes creation list o k response
func (o *V1SafesCreationListOK) Code() int {
	return 200
}

func (o *V1SafesCreationListOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /v1/safes/{address}/creation/][%d] v1SafesCreationListOK %s", 200, payload)
}

func (o *V1SafesCreationListOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /v1/safes/{address}/creation/][%d] v1SafesCreationListOK %s", 200, payload)
}

func (o *V1SafesCreationListOK) GetPayload() *models.SafeCreationInfoResponse {
	return o.Payload
}

func (o *V1SafesCreationListOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.SafeCreationInfoResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewV1SafesCreationListNotFound creates a V1SafesCreationListNotFound with default headers values
func NewV1SafesCreationListNotFound() *V1SafesCreationListNotFound {
	return &V1SafesCreationListNotFound{}
}

/*
V1SafesCreationListNotFound describes a response with status code 404, with default header values.

Safe creation not found
*/
type V1SafesCreationListNotFound struct {
}

// IsSuccess returns true when this v1 safes creation list not found response has a 2xx status code
func (o *V1SafesCreationListNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this v1 safes creation list not found response has a 3xx status code
func (o *V1SafesCreationListNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this v1 safes creation list not found response has a 4xx status code
func (o *V1SafesCreationListNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this v1 safes creation list not found response has a 5xx status code
func (o *V1SafesCreationListNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this v1 safes creation list not found response a status code equal to that given
func (o *V1SafesCreationListNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the v1 safes creation list not found response
func (o *V1SafesCreationListNotFound) Code() int {
	return 404
}

func (o *V1SafesCreationListNotFound) Error() string {
	return fmt.Sprintf("[GET /v1/safes/{address}/creation/][%d] v1SafesCreationListNotFound", 404)
}

func (o *V1SafesCreationListNotFound) String() string {
	return fmt.Sprintf("[GET /v1/safes/{address}/creation/][%d] v1SafesCreationListNotFound", 404)
}

func (o *V1SafesCreationListNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewV1SafesCreationListUnprocessableEntity creates a V1SafesCreationListUnprocessableEntity with default headers values
func NewV1SafesCreationListUnprocessableEntity() *V1SafesCreationListUnprocessableEntity {
	return &V1SafesCreationListUnprocessableEntity{}
}

/*
V1SafesCreationListUnprocessableEntity describes a response with status code 422, with default header values.

Owner address checksum not valid
*/
type V1SafesCreationListUnprocessableEntity struct {
}

// IsSuccess returns true when this v1 safes creation list unprocessable entity response has a 2xx status code
func (o *V1SafesCreationListUnprocessableEntity) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this v1 safes creation list unprocessable entity response has a 3xx status code
func (o *V1SafesCreationListUnprocessableEntity) IsRedirect() bool {
	return false
}

// IsClientError returns true when this v1 safes creation list unprocessable entity response has a 4xx status code
func (o *V1SafesCreationListUnprocessableEntity) IsClientError() bool {
	return true
}

// IsServerError returns true when this v1 safes creation list unprocessable entity response has a 5xx status code
func (o *V1SafesCreationListUnprocessableEntity) IsServerError() bool {
	return false
}

// IsCode returns true when this v1 safes creation list unprocessable entity response a status code equal to that given
func (o *V1SafesCreationListUnprocessableEntity) IsCode(code int) bool {
	return code == 422
}

// Code gets the status code for the v1 safes creation list unprocessable entity response
func (o *V1SafesCreationListUnprocessableEntity) Code() int {
	return 422
}

func (o *V1SafesCreationListUnprocessableEntity) Error() string {
	return fmt.Sprintf("[GET /v1/safes/{address}/creation/][%d] v1SafesCreationListUnprocessableEntity", 422)
}

func (o *V1SafesCreationListUnprocessableEntity) String() string {
	return fmt.Sprintf("[GET /v1/safes/{address}/creation/][%d] v1SafesCreationListUnprocessableEntity", 422)
}

func (o *V1SafesCreationListUnprocessableEntity) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewV1SafesCreationListServiceUnavailable creates a V1SafesCreationListServiceUnavailable with default headers values
func NewV1SafesCreationListServiceUnavailable() *V1SafesCreationListServiceUnavailable {
	return &V1SafesCreationListServiceUnavailable{}
}

/*
V1SafesCreationListServiceUnavailable describes a response with status code 503, with default header values.

Problem connecting to Ethereum network
*/
type V1SafesCreationListServiceUnavailable struct {
}

// IsSuccess returns true when this v1 safes creation list service unavailable response has a 2xx status code
func (o *V1SafesCreationListServiceUnavailable) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this v1 safes creation list service unavailable response has a 3xx status code
func (o *V1SafesCreationListServiceUnavailable) IsRedirect() bool {
	return false
}

// IsClientError returns true when this v1 safes creation list service unavailable response has a 4xx status code
func (o *V1SafesCreationListServiceUnavailable) IsClientError() bool {
	return false
}

// IsServerError returns true when this v1 safes creation list service unavailable response has a 5xx status code
func (o *V1SafesCreationListServiceUnavailable) IsServerError() bool {
	return true
}

// IsCode returns true when this v1 safes creation list service unavailable response a status code equal to that given
func (o *V1SafesCreationListServiceUnavailable) IsCode(code int) bool {
	return code == 503
}

// Code gets the status code for the v1 safes creation list service unavailable response
func (o *V1SafesCreationListServiceUnavailable) Code() int {
	return 503
}

func (o *V1SafesCreationListServiceUnavailable) Error() string {
	return fmt.Sprintf("[GET /v1/safes/{address}/creation/][%d] v1SafesCreationListServiceUnavailable", 503)
}

func (o *V1SafesCreationListServiceUnavailable) String() string {
	return fmt.Sprintf("[GET /v1/safes/{address}/creation/][%d] v1SafesCreationListServiceUnavailable", 503)
}

func (o *V1SafesCreationListServiceUnavailable) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

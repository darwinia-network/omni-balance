// Code generated by go-swagger; DO NOT EDIT.

package transactions

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

// V1SafesIncomingTransfersListReader is a Reader for the V1SafesIncomingTransfersList structure.
type V1SafesIncomingTransfersListReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *V1SafesIncomingTransfersListReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewV1SafesIncomingTransfersListOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 422:
		result := NewV1SafesIncomingTransfersListUnprocessableEntity()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[GET /v1/safes/{address}/incoming-transfers/] v1_safes_incoming-transfers_list", response, response.Code())
	}
}

// NewV1SafesIncomingTransfersListOK creates a V1SafesIncomingTransfersListOK with default headers values
func NewV1SafesIncomingTransfersListOK() *V1SafesIncomingTransfersListOK {
	return &V1SafesIncomingTransfersListOK{}
}

/*
V1SafesIncomingTransfersListOK describes a response with status code 200, with default header values.

V1SafesIncomingTransfersListOK v1 safes incoming transfers list o k
*/
type V1SafesIncomingTransfersListOK struct {
	Payload []*models.TransferWithTokenInfoResponse
}

// IsSuccess returns true when this v1 safes incoming transfers list o k response has a 2xx status code
func (o *V1SafesIncomingTransfersListOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this v1 safes incoming transfers list o k response has a 3xx status code
func (o *V1SafesIncomingTransfersListOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this v1 safes incoming transfers list o k response has a 4xx status code
func (o *V1SafesIncomingTransfersListOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this v1 safes incoming transfers list o k response has a 5xx status code
func (o *V1SafesIncomingTransfersListOK) IsServerError() bool {
	return false
}

// IsCode returns true when this v1 safes incoming transfers list o k response a status code equal to that given
func (o *V1SafesIncomingTransfersListOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the v1 safes incoming transfers list o k response
func (o *V1SafesIncomingTransfersListOK) Code() int {
	return 200
}

func (o *V1SafesIncomingTransfersListOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /v1/safes/{address}/incoming-transfers/][%d] v1SafesIncomingTransfersListOK %s", 200, payload)
}

func (o *V1SafesIncomingTransfersListOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /v1/safes/{address}/incoming-transfers/][%d] v1SafesIncomingTransfersListOK %s", 200, payload)
}

func (o *V1SafesIncomingTransfersListOK) GetPayload() []*models.TransferWithTokenInfoResponse {
	return o.Payload
}

func (o *V1SafesIncomingTransfersListOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewV1SafesIncomingTransfersListUnprocessableEntity creates a V1SafesIncomingTransfersListUnprocessableEntity with default headers values
func NewV1SafesIncomingTransfersListUnprocessableEntity() *V1SafesIncomingTransfersListUnprocessableEntity {
	return &V1SafesIncomingTransfersListUnprocessableEntity{}
}

/*
V1SafesIncomingTransfersListUnprocessableEntity describes a response with status code 422, with default header values.

Safe address checksum not valid
*/
type V1SafesIncomingTransfersListUnprocessableEntity struct {
}

// IsSuccess returns true when this v1 safes incoming transfers list unprocessable entity response has a 2xx status code
func (o *V1SafesIncomingTransfersListUnprocessableEntity) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this v1 safes incoming transfers list unprocessable entity response has a 3xx status code
func (o *V1SafesIncomingTransfersListUnprocessableEntity) IsRedirect() bool {
	return false
}

// IsClientError returns true when this v1 safes incoming transfers list unprocessable entity response has a 4xx status code
func (o *V1SafesIncomingTransfersListUnprocessableEntity) IsClientError() bool {
	return true
}

// IsServerError returns true when this v1 safes incoming transfers list unprocessable entity response has a 5xx status code
func (o *V1SafesIncomingTransfersListUnprocessableEntity) IsServerError() bool {
	return false
}

// IsCode returns true when this v1 safes incoming transfers list unprocessable entity response a status code equal to that given
func (o *V1SafesIncomingTransfersListUnprocessableEntity) IsCode(code int) bool {
	return code == 422
}

// Code gets the status code for the v1 safes incoming transfers list unprocessable entity response
func (o *V1SafesIncomingTransfersListUnprocessableEntity) Code() int {
	return 422
}

func (o *V1SafesIncomingTransfersListUnprocessableEntity) Error() string {
	return fmt.Sprintf("[GET /v1/safes/{address}/incoming-transfers/][%d] v1SafesIncomingTransfersListUnprocessableEntity", 422)
}

func (o *V1SafesIncomingTransfersListUnprocessableEntity) String() string {
	return fmt.Sprintf("[GET /v1/safes/{address}/incoming-transfers/][%d] v1SafesIncomingTransfersListUnprocessableEntity", 422)
}

func (o *V1SafesIncomingTransfersListUnprocessableEntity) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

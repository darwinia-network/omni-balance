// Code generated by go-swagger; DO NOT EDIT.

package messages

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

// V1MessagesReadReader is a Reader for the V1MessagesRead structure.
type V1MessagesReadReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *V1MessagesReadReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewV1MessagesReadOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("[GET /v1/messages/{message_hash}/] v1_messages_read", response, response.Code())
	}
}

// NewV1MessagesReadOK creates a V1MessagesReadOK with default headers values
func NewV1MessagesReadOK() *V1MessagesReadOK {
	return &V1MessagesReadOK{}
}

/*
V1MessagesReadOK describes a response with status code 200, with default header values.

V1MessagesReadOK v1 messages read o k
*/
type V1MessagesReadOK struct {
	Payload *models.SafeMessageResponse
}

// IsSuccess returns true when this v1 messages read o k response has a 2xx status code
func (o *V1MessagesReadOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this v1 messages read o k response has a 3xx status code
func (o *V1MessagesReadOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this v1 messages read o k response has a 4xx status code
func (o *V1MessagesReadOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this v1 messages read o k response has a 5xx status code
func (o *V1MessagesReadOK) IsServerError() bool {
	return false
}

// IsCode returns true when this v1 messages read o k response a status code equal to that given
func (o *V1MessagesReadOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the v1 messages read o k response
func (o *V1MessagesReadOK) Code() int {
	return 200
}

func (o *V1MessagesReadOK) Error() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /v1/messages/{message_hash}/][%d] v1MessagesReadOK %s", 200, payload)
}

func (o *V1MessagesReadOK) String() string {
	payload, _ := json.Marshal(o.Payload)
	return fmt.Sprintf("[GET /v1/messages/{message_hash}/][%d] v1MessagesReadOK %s", 200, payload)
}

func (o *V1MessagesReadOK) GetPayload() *models.SafeMessageResponse {
	return o.Payload
}

func (o *V1MessagesReadOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.SafeMessageResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

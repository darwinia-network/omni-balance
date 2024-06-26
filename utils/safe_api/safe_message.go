// Code generated by go-swagger; DO NOT EDIT.

package safe_api

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// SafeMessage safe message
//
// swagger:model SafeMessage
type SafeMessage struct {

	// Message
	// Required: true
	Message interface{} `json:"message"`

	// Safe app id
	// Minimum: 0
	SafeAppID *int64 `json:"safeAppId,omitempty"`

	// Signature
	// Required: true
	Signature *string `json:"signature"`
}

// Validate validates this safe message
func (m *SafeMessage) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateMessage(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSafeAppID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSignature(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *SafeMessage) validateMessage(formats strfmt.Registry) error {

	if m.Message == nil {
		return errors.Required("message", "body", nil)
	}

	return nil
}

func (m *SafeMessage) validateSafeAppID(formats strfmt.Registry) error {
	if swag.IsZero(m.SafeAppID) { // not required
		return nil
	}

	if err := validate.MinimumInt("safeAppId", "body", *m.SafeAppID, 0, false); err != nil {
		return err
	}

	return nil
}

func (m *SafeMessage) validateSignature(formats strfmt.Registry) error {

	if err := validate.Required("signature", "body", m.Signature); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this safe message based on context it is used
func (m *SafeMessage) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *SafeMessage) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SafeMessage) UnmarshalBinary(b []byte) error {
	var res SafeMessage
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

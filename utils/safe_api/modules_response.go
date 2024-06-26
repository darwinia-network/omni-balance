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

// ModulesResponse modules response
//
// swagger:model ModulesResponse
type ModulesResponse struct {

	// safes
	// Required: true
	Safes []string `json:"safes"`
}

// Validate validates this modules response
func (m *ModulesResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateSafes(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ModulesResponse) validateSafes(formats strfmt.Registry) error {

	if err := validate.Required("safes", "body", m.Safes); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this modules response based on context it is used
func (m *ModulesResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ModulesResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ModulesResponse) UnmarshalBinary(b []byte) error {
	var res ModulesResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

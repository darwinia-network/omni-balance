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

// DelegateDelete delegate delete
//
// swagger:model DelegateDelete
type DelegateDelete struct {

	// Delegate
	// Required: true
	Delegate *string `json:"delegate"`

	// Delegator
	// Required: true
	Delegator *string `json:"delegator"`

	// Signature
	// Required: true
	Signature *string `json:"signature"`
}

// Validate validates this delegate delete
func (m *DelegateDelete) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDelegate(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDelegator(formats); err != nil {
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

func (m *DelegateDelete) validateDelegate(formats strfmt.Registry) error {

	if err := validate.Required("delegate", "body", m.Delegate); err != nil {
		return err
	}

	return nil
}

func (m *DelegateDelete) validateDelegator(formats strfmt.Registry) error {

	if err := validate.Required("delegator", "body", m.Delegator); err != nil {
		return err
	}

	return nil
}

func (m *DelegateDelete) validateSignature(formats strfmt.Registry) error {

	if err := validate.Required("signature", "body", m.Signature); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this delegate delete based on context it is used
func (m *DelegateDelete) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *DelegateDelete) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DelegateDelete) UnmarshalBinary(b []byte) error {
	var res DelegateDelete
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

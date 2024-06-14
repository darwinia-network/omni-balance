// Code generated by go-swagger; DO NOT EDIT.

package safe_api

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// SafeDeployment safe deployment
//
// swagger:model SafeDeployment
type SafeDeployment struct {

	// contracts
	// Required: true
	Contracts []*SafeDeploymentContract `json:"contracts"`

	// Version
	// Required: true
	// Max Length: 10
	// Min Length: 1
	Version *string `json:"version"`
}

// Validate validates this safe deployment
func (m *SafeDeployment) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateContracts(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateVersion(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *SafeDeployment) validateContracts(formats strfmt.Registry) error {

	if err := validate.Required("contracts", "body", m.Contracts); err != nil {
		return err
	}

	for i := 0; i < len(m.Contracts); i++ {
		if swag.IsZero(m.Contracts[i]) { // not required
			continue
		}

		if m.Contracts[i] != nil {
			if err := m.Contracts[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("contracts" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("contracts" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *SafeDeployment) validateVersion(formats strfmt.Registry) error {

	if err := validate.Required("version", "body", m.Version); err != nil {
		return err
	}

	if err := validate.MinLength("version", "body", *m.Version, 1); err != nil {
		return err
	}

	if err := validate.MaxLength("version", "body", *m.Version, 10); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this safe deployment based on the context it is used
func (m *SafeDeployment) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateContracts(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *SafeDeployment) contextValidateContracts(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Contracts); i++ {

		if m.Contracts[i] != nil {

			if swag.IsZero(m.Contracts[i]) { // not required
				return nil
			}

			if err := m.Contracts[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("contracts" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("contracts" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *SafeDeployment) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SafeDeployment) UnmarshalBinary(b []byte) error {
	var res SafeDeployment
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

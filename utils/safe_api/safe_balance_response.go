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

// SafeBalanceResponse safe balance response
//
// swagger:model SafeBalanceResponse
type SafeBalanceResponse struct {

	// Balance
	// Required: true
	// Min Length: 1
	Balance *string `json:"balance"`

	// token
	// Required: true
	Token *Erc20Info `json:"token"`

	// Token address
	// Required: true
	// Min Length: 1
	TokenAddress *string `json:"tokenAddress"`
}

// Validate validates this safe balance response
func (m *SafeBalanceResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateBalance(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateToken(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTokenAddress(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *SafeBalanceResponse) validateBalance(formats strfmt.Registry) error {

	if err := validate.Required("balance", "body", m.Balance); err != nil {
		return err
	}

	if err := validate.MinLength("balance", "body", *m.Balance, 1); err != nil {
		return err
	}

	return nil
}

func (m *SafeBalanceResponse) validateToken(formats strfmt.Registry) error {

	if err := validate.Required("token", "body", m.Token); err != nil {
		return err
	}

	if m.Token != nil {
		if err := m.Token.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("token")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("token")
			}
			return err
		}
	}

	return nil
}

func (m *SafeBalanceResponse) validateTokenAddress(formats strfmt.Registry) error {

	if err := validate.Required("tokenAddress", "body", m.TokenAddress); err != nil {
		return err
	}

	if err := validate.MinLength("tokenAddress", "body", *m.TokenAddress, 1); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this safe balance response based on the context it is used
func (m *SafeBalanceResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateToken(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *SafeBalanceResponse) contextValidateToken(ctx context.Context, formats strfmt.Registry) error {

	if m.Token != nil {

		if err := m.Token.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("token")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("token")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *SafeBalanceResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SafeBalanceResponse) UnmarshalBinary(b []byte) error {
	var res SafeBalanceResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

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

// SafeCreationInfoResponse safe creation info response
//
// swagger:model SafeCreationInfoResponse
type SafeCreationInfoResponse struct {

	// Created
	// Required: true
	// Format: date-time
	Created *strfmt.DateTime `json:"created"`

	// Creator
	// Required: true
	Creator *string `json:"creator"`

	// Data decoded
	// Read Only: true
	DataDecoded string `json:"dataDecoded,omitempty"`

	// Factory address
	// Required: true
	FactoryAddress *string `json:"factoryAddress"`

	// Master copy
	// Required: true
	MasterCopy *string `json:"masterCopy"`

	// Setup data
	// Required: true
	SetupData *string `json:"setupData"`

	// Transaction hash
	// Required: true
	TransactionHash *string `json:"transactionHash"`

	// user operation
	// Required: true
	UserOperation *UserOperationWithSafeOperationResponse `json:"userOperation"`
}

// Validate validates this safe creation info response
func (m *SafeCreationInfoResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCreated(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCreator(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFactoryAddress(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMasterCopy(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSetupData(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTransactionHash(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUserOperation(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *SafeCreationInfoResponse) validateCreated(formats strfmt.Registry) error {

	if err := validate.Required("created", "body", m.Created); err != nil {
		return err
	}

	if err := validate.FormatOf("created", "body", "date-time", m.Created.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *SafeCreationInfoResponse) validateCreator(formats strfmt.Registry) error {

	if err := validate.Required("creator", "body", m.Creator); err != nil {
		return err
	}

	return nil
}

func (m *SafeCreationInfoResponse) validateFactoryAddress(formats strfmt.Registry) error {

	if err := validate.Required("factoryAddress", "body", m.FactoryAddress); err != nil {
		return err
	}

	return nil
}

func (m *SafeCreationInfoResponse) validateMasterCopy(formats strfmt.Registry) error {

	if err := validate.Required("masterCopy", "body", m.MasterCopy); err != nil {
		return err
	}

	return nil
}

func (m *SafeCreationInfoResponse) validateSetupData(formats strfmt.Registry) error {

	if err := validate.Required("setupData", "body", m.SetupData); err != nil {
		return err
	}

	return nil
}

func (m *SafeCreationInfoResponse) validateTransactionHash(formats strfmt.Registry) error {

	if err := validate.Required("transactionHash", "body", m.TransactionHash); err != nil {
		return err
	}

	return nil
}

func (m *SafeCreationInfoResponse) validateUserOperation(formats strfmt.Registry) error {

	if err := validate.Required("userOperation", "body", m.UserOperation); err != nil {
		return err
	}

	if m.UserOperation != nil {
		if err := m.UserOperation.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("userOperation")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("userOperation")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this safe creation info response based on the context it is used
func (m *SafeCreationInfoResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateDataDecoded(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateUserOperation(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *SafeCreationInfoResponse) contextValidateDataDecoded(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "dataDecoded", "body", string(m.DataDecoded)); err != nil {
		return err
	}

	return nil
}

func (m *SafeCreationInfoResponse) contextValidateUserOperation(ctx context.Context, formats strfmt.Registry) error {

	if m.UserOperation != nil {

		if err := m.UserOperation.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("userOperation")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("userOperation")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *SafeCreationInfoResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SafeCreationInfoResponse) UnmarshalBinary(b []byte) error {
	var res SafeCreationInfoResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

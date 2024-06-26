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

// SafeOperation safe operation
//
// swagger:model SafeOperation
type SafeOperation struct {

	// Call data
	// Required: true
	CallData *string `json:"callData"`

	// Call gas limit
	// Required: true
	// Minimum: 0
	CallGasLimit *int64 `json:"callGasLimit"`

	// Entry point
	// Required: true
	EntryPoint *string `json:"entryPoint"`

	// Init code
	// Required: true
	InitCode *string `json:"initCode"`

	// Max fee per gas
	// Required: true
	// Minimum: 0
	MaxFeePerGas *int64 `json:"maxFeePerGas"`

	// Max priority fee per gas
	// Required: true
	// Minimum: 0
	MaxPriorityFeePerGas *int64 `json:"maxPriorityFeePerGas"`

	// Module address
	// Required: true
	ModuleAddress *string `json:"moduleAddress"`

	// Nonce
	// Required: true
	// Minimum: 0
	Nonce *int64 `json:"nonce"`

	// Paymaster and data
	// Required: true
	PaymasterAndData *string `json:"paymasterAndData"`

	// Pre verification gas
	// Required: true
	// Minimum: 0
	PreVerificationGas *int64 `json:"preVerificationGas"`

	// Signature
	// Required: true
	Signature *string `json:"signature"`

	// Valid after
	// Required: true
	// Format: date-time
	ValidAfter *strfmt.DateTime `json:"validAfter"`

	// Valid until
	// Required: true
	// Format: date-time
	ValidUntil *strfmt.DateTime `json:"validUntil"`

	// Verification gas limit
	// Required: true
	// Minimum: 0
	VerificationGasLimit *int64 `json:"verificationGasLimit"`
}

// Validate validates this safe operation
func (m *SafeOperation) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCallData(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCallGasLimit(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateEntryPoint(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateInitCode(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMaxFeePerGas(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMaxPriorityFeePerGas(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateModuleAddress(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNonce(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePaymasterAndData(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePreVerificationGas(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSignature(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateValidAfter(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateValidUntil(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateVerificationGasLimit(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *SafeOperation) validateCallData(formats strfmt.Registry) error {

	if err := validate.Required("callData", "body", m.CallData); err != nil {
		return err
	}

	return nil
}

func (m *SafeOperation) validateCallGasLimit(formats strfmt.Registry) error {

	if err := validate.Required("callGasLimit", "body", m.CallGasLimit); err != nil {
		return err
	}

	if err := validate.MinimumInt("callGasLimit", "body", *m.CallGasLimit, 0, false); err != nil {
		return err
	}

	return nil
}

func (m *SafeOperation) validateEntryPoint(formats strfmt.Registry) error {

	if err := validate.Required("entryPoint", "body", m.EntryPoint); err != nil {
		return err
	}

	return nil
}

func (m *SafeOperation) validateInitCode(formats strfmt.Registry) error {

	if err := validate.Required("initCode", "body", m.InitCode); err != nil {
		return err
	}

	return nil
}

func (m *SafeOperation) validateMaxFeePerGas(formats strfmt.Registry) error {

	if err := validate.Required("maxFeePerGas", "body", m.MaxFeePerGas); err != nil {
		return err
	}

	if err := validate.MinimumInt("maxFeePerGas", "body", *m.MaxFeePerGas, 0, false); err != nil {
		return err
	}

	return nil
}

func (m *SafeOperation) validateMaxPriorityFeePerGas(formats strfmt.Registry) error {

	if err := validate.Required("maxPriorityFeePerGas", "body", m.MaxPriorityFeePerGas); err != nil {
		return err
	}

	if err := validate.MinimumInt("maxPriorityFeePerGas", "body", *m.MaxPriorityFeePerGas, 0, false); err != nil {
		return err
	}

	return nil
}

func (m *SafeOperation) validateModuleAddress(formats strfmt.Registry) error {

	if err := validate.Required("moduleAddress", "body", m.ModuleAddress); err != nil {
		return err
	}

	return nil
}

func (m *SafeOperation) validateNonce(formats strfmt.Registry) error {

	if err := validate.Required("nonce", "body", m.Nonce); err != nil {
		return err
	}

	if err := validate.MinimumInt("nonce", "body", *m.Nonce, 0, false); err != nil {
		return err
	}

	return nil
}

func (m *SafeOperation) validatePaymasterAndData(formats strfmt.Registry) error {

	if err := validate.Required("paymasterAndData", "body", m.PaymasterAndData); err != nil {
		return err
	}

	return nil
}

func (m *SafeOperation) validatePreVerificationGas(formats strfmt.Registry) error {

	if err := validate.Required("preVerificationGas", "body", m.PreVerificationGas); err != nil {
		return err
	}

	if err := validate.MinimumInt("preVerificationGas", "body", *m.PreVerificationGas, 0, false); err != nil {
		return err
	}

	return nil
}

func (m *SafeOperation) validateSignature(formats strfmt.Registry) error {

	if err := validate.Required("signature", "body", m.Signature); err != nil {
		return err
	}

	return nil
}

func (m *SafeOperation) validateValidAfter(formats strfmt.Registry) error {

	if err := validate.Required("validAfter", "body", m.ValidAfter); err != nil {
		return err
	}

	if err := validate.FormatOf("validAfter", "body", "date-time", m.ValidAfter.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *SafeOperation) validateValidUntil(formats strfmt.Registry) error {

	if err := validate.Required("validUntil", "body", m.ValidUntil); err != nil {
		return err
	}

	if err := validate.FormatOf("validUntil", "body", "date-time", m.ValidUntil.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *SafeOperation) validateVerificationGasLimit(formats strfmt.Registry) error {

	if err := validate.Required("verificationGasLimit", "body", m.VerificationGasLimit); err != nil {
		return err
	}

	if err := validate.MinimumInt("verificationGasLimit", "body", *m.VerificationGasLimit, 0, false); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this safe operation based on context it is used
func (m *SafeOperation) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *SafeOperation) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SafeOperation) UnmarshalBinary(b []byte) error {
	var res SafeOperation
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

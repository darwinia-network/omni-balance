// Code generated by go-swagger; DO NOT EDIT.

package safe_api

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// SafeModuleTransactionResponse safe module transaction response
//
// swagger:model SafeModuleTransactionResponse
type SafeModuleTransactionResponse struct {

	// Block number
	// Read Only: true
	BlockNumber int64 `json:"blockNumber,omitempty"`

	// Created
	// Read Only: true
	// Format: date-time
	Created strfmt.DateTime `json:"created,omitempty"`

	// Data
	// Required: true
	Data *string `json:"data"`

	// Data decoded
	// Read Only: true
	DataDecoded string `json:"dataDecoded,omitempty"`

	// Execution date
	// Required: true
	// Format: date-time
	ExecutionDate *strfmt.DateTime `json:"executionDate"`

	// Is successful
	// Read Only: true
	IsSuccessful *bool `json:"isSuccessful,omitempty"`

	// Module
	// Required: true
	// Min Length: 1
	Module *string `json:"module"`

	// Module transaction id
	//
	// Internally calculated parameter to uniquely identify a moduleTransaction
	// `ModuleTransactionId = i+tx_hash+trace_address`
	// Read Only: true
	ModuleTransactionID string `json:"moduleTransactionId,omitempty"`

	// Operation
	// Required: true
	// Enum: [0,1,2]
	Operation *int64 `json:"operation"`

	// Safe
	// Required: true
	// Min Length: 1
	Safe *string `json:"safe"`

	// To
	// Required: true
	// Min Length: 1
	To *string `json:"to"`

	// Transaction hash
	// Read Only: true
	TransactionHash string `json:"transactionHash,omitempty"`

	// Value
	// Required: true
	Value *string `json:"value"`
}

// Validate validates this safe module transaction response
func (m *SafeModuleTransactionResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCreated(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateData(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateExecutionDate(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateModule(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateOperation(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSafe(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTo(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateValue(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *SafeModuleTransactionResponse) validateCreated(formats strfmt.Registry) error {
	if swag.IsZero(m.Created) { // not required
		return nil
	}

	if err := validate.FormatOf("created", "body", "date-time", m.Created.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *SafeModuleTransactionResponse) validateData(formats strfmt.Registry) error {

	if err := validate.Required("data", "body", m.Data); err != nil {
		return err
	}

	return nil
}

func (m *SafeModuleTransactionResponse) validateExecutionDate(formats strfmt.Registry) error {

	if err := validate.Required("executionDate", "body", m.ExecutionDate); err != nil {
		return err
	}

	if err := validate.FormatOf("executionDate", "body", "date-time", m.ExecutionDate.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *SafeModuleTransactionResponse) validateModule(formats strfmt.Registry) error {

	if err := validate.Required("module", "body", m.Module); err != nil {
		return err
	}

	if err := validate.MinLength("module", "body", *m.Module, 1); err != nil {
		return err
	}

	return nil
}

var safeModuleTransactionResponseTypeOperationPropEnum []interface{}

func init() {
	var res []int64
	if err := json.Unmarshal([]byte(`[0,1,2]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		safeModuleTransactionResponseTypeOperationPropEnum = append(safeModuleTransactionResponseTypeOperationPropEnum, v)
	}
}

// prop value enum
func (m *SafeModuleTransactionResponse) validateOperationEnum(path, location string, value int64) error {
	if err := validate.EnumCase(path, location, value, safeModuleTransactionResponseTypeOperationPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *SafeModuleTransactionResponse) validateOperation(formats strfmt.Registry) error {

	if err := validate.Required("operation", "body", m.Operation); err != nil {
		return err
	}

	// value enum
	if err := m.validateOperationEnum("operation", "body", *m.Operation); err != nil {
		return err
	}

	return nil
}

func (m *SafeModuleTransactionResponse) validateSafe(formats strfmt.Registry) error {

	if err := validate.Required("safe", "body", m.Safe); err != nil {
		return err
	}

	if err := validate.MinLength("safe", "body", *m.Safe, 1); err != nil {
		return err
	}

	return nil
}

func (m *SafeModuleTransactionResponse) validateTo(formats strfmt.Registry) error {

	if err := validate.Required("to", "body", m.To); err != nil {
		return err
	}

	if err := validate.MinLength("to", "body", *m.To, 1); err != nil {
		return err
	}

	return nil
}

func (m *SafeModuleTransactionResponse) validateValue(formats strfmt.Registry) error {

	if err := validate.Required("value", "body", m.Value); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this safe module transaction response based on the context it is used
func (m *SafeModuleTransactionResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateBlockNumber(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateCreated(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateDataDecoded(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateIsSuccessful(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateModuleTransactionID(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateTransactionHash(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *SafeModuleTransactionResponse) contextValidateBlockNumber(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "blockNumber", "body", int64(m.BlockNumber)); err != nil {
		return err
	}

	return nil
}

func (m *SafeModuleTransactionResponse) contextValidateCreated(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "created", "body", strfmt.DateTime(m.Created)); err != nil {
		return err
	}

	return nil
}

func (m *SafeModuleTransactionResponse) contextValidateDataDecoded(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "dataDecoded", "body", string(m.DataDecoded)); err != nil {
		return err
	}

	return nil
}

func (m *SafeModuleTransactionResponse) contextValidateIsSuccessful(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "isSuccessful", "body", m.IsSuccessful); err != nil {
		return err
	}

	return nil
}

func (m *SafeModuleTransactionResponse) contextValidateModuleTransactionID(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "moduleTransactionId", "body", string(m.ModuleTransactionID)); err != nil {
		return err
	}

	return nil
}

func (m *SafeModuleTransactionResponse) contextValidateTransactionHash(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.ReadOnly(ctx, "transactionHash", "body", string(m.TransactionHash)); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *SafeModuleTransactionResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SafeModuleTransactionResponse) UnmarshalBinary(b []byte) error {
	var res SafeModuleTransactionResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

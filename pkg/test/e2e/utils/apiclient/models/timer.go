// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Timer Represents all available timers in a vmi.
//
// swagger:model Timer
type Timer struct {

	// hpet
	Hpet *HPETTimer `json:"hpet,omitempty"`

	// hyperv
	Hyperv *HypervTimer `json:"hyperv,omitempty"`

	// kvm
	Kvm *KVMTimer `json:"kvm,omitempty"`

	// pit
	Pit *PITTimer `json:"pit,omitempty"`

	// rtc
	Rtc *RTCTimer `json:"rtc,omitempty"`
}

// Validate validates this timer
func (m *Timer) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateHpet(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateHyperv(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateKvm(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePit(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRtc(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Timer) validateHpet(formats strfmt.Registry) error {
	if swag.IsZero(m.Hpet) { // not required
		return nil
	}

	if m.Hpet != nil {
		if err := m.Hpet.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("hpet")
			}
			return err
		}
	}

	return nil
}

func (m *Timer) validateHyperv(formats strfmt.Registry) error {
	if swag.IsZero(m.Hyperv) { // not required
		return nil
	}

	if m.Hyperv != nil {
		if err := m.Hyperv.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("hyperv")
			}
			return err
		}
	}

	return nil
}

func (m *Timer) validateKvm(formats strfmt.Registry) error {
	if swag.IsZero(m.Kvm) { // not required
		return nil
	}

	if m.Kvm != nil {
		if err := m.Kvm.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("kvm")
			}
			return err
		}
	}

	return nil
}

func (m *Timer) validatePit(formats strfmt.Registry) error {
	if swag.IsZero(m.Pit) { // not required
		return nil
	}

	if m.Pit != nil {
		if err := m.Pit.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("pit")
			}
			return err
		}
	}

	return nil
}

func (m *Timer) validateRtc(formats strfmt.Registry) error {
	if swag.IsZero(m.Rtc) { // not required
		return nil
	}

	if m.Rtc != nil {
		if err := m.Rtc.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("rtc")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this timer based on the context it is used
func (m *Timer) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateHpet(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateHyperv(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateKvm(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidatePit(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateRtc(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Timer) contextValidateHpet(ctx context.Context, formats strfmt.Registry) error {

	if m.Hpet != nil {
		if err := m.Hpet.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("hpet")
			}
			return err
		}
	}

	return nil
}

func (m *Timer) contextValidateHyperv(ctx context.Context, formats strfmt.Registry) error {

	if m.Hyperv != nil {
		if err := m.Hyperv.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("hyperv")
			}
			return err
		}
	}

	return nil
}

func (m *Timer) contextValidateKvm(ctx context.Context, formats strfmt.Registry) error {

	if m.Kvm != nil {
		if err := m.Kvm.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("kvm")
			}
			return err
		}
	}

	return nil
}

func (m *Timer) contextValidatePit(ctx context.Context, formats strfmt.Registry) error {

	if m.Pit != nil {
		if err := m.Pit.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("pit")
			}
			return err
		}
	}

	return nil
}

func (m *Timer) contextValidateRtc(ctx context.Context, formats strfmt.Registry) error {

	if m.Rtc != nil {
		if err := m.Rtc.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("rtc")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Timer) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Timer) UnmarshalBinary(b []byte) error {
	var res Timer
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
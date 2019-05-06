package valueobjects

import (
	"encoding/json"
	"go-iddd/shared"

	"golang.org/x/xerrors"
)

type ConfirmableEmailAddress struct {
	baseEmailAddress *EmailAddress
	confirmationHash *ConfirmationHash
	isConfirmed      bool
}

/*** Factory methods ***/

func NewConfirmableEmailAddress(from string) (*ConfirmableEmailAddress, error) {
	baseEmailAddress, err := NewEmailAddress(from)
	if err != nil {
		return nil, xerrors.Errorf("NewConfirmableEmailAddress: %w", err)
	}

	newEmailAddress := buildConfirmableEmailAddress(baseEmailAddress, GenerateConfirmationHash(from))

	return newEmailAddress, nil
}

func ReconstituteConfirmableEmailAddress(from string, withConfirmationHash string) *ConfirmableEmailAddress {
	return buildConfirmableEmailAddress(
		ReconstituteEmailAddress(from),
		ReconstituteConfirmationHash(withConfirmationHash),
	)
}

func buildConfirmableEmailAddress(from *EmailAddress, with *ConfirmationHash) *ConfirmableEmailAddress {
	return &ConfirmableEmailAddress{
		baseEmailAddress: from,
		confirmationHash: with,
		isConfirmed:      false,
	}
}

/*** Public methods implementing ConfirmableEmailAddress ***/

func (confirmableEmailAddress *ConfirmableEmailAddress) Confirm(given *ConfirmationHash) (*ConfirmableEmailAddress, error) {
	if err := confirmableEmailAddress.confirmationHash.MustMatch(given); err != nil {
		return nil, xerrors.Errorf("confirmableEmailAddress.Confirm: %w", err)
	}

	confirmedEmailAddress := buildConfirmableEmailAddress(
		confirmableEmailAddress.baseEmailAddress,
		confirmableEmailAddress.confirmationHash,
	)

	confirmedEmailAddress.isConfirmed = true

	return confirmedEmailAddress, nil
}

func (confirmableEmailAddress *ConfirmableEmailAddress) IsConfirmed() bool {
	return confirmableEmailAddress.isConfirmed
}

func (confirmableEmailAddress *ConfirmableEmailAddress) EmailAddress() string {
	return confirmableEmailAddress.baseEmailAddress.EmailAddress()
}

func (confirmableEmailAddress *ConfirmableEmailAddress) Equals(other *ConfirmableEmailAddress) bool {
	return confirmableEmailAddress.baseEmailAddress.Equals(other.baseEmailAddress)
}

func (confirmableEmailAddress *ConfirmableEmailAddress) EqualsAny(other *EmailAddress) bool {
	return confirmableEmailAddress.baseEmailAddress.Equals(other)
}

func (confirmableEmailAddress *ConfirmableEmailAddress) MarshalJSON() ([]byte, error) {
	data := &struct {
		EmailAddress     *EmailAddress     `json:"emailAddress"`
		ConfirmationHash *ConfirmationHash `json:"confirmationHash"`
	}{
		EmailAddress:     confirmableEmailAddress.baseEmailAddress,
		ConfirmationHash: confirmableEmailAddress.confirmationHash,
	}

	return json.Marshal(data)
}

func UnmarshalConfirmableEmailAddress(input interface{}) (*ConfirmableEmailAddress, error) {
	var err error

	values, ok := input.(map[string]interface{})
	if !ok {
		return nil,
			xerrors.Errorf(
				"UnmarshalConfirmableEmailAddress: input is not [map[string]interface{}]: %w",
				shared.ErrUnmarshaling,
			)
	}

	confirmableEmailAddress := &ConfirmableEmailAddress{}

	for key, value := range values {
		switch key {
		case "emailAddress":
			if confirmableEmailAddress.baseEmailAddress, err = UnmarshalEmailAddress(value); err != nil {
				return nil, xerrors.Errorf("UnmarshalConfirmableEmailAddress: %w", err)
			}
		case "confirmationHash":
			if confirmableEmailAddress.confirmationHash, err = UnmarshalConfirmationHash(value); err != nil {
				return nil, xerrors.Errorf("UnmarshalConfirmableEmailAddress: %w", err)
			}
		}
	}

	return confirmableEmailAddress, nil
}

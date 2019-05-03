package valueobjects

import (
	"encoding/json"
	"errors"
)

type ConfirmableEmailAddress interface {
	Confirm(given ConfirmationHash) (*confirmableEmailAddress, error)
	IsConfirmed() bool

	EmailAddress
}

type confirmableEmailAddress struct {
	baseEmailAddress EmailAddress
	confirmationHash ConfirmationHash
	isConfirmed      bool
}

/*** Factory methods ***/

func NewConfirmableEmailAddress(from string) (*confirmableEmailAddress, error) {
	baseEmailAddress, err := NewEmailAddress(from)
	if err != nil {
		// TODO: map error?
		return nil, err
	}

	newEmailAddress := buildConfirmableEmailAddress(baseEmailAddress, GenerateConfirmationHash(from))

	return newEmailAddress, nil
}

func ReconstituteConfirmableEmailAddress(from string, withConfirmationHash string) *confirmableEmailAddress {
	return buildConfirmableEmailAddress(
		ReconstituteEmailAddress(from),
		ReconstituteConfirmationHash(withConfirmationHash),
	)
}

func buildConfirmableEmailAddress(from EmailAddress, with ConfirmationHash) *confirmableEmailAddress {
	return &confirmableEmailAddress{
		baseEmailAddress: from,
		confirmationHash: with,
		isConfirmed:      false,
	}
}

/*** Public methods implementing ConfirmableEmailAddress (own methods) ***/

func (confirmableEmailAddress *confirmableEmailAddress) Confirm(given ConfirmationHash) (*confirmableEmailAddress, error) {
	if err := confirmableEmailAddress.confirmationHash.MustMatch(given); err != nil {
		return nil, err
	}

	confirmedEmailAddress := buildConfirmableEmailAddress(
		confirmableEmailAddress.baseEmailAddress,
		confirmableEmailAddress.confirmationHash,
	)

	confirmedEmailAddress.isConfirmed = true

	return confirmedEmailAddress, nil
}

func (confirmableEmailAddress *confirmableEmailAddress) IsConfirmed() bool {
	return confirmableEmailAddress.isConfirmed
}

/*** Public methods implementing ConfirmableEmailAddress (methods for EmailAddress) ***/

func (confirmableEmailAddress *confirmableEmailAddress) EmailAddress() string {
	return confirmableEmailAddress.baseEmailAddress.EmailAddress()
}

func (confirmableEmailAddress *confirmableEmailAddress) Equals(other EmailAddress) bool {
	return confirmableEmailAddress.baseEmailAddress.Equals(other)
}

func (confirmableEmailAddress *confirmableEmailAddress) MarshalJSON() ([]byte, error) {
	data := &struct {
		EmailAddress     EmailAddress     `json:"emailAddress"`
		ConfirmationHash ConfirmationHash `json:"confirmationHash"`
	}{
		EmailAddress:     confirmableEmailAddress.baseEmailAddress,
		ConfirmationHash: confirmableEmailAddress.confirmationHash,
	}

	return json.Marshal(data)
}

func UnmarshalConfirmableEmailAddress(data interface{}) (*confirmableEmailAddress, error) {
	var err error

	values, ok := data.(map[string]interface{})
	if !ok {
		return nil, errors.New("zefix")
	}

	confirmableEmailAddress := &confirmableEmailAddress{}

	for key, value := range values {
		value, ok := value.(string)
		if !ok {
			return nil, errors.New("zefix")
		}

		switch key {
		case "emailAddress":
			if confirmableEmailAddress.baseEmailAddress, err = UnmarshalEmailAddress(value); err != nil {
				return nil, err
			}
		case "confirmationHash":
			if confirmableEmailAddress.confirmationHash, err = UnmarshalConfirmationHash(value); err != nil {
				return nil, err
			}
		}
	}

	return confirmableEmailAddress, nil
}

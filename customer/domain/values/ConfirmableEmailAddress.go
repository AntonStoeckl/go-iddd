package values

import (
	"encoding/json"
	"errors"
	"go-iddd/shared"

	"golang.org/x/xerrors"
)

type ConfirmableEmailAddress struct {
	baseEmailAddress *EmailAddress
	confirmationHash *ConfirmationHash
	isConfirmed      bool
}

/*** Factory methods ***/

func buildConfirmableEmailAddress(from *EmailAddress, with *ConfirmationHash) *ConfirmableEmailAddress {
	return &ConfirmableEmailAddress{
		baseEmailAddress: from,
		confirmationHash: with,
		isConfirmed:      false,
	}
}

/*** Getter Methods ***/

func (confirmableEmailAddress *ConfirmableEmailAddress) EmailAddress() string {
	return confirmableEmailAddress.baseEmailAddress.EmailAddress()
}

func (confirmableEmailAddress *ConfirmableEmailAddress) ConfirmationHash() string {
	return confirmableEmailAddress.confirmationHash.Hash()
}

func (confirmableEmailAddress *ConfirmableEmailAddress) IsConfirmed() bool {
	return confirmableEmailAddress.isConfirmed
}

/*** Comparison Methods ***/

func (confirmableEmailAddress *ConfirmableEmailAddress) Equals(other *EmailAddress) bool {
	return confirmableEmailAddress.baseEmailAddress.Equals(other)
}

/*** Modification Methods ***/

func (confirmableEmailAddress *ConfirmableEmailAddress) Confirm(
	given *EmailAddress,
	with *ConfirmationHash,
) (*ConfirmableEmailAddress, error) {

	if !confirmableEmailAddress.baseEmailAddress.Equals(given) {
		return nil, errors.New("confirmableEmailAddress.Confirm: emailAddress is not equal")
	}

	if err := confirmableEmailAddress.confirmationHash.ShouldEqual(with); err != nil {
		return nil, xerrors.Errorf("confirmableEmailAddress.Confirm: %w", err)
	}

	confirmedEmailAddress := buildConfirmableEmailAddress(
		confirmableEmailAddress.baseEmailAddress,
		confirmableEmailAddress.confirmationHash,
	)

	confirmedEmailAddress.isConfirmed = true

	return confirmedEmailAddress, nil
}

/*** Implement json.Marshaler ***/

func (confirmableEmailAddress *ConfirmableEmailAddress) MarshalJSON() ([]byte, error) {
	data := &struct {
		EmailAddress     *EmailAddress     `json:"emailAddress"`
		ConfirmationHash *ConfirmationHash `json:"confirmationHash"`
	}{
		EmailAddress:     confirmableEmailAddress.baseEmailAddress,
		ConfirmationHash: confirmableEmailAddress.confirmationHash,
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return bytes, xerrors.Errorf("confirmableEmailAddress.MarshalJSON: %s: %w", err, shared.ErrMarshalingFailed)
	}

	return bytes, nil
}

/*** Implement json.Unmarshaler ***/

func (confirmableEmailAddress *ConfirmableEmailAddress) UnmarshalJSON(data []byte) error {
	values := &struct {
		EmailAddress     *EmailAddress     `json:"emailAddress"`
		ConfirmationHash *ConfirmationHash `json:"confirmationHash"`
	}{}

	if err := json.Unmarshal(data, values); err != nil {
		return xerrors.Errorf("confirmableEmailAddress.UnmarshalJSON: %s: %w", err, shared.ErrUnmarshalingFailed)
	}

	confirmableEmailAddress.baseEmailAddress = values.EmailAddress
	confirmableEmailAddress.confirmationHash = values.ConfirmationHash

	return nil
}

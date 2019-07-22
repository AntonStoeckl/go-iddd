package values

import (
	"go-iddd/shared"

	jsoniter "github.com/json-iterator/go"
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

func (confirmableEmailAddress *ConfirmableEmailAddress) ShouldConfirm(given *EmailAddress, with *ConfirmationHash) error {
	if err := confirmableEmailAddress.baseEmailAddress.ShouldEqual(given); err != nil {
		return xerrors.Errorf("confirmableEmailAddress.ShouldConfirm: %w", err)
	}

	if err := confirmableEmailAddress.confirmationHash.ShouldEqual(with); err != nil {
		return xerrors.Errorf("confirmableEmailAddress.ShouldConfirm: %w", err)
	}

	confirmedEmailAddress := buildConfirmableEmailAddress(
		confirmableEmailAddress.baseEmailAddress,
		confirmableEmailAddress.confirmationHash,
	)

	confirmedEmailAddress.isConfirmed = true

	return nil
}

/*** Modification Methods ***/

func (confirmableEmailAddress *ConfirmableEmailAddress) MarkAsConfirmed() *ConfirmableEmailAddress {
	confirmedEmailAddress := buildConfirmableEmailAddress(
		confirmableEmailAddress.baseEmailAddress,
		confirmableEmailAddress.confirmationHash,
	)

	confirmedEmailAddress.isConfirmed = true

	return confirmedEmailAddress
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

	bytes, err := jsoniter.Marshal(data)
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

	if err := jsoniter.Unmarshal(data, values); err != nil {
		return xerrors.Errorf("confirmableEmailAddress.UnmarshalJSON: %s: %w", err, shared.ErrUnmarshalingFailed)
	}

	confirmableEmailAddress.baseEmailAddress = values.EmailAddress
	confirmableEmailAddress.confirmationHash = values.ConfirmationHash

	return nil
}

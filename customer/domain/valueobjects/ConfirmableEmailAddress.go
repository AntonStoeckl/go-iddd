package valueobjects

type ConfirmableEmailAddress interface {
    Confirm(given ConfirmationHash) (*confirmableEmailAddress, error)
    IsConfirmed() bool

    EmailAddress
}

type confirmableEmailAddress struct {
    baseEmailAddress *emailAddress
    confirmationHash ConfirmationHash
    isConfirmed      bool
}

func NewConfirmableEmailAddress(from string) (*confirmableEmailAddress, error) {
    baseEmailAddress, err := NewEmailAddress(from)
    if err != nil {
        // TODO: map error?
        return nil, err
    }

    newEmailAddress := newConfirmableEmailAddress(baseEmailAddress, GenerateConfirmationHash(from))

    return newEmailAddress, nil
}

func newConfirmableEmailAddress(from *emailAddress, with ConfirmationHash) *confirmableEmailAddress {
    return &confirmableEmailAddress{
        baseEmailAddress: from,
        confirmationHash: with,
    }
}

func ReconstituteConfirmableEmailAddress(from string, withConfirmationHash string) *confirmableEmailAddress {
    return newConfirmableEmailAddress(
        ReconstituteEmailAddress(from),
        ReconstituteConfirmationHash(withConfirmationHash),
    )
}

func (confirmableEmailAddress *confirmableEmailAddress) Confirm(given ConfirmationHash) (*confirmableEmailAddress, error) {
    if err := confirmableEmailAddress.confirmationHash.MustMatch(given); err != nil {
        return nil, err
    }

    confirmedEmailAddress := newConfirmableEmailAddress(
        confirmableEmailAddress.baseEmailAddress,
        confirmableEmailAddress.confirmationHash,
    )

    confirmedEmailAddress.isConfirmed = true

    return confirmedEmailAddress, nil
}

func (confirmableEmailAddress *confirmableEmailAddress) String() string {
    return confirmableEmailAddress.baseEmailAddress.String()
}

func (confirmableEmailAddress *confirmableEmailAddress) Equals(other EmailAddress) bool {
    return confirmableEmailAddress.baseEmailAddress.Equals(other)
}

func (confirmableEmailAddress *confirmableEmailAddress) IsConfirmed() bool {
    return confirmableEmailAddress.isConfirmed
}

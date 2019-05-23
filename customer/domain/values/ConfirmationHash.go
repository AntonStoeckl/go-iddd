package values

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"go-iddd/shared"
	"math/rand"
	"strconv"
	"time"

	"golang.org/x/xerrors"
)

type ConfirmationHash struct {
	value string
}

/*** Factory methods ***/

func GenerateConfirmationHash(using string) *ConfirmationHash {
	randomInt := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	md5Sum := md5.Sum([]byte(using + strconv.Itoa(randomInt)))
	value := fmt.Sprintf("%x", md5Sum)

	return buildConfirmationHash(value)
}

func RebuildConfirmationHash(from string) (*ConfirmationHash, error) {
	rebuiltConfirmationHash := buildConfirmationHash(from)

	if err := rebuiltConfirmationHash.shouldBeValid(); err != nil {
		return nil, xerrors.Errorf("confirmationHash.New: %s: %w", err, shared.ErrInputIsInvalid)
	}

	return rebuiltConfirmationHash, nil
}

func (confirmationHash *ConfirmationHash) shouldBeValid() error {
	if confirmationHash.value == "" {
		return errors.New("empty input for confirmationHash")
	}

	return nil
}

func buildConfirmationHash(from string) *ConfirmationHash {
	return &ConfirmationHash{value: from}
}

/*** Getter Methods ***/

func (confirmationHash *ConfirmationHash) Hash() string {
	return confirmationHash.value
}

/*** Comparison Methods ***/

func (confirmationHash *ConfirmationHash) ShouldEqual(other *ConfirmationHash) error {
	if confirmationHash.value != other.value {
		return xerrors.Errorf("confirmationHash.ShouldEqual: %w", shared.ErrNotEqual)
	}

	return nil
}

func (confirmationHash *ConfirmationHash) Equals(other *ConfirmationHash) bool {
	return confirmationHash.value == other.value
}

/*** Implement json.Marshaler ***/

func (confirmationHash *ConfirmationHash) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(confirmationHash.value)
	if err != nil {
		return bytes, xerrors.Errorf("confirmationHash.MarshalJSON: %s: %w", err, shared.ErrMarshalingFailed)
	}

	return bytes, nil
}

/*** Implement json.Unmarshaler ***/

func (confirmationHash *ConfirmationHash) UnmarshalJSON(data []byte) error {
	var value string

	if err := json.Unmarshal(data, &value); err != nil {
		return xerrors.Errorf("confirmationHash.UnmarshalJSON: %s: %w", err, shared.ErrUnmarshalingFailed)
	}

	confirmationHash.value = value

	return nil
}

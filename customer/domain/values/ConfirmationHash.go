package values

import (
	"crypto/md5"
	"fmt"
	"go-iddd/shared"
	"math/rand"
	"strconv"
	"time"

	"github.com/cockroachdb/errors"
	jsoniter "github.com/json-iterator/go"
)

type ConfirmationHash struct {
	value string
}

/*** Factory methods ***/

func GenerateConfirmationHash(using string) *ConfirmationHash {
	randomInt := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	md5Sum := md5.Sum([]byte(using + strconv.Itoa(randomInt)))
	value := fmt.Sprintf("%x", md5Sum)

	return &ConfirmationHash{value: value}
}

func ConfirmationHashFrom(input string) (*ConfirmationHash, error) {
	rebuiltConfirmationHash := &ConfirmationHash{value: input}

	if err := rebuiltConfirmationHash.shouldBeValid(); err != nil {
		return nil, errors.Wrap(errors.Mark(err, shared.ErrInputIsInvalid), "confirmationHash.ConfirmationHashFrom")
	}

	return rebuiltConfirmationHash, nil
}

func (confirmationHash *ConfirmationHash) shouldBeValid() error {
	if confirmationHash.value == "" {
		return errors.New("empty input for confirmationHash")
	}

	return nil
}

/*** Getter Methods ***/

func (confirmationHash *ConfirmationHash) Hash() string {
	return confirmationHash.value
}

/*** Comparison Methods ***/

func (confirmationHash *ConfirmationHash) Equals(other *ConfirmationHash) bool {
	return confirmationHash.value == other.value
}

/*** Implement json.Marshaler ***/

func (confirmationHash *ConfirmationHash) MarshalJSON() ([]byte, error) {
	bytes, err := jsoniter.Marshal(confirmationHash.value)
	if err != nil {
		return nil, errors.Wrap(errors.Mark(err, shared.ErrMarshalingFailed), "confirmationHash.MarshalJSON")
	}

	return bytes, nil
}

/*** Implement json.Unmarshaler ***/

func (confirmationHash *ConfirmationHash) UnmarshalJSON(data []byte) error {
	var value string

	if err := jsoniter.Unmarshal(data, &value); err != nil {
		return errors.Wrap(errors.Mark(err, shared.ErrUnmarshalingFailed), "confirmationHash.UnmarshalJSON")
	}

	confirmationHash.value = value

	return nil
}

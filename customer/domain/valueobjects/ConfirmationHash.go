package valueobjects

import (
	"crypto/md5"
	"encoding/json"
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
	md5Sum := md5.Sum([]byte(strconv.Itoa(randomInt) + using))
	value := fmt.Sprintf("%x", md5Sum)

	return buildConfirmationHash(value)
}

func ReconstituteConfirmationHash(from string) *ConfirmationHash {
	return buildConfirmationHash(from)
}

func buildConfirmationHash(from string) *ConfirmationHash {
	return &ConfirmationHash{value: from}
}

/*** Public methods implementing ConfirmationHash ***/

func (confirmationHash *ConfirmationHash) Hash() string {
	return confirmationHash.value
}

func (confirmationHash *ConfirmationHash) MustMatch(other *ConfirmationHash) error {
	if confirmationHash.Hash() != other.Hash() {
		return xerrors.Errorf("confirmationHash.MustMatch: input does not match: %w", shared.ErrInvalidInput) // TODO: use a distinct error type?
	}

	return nil
}

func (confirmationHash *ConfirmationHash) MarshalJSON() ([]byte, error) {
	return json.Marshal(confirmationHash.value)
}

func UnmarshalConfirmationHash(data interface{}) (*ConfirmationHash, error) {
	value, ok := data.(string)
	if !ok {
		return nil, xerrors.Errorf("UnmarshalConfirmationHash: input is not [string]: %w", shared.ErrUnmarshaling)
	}

	return &ConfirmationHash{value: value}, nil
}

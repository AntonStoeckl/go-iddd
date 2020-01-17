package values

import (
	"crypto/md5"
	"fmt"
	"go-iddd/shared"
	"math/rand"
	"strconv"
	"time"

	"github.com/cockroachdb/errors"
)

type ConfirmationHash struct {
	value string
}

func GenerateConfirmationHash(using string) ConfirmationHash {
	randomInt := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	md5Sum := md5.Sum([]byte(using + strconv.Itoa(randomInt)))
	value := fmt.Sprintf("%x", md5Sum)

	return ConfirmationHash{value: value}
}

func BuildConfirmationHash(input string) (ConfirmationHash, error) {
	if input == "" {
		err := shared.MarkAndWrapError(
			errors.New("empty input for confirmationHash"),
			shared.ErrInputIsInvalid,
			"BuildConfirmationHash",
		)

		return ConfirmationHash{}, err
	}

	confirmationHash := ConfirmationHash{value: input}

	return confirmationHash, nil
}

func RebuildConfirmationHash(input string) ConfirmationHash {
	return ConfirmationHash{value: input}
}

func (confirmationHash ConfirmationHash) Hash() string {
	return confirmationHash.value
}

func (confirmationHash ConfirmationHash) Equals(other ConfirmationHash) bool {
	return confirmationHash.value == other.value
}

package values

import (
	"crypto/sha256"
	"fmt"
	"go-iddd/service/lib"
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
	sha256Sum := sha256.Sum256([]byte(using + strconv.Itoa(randomInt)))
	value := fmt.Sprintf("%x", sha256Sum)

	return ConfirmationHash{value: value}
}

func BuildConfirmationHash(input string) (ConfirmationHash, error) {
	if input == "" {
		err := lib.MarkAndWrapError(
			errors.New("empty input for confirmationHash"),
			lib.ErrInputIsInvalid,
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

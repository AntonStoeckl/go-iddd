package value

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/cockroachdb/errors"
)

type ConfirmationHash string

func GenerateConfirmationHash(using string) ConfirmationHash {
	//nolint:gosec // no super secure random number needed here
	randomInt := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	sha256Sum := sha256.Sum256([]byte(using + strconv.Itoa(randomInt)))
	value := fmt.Sprintf("%x", sha256Sum)

	return ConfirmationHash(value)
}

func BuildConfirmationHash(input string) (ConfirmationHash, error) {
	if input == "" {
		err := errors.New("empty input for confirmationHash")
		err = shared.MarkAndWrapError(err, shared.ErrInputIsInvalid, "BuildConfirmationHash")

		return "", err
	}

	confirmationHash := ConfirmationHash(input)

	return confirmationHash, nil
}

func RebuildConfirmationHash(input string) ConfirmationHash {
	return ConfirmationHash(input)
}

func (confirmationHash ConfirmationHash) String() string {
	return string(confirmationHash)
}

func (confirmationHash ConfirmationHash) Equals(other ConfirmationHash) bool {
	return confirmationHash.String() == other.String()
}

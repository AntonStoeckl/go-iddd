package valueobjects

import (
	"crypto/md5"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type ConfirmationHash interface {
	Hash() string
	MustMatch(other ConfirmationHash) error
}

type confirmationHash struct {
	value string
}

/*** Factory methods ***/

func GenerateConfirmationHash(using string) *confirmationHash {
	randomInt := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	md5Sum := md5.Sum([]byte(strconv.Itoa(randomInt) + using))
	value := fmt.Sprintf("%x", md5Sum)

	return buildConfirmationHash(value)
}

func ReconstituteConfirmationHash(from string) *confirmationHash {
	return buildConfirmationHash(from)
}

func buildConfirmationHash(from string) *confirmationHash {
	return &confirmationHash{value: from}
}

/*** Public methods implementing ConfirmationHash ***/

func (confirmationHash *confirmationHash) Hash() string {
	return confirmationHash.value
}

func (confirmationHash *confirmationHash) MustMatch(other ConfirmationHash) error {
	if confirmationHash.Hash() != other.Hash() {
		return errors.New("confirmationHash - is not equal")
	}

	return nil
}

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
	String() string
	MustMatch(other ConfirmationHash) error
}

type confirmationHash struct {
	value string
}

func GenerateConfirmationHash(using string) *confirmationHash {
	randomInt := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	md5Sum := md5.Sum([]byte(strconv.Itoa(randomInt) + using))
	value := fmt.Sprintf("%x", md5Sum)
	newConfirmationHash := ReconstituteConfirmationHash(value)

	return newConfirmationHash
}

func ReconstituteConfirmationHash(from string) *confirmationHash {
	return &confirmationHash{value: from}
}

func (confirmationHash *confirmationHash) String() string {
	return confirmationHash.value
}

func (confirmationHash *confirmationHash) MustMatch(other ConfirmationHash) error {
	if confirmationHash.String() != other.String() {
		return errors.New("confirmationHash - is not equal")
	}

	return nil
}

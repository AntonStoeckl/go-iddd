package value

import (
	"strings"

	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/cockroachdb/errors"
)

const minPasswordLength = 12
const maxPasswordLength = 250

type PlainPassword string

func BuildPlainPassword(input string) (PlainPassword, error) {
	input = strings.TrimSpace(input)

	if err := validate(input); err != nil {
		return "", shared.MarkAndWrapError(err, shared.ErrInputIsInvalid, "BuildPlainPassword")
	}

	return PlainPassword(input), nil
}

func RebuildPlainPassword(input string) (PlainPassword, error) {
	return PlainPassword(input), nil
}

func validate(input string) error {
	if input == "" {
		err := errors.New("empty input for PlainPassword")

		return err
	}

	if len(input) < minPasswordLength {
		err := errors.Newf(
			"input for PlainPassword is too short, min [%d] characters required, but [%d] supplied",
			minPasswordLength,
			len(input),
		)

		return err
	}

	if len(input) > maxPasswordLength {
		err := errors.Newf(
			"input for PlainPassword is too long, max [%d] characters allowed, but [%d] supplied",
			maxPasswordLength,
			len(input),
		)

		return err
	}

	return nil
}

func (pw PlainPassword) String() string {
	return string(pw)
}

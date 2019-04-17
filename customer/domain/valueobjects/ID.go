package valueobjects

import "github.com/google/uuid"

type ID interface {
	String() string
}

type id struct {
	value string
}

func GenerateID() *id {
	uid, err := uuid.NewRandom()
	if err != nil {
		panic("id - could not generate uid: " + err.Error())
	}

	return ReconstituteID(uid.String())
}

func ReconstituteID(from string) *id {
	return &id{value: from}
}

func (id *id) String() string {
	return id.value
}

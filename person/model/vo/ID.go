package vo

type ID interface {
	ID() string
}

type id struct {
	value string
}

func NewID(value string) *id {
	newID := &id{value: value}

	return newID
}

func (id *id) ID() string {
	return id.value
}

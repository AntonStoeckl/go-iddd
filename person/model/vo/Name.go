package vo

type Name interface {
	GivenName() string
	FamilyName() string
}

type name struct {
	givenName  string
	familyName string
}

func NewName(givenName string, familyName string) *name {
	newName := &name{
		givenName:  givenName,
		familyName: familyName,
	}

	return newName
}

func (name *name) GivenName() string {
	return name.givenName
}

func (name *name) FamilyName() string {
	return name.familyName
}

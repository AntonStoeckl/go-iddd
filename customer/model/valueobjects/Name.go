package valueobjects

type Name interface {
	GivenName() string
	FamilyName() string
}

type name struct {
	givenName  string
	familyName string
}

func NewName(givenName string, familyName string) *name {
	newName := ReconstituteName(givenName, familyName)
	// TODO: validate

	return newName
}

func ReconstituteName(givenName string, familyName string) *name {
	reconstitutedName := &name{
		givenName:  givenName,
		familyName: familyName,
	}

	return reconstitutedName
}

func (name *name) GivenName() string {
	return name.givenName
}

func (name *name) FamilyName() string {
	return name.familyName
}

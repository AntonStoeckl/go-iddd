package valueobjects

type PersonName interface {
    GivenName() string
    FamilyName() string
}

type personName struct {
    givenName  string
    familyName string
}

func NewPersonName(givenName string, familyName string) *personName {
    newPersonName := ReconstitutePersonName(givenName, familyName)
    // TODO: validate

    return newPersonName
}

func ReconstitutePersonName(givenName string, familyName string) *personName {
    reconstitutedPersonName := &personName{
        givenName:  givenName,
        familyName: familyName,
    }

    return reconstitutedPersonName
}

func (personName *personName) GivenName() string {
    return personName.givenName
}

func (personName *personName) FamilyName() string {
    return personName.familyName
}

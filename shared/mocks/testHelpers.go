package mocks

type SomethingHappend struct{}

func (event *SomethingHappend) Identifier() string {
	return "something"
}

func (event *SomethingHappend) EventName() string {
	return "something"
}

func (event *SomethingHappend) OccurredAt() string {
	return "something"
}

func (event *SomethingHappend) StreamVersion() uint {
	return 1
}

type SomethingElseHappend struct{}

func (event *SomethingElseHappend) Identifier() string {
	return "something"
}

func (event *SomethingElseHappend) EventName() string {
	return "something"
}

func (event *SomethingElseHappend) OccurredAt() string {
	return "something"
}

func (event *SomethingElseHappend) StreamVersion() uint {
	return 1
}

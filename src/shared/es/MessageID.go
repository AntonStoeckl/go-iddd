package es

import (
	"github.com/google/uuid"
)

type MessageID string

func GenerateMessageID() MessageID {
	return MessageID(uuid.New().String())
}

func BuildMessageID(from MessageID) MessageID {
	return from
}

func RebuildMessageID(from string) MessageID {
	return MessageID(from)
}

func (messageID MessageID) String() string {
	return string(messageID)
}

func (messageID MessageID) Equals(other MessageID) bool {
	return messageID.String() == other.String()
}

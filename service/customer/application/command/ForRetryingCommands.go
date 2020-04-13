package command

type ForRetryingCommands func(originalFunc func() error, maxRetries uint8) error

package application

type ForRetryingCommands func(originalFunc func() error, maxRetries uint8) error

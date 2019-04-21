package shared

type CommandHandler interface {
    Handle(command Command) error
}

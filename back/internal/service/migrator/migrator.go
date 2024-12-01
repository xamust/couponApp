package migrator

type Migrator interface {
	Commands(command string, args ...string) error
}

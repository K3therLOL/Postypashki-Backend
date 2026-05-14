package domain

type BrockerProvider interface {
	Send(message string) error
}

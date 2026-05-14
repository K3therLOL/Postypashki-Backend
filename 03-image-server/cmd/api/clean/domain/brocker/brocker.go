package domain

type BrockerProvider interface {
	Send(message string) error
	Recv() (string, error)
}

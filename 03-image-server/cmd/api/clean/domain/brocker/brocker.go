package domain

type BrockerProvider interface {
	Send(body string) error
}

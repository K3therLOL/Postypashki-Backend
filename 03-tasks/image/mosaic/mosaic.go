package image

import "log"

type Mosaic struct {
	logger *log.Logger
}

func NewMosaic() *Mosaic {
	return &Mosaic{
		logger: log.Default(),
	}
}

func (m *Mosaic) Build() error {
	return nil
}

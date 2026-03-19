package usecase

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"task/clean/domain"
)

var (
	ErrNotImage = errors.New("Not a image.")
)

type ImageUsecase struct {
	image domain.Image
}

func isImage(body io.ReadCloser) (io.Reader, bool) {
	const mimeLength = 512
	buff := make([]byte, mimeLength)

	n, _ := io.ReadFull(body, buff)
	contentType := http.DetectContentType(buff[:n])
	if !strings.HasPrefix(contentType, "image/") {
		return nil, false
	}

	originalBody := io.MultiReader(bytes.NewReader(buff[:n]), body)
	return originalBody, true
}

func (iu *ImageUsecase) Process(url string) (*domain.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, ok := isImage(resp.Body)
	if !ok {
		return nil, ErrNotImage
	}

	/* TO DO */
}

package parser

import (
	"errors"
	"io"
	"mime/multipart"

	"github.com/bytedance/sonic"
)

func ParseMultipartFileToBytes(file *multipart.FileHeader) (*[]byte, *[]byte, error) {
	if file == nil {
		err := errors.New("file cannot be nil")
		logger.Error(err)
		return nil, nil, err
	}

	f, err := file.Open()
	if err != nil {
		logger.Error(err)
		return nil, nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			logger.Error(err)
		}
	}()

	fileBytes, err := io.ReadAll(f)
	if err != nil {
		logger.Error(err)
		return nil, nil, err
	}

	headerBytes, err := sonic.Marshal(file.Header)
	if err != nil {
		logger.Error(err)
		return nil, nil, err
	}

	return &fileBytes, &headerBytes, nil
}

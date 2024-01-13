package utils

import (
	"io"
	"mime/multipart"
	"os"
)

func WriteBytes(values map[string]io.Reader, w *multipart.Writer) error {
	for key, r := range values {
		var fw io.Writer
		var err error
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return err
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return err
			}
		}
		_, err = io.Copy(fw, r)
		if err != nil {
			return err
		}
	}
	w.Close()
	return nil
}

// Author: Steve Zhang
// Date: 2020/10/17 9:46 下午

package file

import (
	"io"
	"mime/multipart"
	"os"
)

func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

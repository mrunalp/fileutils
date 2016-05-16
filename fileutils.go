package fileutils

import (
	"io"
	"os"
)

// CopyFile copies the file at source to dest
func CopyFile(source string, dest string) error {
	sf, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sf.Close()

	df, err := os.Open(dest)
	if err != nil {
		return err
	}
	defer df.Close()

	_, err = io.Copy(sf, df)
	if err != nil {
		return err
	}

	return nil
}

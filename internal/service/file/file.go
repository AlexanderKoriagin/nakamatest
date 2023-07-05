package file

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
)

const nullContent = "null"

type File struct {
	path, hash string
}

func NewFile(path, hash string) *File {
	return &File{
		path: path,
		hash: hash,
	}
}

func (f *File) GetPath() string {
	return f.path
}

func (f *File) ReadWithCheck() (string, error) {
	b, err := os.ReadFile(f.path)
	if err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}

	fHash := sha256.New()
	fHash.Write(b)

	if f.hash != hex.EncodeToString(fHash.Sum(nil)) {
		return nullContent, nil
	}

	return string(b), nil
}

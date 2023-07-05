package file_test

import (
	"crypto/sha256"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/akrillis/nakamatest/internal/service/file"
)

func TestFile_GetPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "ok",
			path: "path",
			want: "path",
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			f := file.NewFile(tt.path, "")
			require.Equal(t, tt.want, f.GetPath())
		})
	}
}

func TestFile_ReadWithCheck(t *testing.T) {
	content := `{"field0": "value0", "field1": 1}`
	hashIncorrect := "null"

	hash := sha256.New()
	hash.Write([]byte(content))
	hashCorrect := string(hash.Sum(nil))

	incorrectType := "incorrectType1234567890"
	correctType := "correctType1234567890"
	correctVersion := "correctVersion1234567890"

	require.NoError(t, os.MkdirAll(correctType, 0755))
	defer func() {
		_ = os.RemoveAll(correctType)
	}()

	f, err := os.OpenFile(correctType+"/"+correctVersion+".json", os.O_CREATE|os.O_WRONLY, 0644)
	require.NoError(t, err)
	_, err = f.WriteString(content)
	require.NoError(t, err)
	require.NoError(t, f.Close())

	tests := []struct {
		name    string
		path    string
		hash    string
		want    string
		wantErr bool
	}{
		{
			name:    "file doesn't exist",
			path:    incorrectType + "/" + correctVersion + ".json",
			hash:    hashCorrect,
			want:    "",
			wantErr: true,
		},
		{
			name:    "wrong hash",
			path:    correctType + "/" + correctVersion + ".json",
			hash:    hashIncorrect,
			want:    "null",
			wantErr: false,
		},
		{
			name:    "correct file",
			path:    correctType + "/" + correctVersion + ".json",
			hash:    hashCorrect,
			want:    content,
			wantErr: false,
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {

			fi := file.NewFile(test.path, test.hash)
			res, err := fi.ReadWithCheck()

			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.want, res)
			}
		})
	}
}

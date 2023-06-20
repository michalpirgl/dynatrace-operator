package support_archive

import (
	"bytes"
	"github.com/klauspost/compress/zip"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestZipAddFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	tarFile, err := createTarFile(tmpDir)
	require.NoError(t, err)
	tarball := newZipball(tarFile)

	fileName := tarFile.Name()

	testString := []byte(`Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.`)
	require.NoError(t, tarball.addFile("lorem-ipsum.txt", bytes.NewReader(testString)))
	tarball.close()
	tarFile.Close()

	resultFile, err := os.OpenFile(fileName, os.O_RDONLY, os.ModeTemporary)
	require.NoError(t, err)
	defer tarFile.Close()
	stat, err := resultFile.Stat()
	assert.NoError(t, err)
	zipReader, err := zip.NewReader(resultFile, stat.Size())
	require.NoError(t, err)

	zippedFile := zipReader.File[0]
	require.NoError(t, err)
	assert.Equal(t, "lorem-ipsum.txt", zippedFile.Name)
	r, err := zippedFile.Open()
	assert.NoError(t, err)

	resultString := make([]byte, 1024)
	resultLen, err := r.Read(resultString)
	require.Equal(t, io.EOF, err)
	assert.Equal(t, len(testString), resultLen)
	assert.Equal(t, testString, resultString[:resultLen])

	resultFile.Close()
}

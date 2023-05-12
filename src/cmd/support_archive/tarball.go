package support_archive

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const tarFileName = "%s/operator-support-archive-%s.tgz"

type tarball struct {
	tarWriter   *tar.Writer
	gzipWriter  *gzip.Writer
	memoryLimit int64
}

func newTarball(target io.Writer) tarball {
	newTarball := tarball{
		gzipWriter:  gzip.NewWriter(target),
		memoryLimit: defaultMemoryLimit / 10,
	}
	newTarball.tarWriter = tar.NewWriter(newTarball.gzipWriter)
	return newTarball
}

func (t *tarball) setMemoryLimit(limit int64) {
	t.memoryLimit = limit
}

func (t tarball) close() {
	if t.tarWriter != nil {
		t.tarWriter.Close()
	}
	if t.gzipWriter != nil {
		t.gzipWriter.Close()
	}
}

func (t tarball) addFileTmpIntermediate(fileName string, file io.Reader) error {
	tmpFile, err := os.CreateTemp(defaultSupportArchiveTargetDir, "support-archive-tmp")
	if err != nil {
		return errors.WithMessagef(err, "could not create temp file for %s", fileName)
	}
	defer func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}()

	_, err = io.Copy(tmpFile, file)
	if err != nil {
		return errors.WithMessagef(err, "could not copy data from source for '%s'", fileName)
	}

	_, err = tmpFile.Seek(0, 0)
	if err != nil {
		return errors.WithMessagef(err, "could not go back to start of tmp file for '%s'", fileName)
	}

	fileInfo, err := tmpFile.Stat()
	if err != nil {
		return errors.WithMessagef(err, "could get file size for '%s'", fileName)
	}

	header := &tar.Header{
		Name: fileName,
		Size: fileInfo.Size(),
		Mode: int64(fs.ModePerm),
	}

	err = t.tarWriter.WriteHeader(header)
	if err != nil {
		return errors.WithMessagef(err, "could not write header for file '%s'", fileName)
	}

	_, err = io.Copy(t.tarWriter, tmpFile)
	if err != nil {
		return errors.WithMessagef(err, "could not copy the file '%s' data to the tarball", fileName)
	}
	return nil
}

func (t tarball) addFileInParts(fileName string, file io.Reader) error {
	fileCount := 0
	done := false
	adaptedFileName := fileName
	for !done {
		buffer := &bytes.Buffer{}
		copied, err := io.CopyN(buffer, file, t.memoryLimit)

		switch {
		case errors.Is(err, io.EOF):
			done = true
		case err != nil:
			return errors.WithMessagef(err, "could not copy data from source for '%s'", fileName)
		}

		if copied == 0 && fileCount > 0 {
			// we want to make sure to not suppress 0 byte files at all, but don't want to create
			// a last 0-byte part
			return nil
		}

		// only rename if a file is not done in one rush
		if !(fileCount == 0 && done) {
			adaptedFileName = fmt.Sprintf("%s.%d", fileName, fileCount)
		}

		header := &tar.Header{
			Name: adaptedFileName,
			Size: int64(buffer.Len()),
			Mode: int64(fs.ModePerm),
		}

		err = t.tarWriter.WriteHeader(header)
		if err != nil {
			return errors.WithMessagef(err, "could not write header for file '%s'", fileName)
		}

		_, err = io.Copy(t.tarWriter, buffer)
		if err != nil {
			return errors.WithMessagef(err, "could not copy the file '%s' data to the tarball", fileName)
		}
		fileCount++
	}
	return nil
}

func createTarballTargetFile(useStdout bool, targetDir string) (*os.File, error) {
	if useStdout {
		return os.Stdout, nil
	} else {
		tarFile, err := createTarFile(targetDir)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return tarFile, nil
	}
}

func createTarFile(targetDir string) (*os.File, error) {
	tarballFilePath := fmt.Sprintf(tarFileName, targetDir, time.Now().Format(time.RFC3339))
	tarballFilePath = strings.ReplaceAll(tarballFilePath, ":", "_")

	tarFile, err := os.Create(tarballFilePath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return tarFile, nil
}

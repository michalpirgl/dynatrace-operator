package support_archive

import (
	"fmt"
	"github.com/klauspost/compress/zip"
	"github.com/pkg/errors"
	"io"
	"os"
	"strings"
	"time"
)

type tarballer interface {
	close()
	addFile(fileName string, file io.Reader) error
}

const zipFileName = "%s/operator-support-archive-%s.zip"

type zipball struct {
	writer *zip.Writer
}

func newZipball(target io.Writer) zipball {
	newZipball := zipball{writer: zip.NewWriter(target)}

	return newZipball
}

func (z zipball) close() {
	if z.writer != nil {
		z.writer.Close()
	}
}

func (z zipball) addFile(fileName string, reader io.Reader) error {
	w, err := z.writer.Create(fileName)
	if err != nil {
		return errors.WithMessagef(err, "could not write header for file '%s'", fileName)
	}

	_, err = io.Copy(w, reader)
	if err != nil {
		return errors.WithMessagef(err, "could not copy the file '%s' data to the tarball", fileName)
	}
	return nil
}

func createZipballTargetFile(useStdout bool, targetDir string) (*os.File, error) {
	if useStdout {
		return os.Stdout, nil
	} else {
		tarFile, err := createZipFile(targetDir)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return tarFile, nil
	}
}

func createZipFile(targetDir string) (*os.File, error) {
	zipballFilePath := fmt.Sprintf(zipFileName, targetDir, time.Now().Format(time.RFC3339))
	zipballFilePath = strings.ReplaceAll(zipballFilePath, ":", "_")

	tarFile, err := os.Create(zipballFilePath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return tarFile, nil
}

package zipexp

import (
	"archive/zip"
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestMultiFileStdZip(t *testing.T) {
	testLogFileCompressionStdZipInMem(logFileSizeInMB, numberOfLogFiles, t)
}

func testLogFileCompressionStdZipInMem(testSizeInMB int64, numberOfFiles int, t *testing.T) {

	testSize := 1024 * 1024 * testSizeInMB
	writer := countingDevNullWriter{}

	zw := zip.NewWriter(&writer)
	defer zw.Close()

	for i := 0; i < numberOfFiles; i++ {
		generatingReader := loadGenerator{bytesLeftToRead: testSize}
		fn := fmt.Sprintf("file%d.txt", i)

		t.Logf("Writing %s", fn)

		w, err := zw.Create(fn)
		assert.NoError(t, err)

		_, err = io.Copy(w, &generatingReader)
		assert.NoError(t, err)
	}
	LogWrittenBytes(writer.count, t)
}

func testLogFileCompressionStdZipToFile(testSizeInMB int64, numberOfFiles int, t *testing.T) {

	testSize := 1024 * 1024 * testSizeInMB
	file, err := os.Create("logs.zip")
	assert.NoError(t, err)
	defer file.Close()
	writer := bufio.NewWriter(file)

	zw := zip.NewWriter(writer)
	defer zw.Close()
	for i := 0; i < numberOfFiles; i++ {
		generatingReader := loadGenerator{bytesLeftToRead: testSize}
		fn := fmt.Sprintf("file%d.txt", i)

		t.Logf("Writing %s", fn)

		w, err := zw.Create(fn)
		assert.NoError(t, err)
		_, err = io.Copy(w, &generatingReader)
		assert.NoError(t, err)
	}

	assert.NoError(t, err)
	err = writer.Flush()
	assert.NoError(t, err)
	stat, err := file.Stat()
	assert.NoError(t, err)
	LogWrittenBytes(stat.Size(), t)
}

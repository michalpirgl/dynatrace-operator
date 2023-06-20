package zipexp

import (
	"bufio"
	"fmt"
	"github.com/klauspost/compress/zip"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestMultiFileKZip(t *testing.T) { testFileCompressionKZipInMem(100, 5, t) }

func testFileCompressionKZipInMem(testSizeInMB int64, numberOfFiles int, t *testing.T) {

	testSize := 1024 * 1024 * testSizeInMB
	writer := countingDevNullWriter{}

	zw := zip.NewWriter(&writer)

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
	zw.Close()
}

func testFileCompressionKZipToFile(testSizeInMB int64, numberOfFiles int, t *testing.T) {

	testSize := 1024 * 1024 * testSizeInMB
	file, err := os.Create("logs.zip")
	assert.NoError(t, err)
	defer file.Close()
	writer := bufio.NewWriter(file)

	zw := zip.NewWriter(writer)

	for i := 0; i < numberOfFiles; i++ {
		generatingReader := loadGenerator{bytesLeftToRead: testSize}
		fn := fmt.Sprintf("file%d.txt", i)

		t.Logf("Writing %s", fn)

		w, err := zw.Create(fn)
		assert.NoError(t, err)

		_, err = io.Copy(w, &generatingReader)
		assert.NoError(t, err)
	}
	zw.Close()
	assert.NoError(t, err)
	err = writer.Flush()
	assert.NoError(t, err)
	stat, err := file.Stat()
	assert.NoError(t, err)
	LogWrittenBytes(stat.Size(), t)
}

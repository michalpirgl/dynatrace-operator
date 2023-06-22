package zipexp

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMultiFileTar(t *testing.T) { testLogFileCompressionTar(logFileSizeInMB, numberOfLogFiles, t) }

func testLogFileCompressionTar(testSizeInMB int64, numberOfFiles int, t *testing.T) {
	testSize := 1024 * 1024 * testSizeInMB
	writer := countingDevNullWriter{}
	tarball := newTarball(&writer)
	defer tarball.close()

	for i := 0; i < numberOfFiles; i++ {
		generatingReader := loadGenerator{bytesLeftToRead: testSize}
		fn := fmt.Sprintf("file%d.txt", i)

		t.Logf("Writing %s", fn)
		err := tarball.addFile(fn, &generatingReader)
		assert.NoError(t, err)
	}

	LogWrittenBytes(writer.count, t)
}

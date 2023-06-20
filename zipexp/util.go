package zipexp

import (
	"github.com/bcicen/go-units"
	"golang.org/x/exp/rand"
	"io"
	"testing"
)

type loadGenerator struct {
	bytesLeftToRead int64
}

// var alphaNumeric = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 ")
var alphaNumeric = []byte("a")

func (lg *loadGenerator) Read(p []byte) (n int, err error) {
	if lg.bytesLeftToRead <= 0 {
		return 0, io.EOF
	}

	i := 0

	for ; i < len(p) && lg.bytesLeftToRead > 0; i++ {
		p[i] = alphaNumeric[rand.Intn(len(alphaNumeric))]
		lg.bytesLeftToRead--
	}
	return i, nil
}

type countingDevNullWriter struct {
	count int64
}

func (d *countingDevNullWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	//fmt.Printf("Write called %d\n", n)
	d.count = d.count + int64(n)
	return n, nil
}

func LogWrittenBytes(bytesWritten int64, t *testing.T) {
	val := units.NewValue(float64(bytesWritten), units.Byte)
	opts := units.FmtOptions{
		Label:     true, // append unit name/symbol
		Short:     true, // use unit symbol
		Precision: 3,
	}
	t.Logf("Bytes written: %s", val.MustConvert(units.MegaByte).Fmt(opts))
}

package support_archive

import (
	"golang.org/x/exp/rand"
	"io"
)

type randomAlphanumericCharGenerator struct {
	bytesLeftToRead int64
}

var alphaNumeric = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 ")

//var alphaNumeric = []byte("a")

func (lg *randomAlphanumericCharGenerator) Read(p []byte) (n int, err error) {
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

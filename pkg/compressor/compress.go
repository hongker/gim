package compressor

import (
	"errors"
	"io"
)

const (
	ENCODING_GZIP    = "gzip"
	ENCODING_DEFLATE = "deflate"
)

type CompressingResponseWriter struct {
	compressor io.WriteCloser
	encoding   string
}

func (w CompressingResponseWriter) Write(p []byte) (int, error) {
	return w.compressor.Write(p)
}

func NewCompressingResponseWriter(writer io.Writer, encoding string) (*CompressingResponseWriter, error) {
	c := &CompressingResponseWriter{}

	c.encoding = encoding
	if ENCODING_GZIP == encoding {
		w := currentCompressorProvider.AcquireGzipWriter()
		w.Reset(writer)
		c.compressor = w
	} else if ENCODING_DEFLATE == encoding {
		w := currentCompressorProvider.AcquireZlibWriter()
		w.Reset(writer)
		c.compressor = w
	} else {
		return nil, errors.New("Unknown encoding:" + encoding)
	}

	return c, nil
}

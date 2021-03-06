package sivafs

import (
	"io"
	"os"

	"gopkg.in/src-d/go-billy.v3"
	"gopkg.in/src-d/go-siva.v1"
)

type file struct {
	name        string
	closeNotify func() error
	isClosed    bool

	w siva.Writer
	r *io.SectionReader
}

func newFile(filename string, w siva.Writer, closeNotify func() error) billy.File {
	return &file{
		name:        filename,
		closeNotify: closeNotify,
		w:           w,
	}
}

func openFile(filename string, r *io.SectionReader) billy.File {
	return &file{
		name: filename,
		r:    r,
	}
}

func (f *file) Name() string {
	return f.name
}

func (f *file) Read(p []byte) (int, error) {
	if f.isClosed {
		return 0, os.ErrClosed
	}

	if f.r == nil {
		return 0, ErrWriteOnlyFile
	}

	return f.r.Read(p)
}

func (f *file) ReadAt(b []byte, off int64) (int, error) {
	if f.isClosed {
		return 0, os.ErrClosed
	}

	if f.r == nil {
		return 0, ErrWriteOnlyFile
	}

	return f.r.ReadAt(b, off)
}

func (f *file) Seek(offset int64, whence int) (int64, error) {
	if f.isClosed {
		return 0, os.ErrClosed
	}

	if f.r == nil {
		return 0, ErrNonSeekableFile
	}

	return f.r.Seek(offset, whence)
}

func (f *file) Write(p []byte) (int, error) {
	if f.isClosed {
		return 0, os.ErrClosed
	}

	if f.w == nil {
		return 0, ErrReadOnlyFile
	}

	return f.w.Write(p)
}

func (f *file) Close() error {
	if f.isClosed {
		return os.ErrClosed
	}

	defer func() { f.isClosed = true }()

	if f.closeNotify == nil {
		return nil
	}

	return f.closeNotify()
}

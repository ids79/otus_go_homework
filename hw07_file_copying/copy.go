package main

import (
	"errors"
	"io"
	"os"
	"time"

	"github.com/cheggaaa/pb"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fSrc, err := os.Open(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}
	defer fSrc.Close()
	fDst, err := os.Create(toPath)
	if err != nil {
		return ErrUnsupportedFile
	}
	info, err := fSrc.Stat()
	if err != nil {
		return ErrUnsupportedFile
	}
	size := info.Size()
	if size == 0 {
		return ErrUnsupportedFile
	}
	if size < offset {
		return ErrOffsetExceedsFileSize
	}
	if offset > 0 {
		fSrc.Seek(offset, 0)
	}
	count := 100
	bar := pb.StartNew(count)
	needToCopy := size
	if limit > 0 {
		needToCopy = limit
	}
	copied := int64(0)
	for i := 1; i <= count; i++ {
		part := needToCopy / int64(count)
		if i == count {
			part = needToCopy - copied
		}
		copied += part
		io.CopyN(fDst, fSrc, part)
		bar.Increment()
		time.Sleep(time.Millisecond * 1)
	}
	bar.Finish()
	err = fDst.Close()
	if err != nil {
		return ErrUnsupportedFile
	}
	return nil
}

package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile          = errors.New("unsupported file")
	ErrOffsetExceedsFileSize    = errors.New("offset exceeds file size")
	ErrSameSourceAndDestination = errors.New("same source and destination")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if fromPath == toPath {
		return ErrSameSourceAndDestination
	}
	fSrc, err := os.Open(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}
	defer fSrc.Close()
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
	fDst, err := os.Create(toPath)
	if err != nil {
		return ErrUnsupportedFile
	}
	if offset > 0 {
		_, err := fSrc.Seek(offset, io.SeekStart)
		if err != nil {
			return ErrUnsupportedFile
		}
	}
	if limit == 0 {
		limit = size
	}
	if offset+limit > size {
		limit = size - offset
	}
	bar := pb.Full.Start64(limit)
	barReader := bar.NewProxyReader(fSrc)
	copied, err := io.CopyN(fDst, barReader, limit)
	if err != nil || copied < limit {
		return err
	}
	bar.Finish()
	err = fDst.Close()
	if err != nil {
		return ErrUnsupportedFile
	}
	return nil
}

package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath string, toPath string, offset, limit int64) error {
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("error open fromPath file: %w", err)
	}
	defer fromFile.Close()

	stat, err := fromFile.Stat()
	if err != nil {
		return fmt.Errorf("error getting fromPath file stat: %w", err)
	}

	fromFileSize := stat.Size()
	if offset > fromFileSize {
		return ErrOffsetExceedsFileSize
	}

	_, err = fromFile.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("seek error: %w", err)
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("error creating toPath file: %w", err)
	}
	defer toFile.Close()

	if limit == 0 || offset+limit > fromFileSize {
		limit = fromFileSize - offset
	}

	bar := pb.Start64(limit)
	barReader := bar.NewProxyReader(fromFile)
	_, err = io.CopyN(toFile, barReader, limit)
	bar.Finish()

	if errors.Is(err, io.EOF) || err == nil {
		return nil
	}

	return fmt.Errorf("copy error: %w", err)
}

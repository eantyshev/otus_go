package main

import (
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"io"
	"os"
)

func getReader(from string, offset, limit int64) (r *io.LimitedReader, err error) {
	var fInfo os.FileInfo
	var f *os.File
	var maxbytes int64
	if f, err = os.Open(from); err != nil {
		return
	}
	if fInfo, err = f.Stat(); err != nil {
		return
	}
	// 0 size for character files is a special case
	maxbytes = fInfo.Size()
	if offset > 0 {
		if offset >= maxbytes {
			return nil, fmt.Errorf("offset %d should be less than the file size (%d)", offset, maxbytes)
		}
		if _, err = f.Seek(offset, io.SeekStart); err != nil {
			return
		}
		maxbytes -= offset
	}
	if limit > 0 {
		if maxbytes > limit {
			maxbytes = limit
		}
	}
	r = io.LimitReader(f, maxbytes).(*io.LimitedReader)
	return
}

func copyData(from, to string, offset, limit int64) error {
	var r *io.LimitedReader
	var w io.Writer
	var err error
	if r, err = getReader(from, offset, limit); err != nil {
		return err
	}
	if w, err = os.Create(to); err != nil {
		return err
	}
	bar := pb.Full.Start64(r.N)
	defer bar.Finish()
	barReader := bar.NewProxyReader(r)
	if _, err = io.Copy(w, barReader); err != nil {
		return err
	}
	return err
}

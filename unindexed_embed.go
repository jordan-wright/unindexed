package unindexed

import (
	"embed"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"
)

// EmbedFS is an implementation of a standard http.FileSystem
// without the ability to list files in the directory.
// This implementation is largely inspired by
// https://www.alexedwards.net/blog/disable-http-fileserver-directory-listings
type EmbedFS struct {
	fs embed.FS
}

func (ef EmbedFS) Open(name string) (http.File, error) {
	if name == "/" {
		name = "."
	} else {
		name = strings.TrimPrefix(name, "/")
	}
	f, err := ef.fs.Open(name)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
		index := strings.TrimPrefix(name, "/") + "/index.html"
		_, err := ef.fs.Open(index)
		if err != nil {
			return nil, os.ErrPermission
		}
	}
	return ioFile{f}, nil
}

// ioFile is copied from net/http Package
type ioFile struct {
	file fs.File
}

func (f ioFile) Close() error               { return f.file.Close() }
func (f ioFile) Read(b []byte) (int, error) { return f.file.Read(b) }
func (f ioFile) Stat() (fs.FileInfo, error) { return f.file.Stat() }

var errMissingSeek = errors.New("io.File missing Seek method")
var errMissingReadDir = errors.New("io.File directory missing ReadDir method")

func (f ioFile) Seek(offset int64, whence int) (int64, error) {
	s, ok := f.file.(io.Seeker)
	if !ok {
		return 0, errMissingSeek
	}
	return s.Seek(offset, whence)
}

func (f ioFile) ReadDir(count int) ([]fs.DirEntry, error) {
	d, ok := f.file.(fs.ReadDirFile)
	if !ok {
		return nil, errMissingReadDir
	}
	return d.ReadDir(count)
}

func (f ioFile) Readdir(count int) ([]fs.FileInfo, error) {
	d, ok := f.file.(fs.ReadDirFile)
	if !ok {
		return nil, errMissingReadDir
	}
	var list []fs.FileInfo
	for {
		dirs, err := d.ReadDir(count - len(list))
		for _, dir := range dirs {
			info, err := dir.Info()
			if err != nil {
				// Pretend it doesn't exist, like (*os.File).Readdir does.
				continue
			}
			list = append(list, info)
		}
		if err != nil {
			return list, err
		}
		if count < 0 || len(list) >= count {
			break
		}
	}
	return list, nil
}

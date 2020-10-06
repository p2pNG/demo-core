package file_store

import (
	"crypto/sha512"
	"errors"
	"io"
	"os"
	"time"
)

type FileInfo struct {
	Name string
	Size int64
	Hash []byte

	BlockSize int64
	BlockHash [][]byte
}

type LocalFileInfo struct {
	FileInfo
	Path       string
	LastModify time.Time
}

const DefaultHashBufferSize = 4 * 1024 * 1024

func StatLocalFile(filepath string, blockSize int64) (lf *LocalFileInfo, err error) {
	stat, err := os.Stat(filepath)
	if err != nil {
		return
	}
	if stat.IsDir() {
		err = errors.New("not a valid file")
		return
	}
	fi, err := StatFile(filepath, blockSize)
	if err != nil {
		return
	}

	lf = new(LocalFileInfo)
	lf.FileInfo = *fi
	lf.LastModify = stat.ModTime()
	lf.Path = filepath
	return
}
func StatFile(filepath string, blockSize int64) (fi *FileInfo, err error) {
	if blockSize <= 1024*1024 {
		blockSize = DefaultHashBufferSize
	}

	stat, err := os.Stat(filepath)
	if err != nil {
		return
	}
	if stat.IsDir() {
		err = errors.New("not a valid file")
		return
	}
	fi = new(FileInfo)
	fi.BlockSize = blockSize
	fi.Name, fi.Size = stat.Name(), stat.Size()

	f, err := os.Open(filepath)
	if err != nil {
		return
	}
	buf := make([]byte, blockSize)
	fileSum := sha512.New()
	blockHash := sha512.New512_256()
	flagTail := false
	var n int

	for {
		n, err = f.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return
			}
		}
		if int64(n) != blockSize {
			flagTail = true
			break
		}

		fileSum.Write(buf)
		blockHash.Reset()
		blockHash.Write(buf)
		fi.BlockHash = append(fi.BlockHash, blockHash.Sum(nil))
	}
	if flagTail {
		if int64(len(fi.BlockHash))*blockSize+int64(n) != fi.Size {
			err = errors.New("read file error, length not matched")
			return
		}
		fileSum.Write(buf)
		blockHash.Reset()
		blockHash.Write(buf)
		fi.BlockHash = append(fi.BlockHash, blockHash.Sum(nil))
	}
	fi.Hash = fileSum.Sum(nil)
	return fi, nil
}

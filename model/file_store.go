package model

import "time"

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

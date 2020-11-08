package Utils

import (
	"NPProj3/ORM"
	"NPProj3/Wigets"
	"errors"
	"os"
	"sync"
)

type FileItem struct {
	Uuid           string
	FileDescriptor *os.File
	FileInfo       ORM.FileInfo
	Offset         int64
	State          int
	Path           string
}

type FileList struct {
	Length int
	List   []FileItem
	Lock   sync.Mutex
	Err    FileErr
}

type fileState struct {
	Finish int
	Wrong  int
	Wait   int
}

var FileState = fileState{1, 2, 3}

type FileErr struct {
	errDescription string
}

func (f *FileErr) Error() error {
	return errors.New(f.errDescription)
}

func (f *FileItem) WriteIn(data []byte) error {
	cnt, err := f.FileDescriptor.WriteAt(data, f.Offset)
	f.Offset += int64(cnt)
	return err
}

func (f *FileList) AddFile(uuid string, fileDescriptor *os.File, fileInfo ORM.FileInfo) error {
	fileItem := FileItem{
		Uuid:           uuid,
		FileInfo:       fileInfo,
		FileDescriptor: fileDescriptor,
		Offset:         int64(0),
		State:          FileState.Wait,
	}
	f.Lock.Lock()
	defer f.Lock.Unlock()
	if f.isExist(fileInfo.MD5) {
		f.Err.errDescription = "duplicate file"
		return f.Err.Error()
	} else {
		f.List = append(f.List, fileItem)
		f.Length += 1
		return nil
	}
}

//this func haven't lock for duplicated lock operation
// if you want a lock , call IsExist()
func (f *FileList) isExist(md5 string) bool {
	for _, file := range f.List {
		if file.FileInfo.MD5 == md5 {
			return true
		}
	}
	return false
}

func (f *FileList) IsExist(md5 string) bool {
	f.Lock.Lock()
	defer f.Lock.Unlock()
	for _, file := range f.List {
		if file.FileInfo.MD5 == md5 {
			return false
		}
	}
	return true
}

func (f *FileList) FindFileItemByUUID(uuid string) *FileItem {
	f.Lock.Lock()
	defer f.Lock.Unlock()
	for _, file := range f.List {
		if file.Uuid == uuid {
			return &file
		}
	}
	return nil
}

func (f *FileList) IsUUIDExist(uuid string) bool {
	for _, file := range f.List {
		if file.Uuid == uuid {
			return true
		}
	}
	return false
}

func (f *FileList) WriteBytes(uuid string, data []byte) bool {
	file := f.FindFileItemByUUID(uuid)
	if file != nil {
		err := file.WriteIn(data)
		Wigets.ErrHandle(err)
		return true
	} else {
		return false
	}
}

// You should call IsUUIDExist() First Before Use this func
func (f *FileList) Finish(uuid string) error {
	for _, file := range f.List {
		if file.Uuid == uuid {
			md5Value := FileMD5FileDescriptor(file.FileDescriptor)
			if md5Value == file.FileInfo.MD5 {
				file.State = FileState.Finish
				err := file.FileDescriptor.Close()
				Wigets.ErrHandle(err)
				return nil
			} else {
				file.State = FileState.Wrong
				return errors.New("md5 checksum failed")
			}

		}
	}
	return errors.New("no such file")
}

func (f *FileList) FindFileItemByMD5(md5 string) *FileItem {
	for _, file := range f.List {
		if file.FileInfo.MD5 == md5 {
			return &file
		}
	}
	return nil
}

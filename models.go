package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var fileEntries = FileEntryDb{}

type (
	FileEntryDb map[int64]FileEntry
)

// FileEntry represents a single uploaded file
type FileEntry struct {
	Id         int64
	Name       string
	Size       int64 // in bytes
	UploadDate time.Time
}

func GetEntryById(id int64) FileEntry {
	return fileEntries[id]
}

func GetFileEntries() []FileEntry {
	var ret []FileEntry
	for _, v := range fileEntries {
		ret = append(ret, v)
	}
	return ret
}

func getNextId() int64 {
	i := int64(0)
	for j := range fileEntries {
		i = max(i, j)
	}
	return i + 1
}

func HasFileEntry(id int64) bool {
	_, ok := fileEntries[id]
	return ok
}

func CreateFileEntry(filename string, filesize int64) FileEntry {
	id := getNextId()
	entry := FileEntry{
		getNextId(),
		filename,
		filesize,
		time.Now(),
	}

	fileEntries[id] = entry
	return entry
}

func RemoveFileEntry(id int64) {
	delete(fileEntries, id)
}

func (f *FileEntry) GetFilepath() string {
	return filepath.Join(Config.FilesPath, i64ToStr(f.Id))
}

type View struct {
	viewid int
	files  []FileEntry
	token  string
}

func loadFileEntries(path string) FileEntryDb {
	file := panicOnErr(os.OpenFile(path, os.O_CREATE|os.O_RDONLY, 0644))
	defer closeOrPanic(file)

	var js []FileEntry
	content := panicOnErr(ioutil.ReadAll(file))
	checkErr(json.Unmarshal(content, &js))

	ret := FileEntryDb{}
	for _, v := range js {
		ret[v.Id] = v
	}
	return ret
}

func writeFileEntries(path string) {
	fmt.Println("Writing: ", fileEntries)
	file := panicOnErr(os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644))
	defer closeOrPanic(file)

	out := panicOnErr(json.Marshal(GetFileEntries()))
	_ = panicOnErr(file.Write(out))
}

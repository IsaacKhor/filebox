package main

import (
	"encoding/csv"
	"github.com/labstack/gommon/log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var fileEntries = FileEntryDb{}

type FileEntryDb map[int]FileEntry

// FileEntry represents a single uploaded file
type FileEntry struct {
	Id         int
	Name       string
	Size       int // in bytes
	UploadDate time.Time
}

func init() {
	fileEntries = loadFileEntries(DbPath)
}

func GetEntryById(id int) FileEntry {
	return fileEntries[id]
}

func GetFileEntries() []FileEntry {
	var ret []FileEntry
	for _, v := range fileEntries {
		ret = append(ret, v)
	}
	return ret
}

func (f *FileEntry) toCsvRow() []string {
	return []string{
		strconv.Itoa(f.Id),
		f.Name,
		strconv.Itoa(f.Size),
		strconv.Itoa(int(f.UploadDate.Unix())),
	}
}

func (f *FileEntry) GetFilepath() string {
	return filepath.Join(FilesPath, strconv.Itoa(f.Id))
}

func fromCsvRow(row []string) FileEntry {
	if len(row) != 4 {
		log.Fatal("Invalid row length")
	}
	id := dieOnErr2(strconv.Atoi(row[0]))
	size := dieOnErr2(strconv.Atoi(row[2]))
	unixts := dieOnErr2(strconv.Atoi(row[3]))

	return FileEntry{
		id,
		row[1],
		size,
		time.Unix(int64(unixts), 0),
	}
}

type View struct {
	viewid int
	files  []FileEntry
	token  string
}

func loadFileEntries(path string) FileEntryDb {
	in := dieOnErr2(os.OpenFile(path, os.O_CREATE|os.O_RDONLY, 0644))
	defer in.Close()
	reader := csv.NewReader(in)
	rows := dieOnErr2(reader.ReadAll())

	var entriesDb = map[int]FileEntry{}
	for _, v := range rows {
		r := fromCsvRow(v)
		entriesDb[r.Id] = r
	}

	return entriesDb
}

func writeFileEntries(path string) {
	log.Debugf("writing: %s", fileEntries)
	o := dieOnErr2(os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644))
	defer o.Close()

	writer := csv.NewWriter(o)
	var result [][]string
	for _, fileEntry := range fileEntries {
		result = append(result, fileEntry.toCsvRow())
	}

	checkErr(writer.WriteAll(result))
}

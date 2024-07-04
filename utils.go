package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"

	"github.com/labstack/gommon/log"
	"golang.org/x/exp/constraints"
)

type FileboxConfig struct {
	DbPath    string
	FilesPath string
	AuthToken string
	Port      int
	Host      string
}

const (
	ConfigPath = "./config.json"
)

var (
	Config = FileboxConfig{
		DbPath:    "./files/filesdb.json",
		FilesPath: "./files",
		AuthToken: "DEFAULT_TOKEN",
		Port:      7001,
		Host:      "filebox.isaackhor.com",
	}
)

func loadConfig() {
	file := panicOnErr(os.OpenFile(ConfigPath, os.O_CREATE|os.O_RDWR,
		0644))
	defer closeOrPanic(file)

	content := panicOnErr(ioutil.ReadAll(file))
	checkErr(json.Unmarshal(content, &Config))
}

// ToBinarySuffix assumes that we don't have sizes > 1TiB
func ToBinarySuffix(filesize int64) string {
	suffixes := []string{"", "KiB", "MiB", "GiB", "TiB"}
	newsize := float64(filesize)

	i := 0
	for newsize > 1024 {
		newsize = newsize / 1024
		i++
	}

	num := strconv.FormatFloat(newsize, 'f', 3, 64)
	return num + suffixes[i]
}

func checkErr(e error) {
	if e != nil {
		log.Error(e)
		panic(e)
	}
}

func panicOnErr[T any](x T, e error) T {
	if e != nil {
		log.Error(e)
		panic(e)
	}
	return x
}

func closeOrPanic(x io.Closer) {
	checkErr(x.Close())
}

func max[T constraints.Ordered](a T, b T) T {
	if a > b {
		return a
	}
	return b
}

func i64ToStr(i int64) string {
	return strconv.FormatInt(i, 10)
}

func strToI64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func isRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("[isRoot] Unable to get current user: %s", err)
	}
	return currentUser.Username == "root"
}

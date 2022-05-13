package main

import (
	"github.com/labstack/gommon/log"
	"strconv"
)

// ToBinarySuffix assumes that we don't have sizes > 1TiB
func ToBinarySuffix(filesize int) string {
	suffixes := []string{"", "KiB", "MiB", "GiB", "TiB"}
	newsize := float64(filesize)

	i := 0
	for newsize < 1024 {
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

func dieOnErr2[T any](x T, e error) T {
	checkErr(e)
	return x
}

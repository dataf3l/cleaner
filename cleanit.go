package main

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"syscall"
	"time"
)

func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}

func checkFile(filePath string, info os.FileInfo, days int) {

	// Sys () returns interface {}, so you need a type assertion. Different platforms need different types. On linux, * syscall. Stat_t
	stat := info.Sys().(*syscall.Stat_t)
	if stat == nil {
		return // whoops
	}

	// atime, CTime and mtime are access time, creation time and modification time, respectively. See Man 2 stat for details.
	//fmt.Println(timespecToTime(stat.Atimespec))
	created := timespecToTime(stat.Mtim) // change to Mtimespec for osx

	now := time.Now()
	cutoffTime := now.AddDate(0, 0, -days)
	//log.Println(cutoffTime)
	if created.Before(cutoffTime) {
		log.Println("DELETE:\t" + path.Base(filePath))
		os.Remove(filePath)
	} else {
		//log.Println("MAYBE WAIT:\t" + path.Base(filePath))
	}
}
func main() {
	if len(os.Args) < 4 {
		log.Println("Usage:    cleaner a-folder       'some-regexp' DAYS")
		log.Println(" ")
		log.Println(" Example: cleaner ./your-files/  '/deleteme/' 365")
		log.Println(" Example: cleaner ./backups/     '/.*.zip$/'    30")
		log.Println(" Example: cleaner ./dump-files/  '/.*.sql$/'     1")
		log.Println(" ")
		log.Println(" Notice the .* instead of * and also notice the use of single quotes around")
		log.Println(" The Regular Expression, this is to avoid BASH SHELL expansion ")
		log.Println(" ")
		return
	}
	location := os.Args[1]
	regex := os.Args[2]
	daysAsString := os.Args[3]
	re, err := regexp.Compile(regex)
	if err != nil {
		log.Println("Bad Regular Expression, see golang's reference. ")
		log.Println(err)
		return
	}
	days, err := strconv.Atoi(daysAsString)
	if err != nil {
		log.Println("Days should be int")
		log.Println(err)
	}

	//log.Println(re)
	//return

	err = filepath.Walk(location, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			//log.Println("DIRECTORY IGNORING:\t" + path)
			return nil // continue
		}
		base := path.Base(filePath)

		if !re.MatchString(base) {
			//log.Println("MAYBE IGNR:\t" + base)
			return nil // continue
		}
		checkFile(filePath, info, days)
		return nil
	})
	if err != nil {
		panic(err)
	}

}

// https://developpaper.com/getting-access-creation-modification-time-of-files-on-linux-using-golang/
// https://superuser.com/questions/146125/how-to-preserve-file-attributes-when-one-copies-files-in-windows

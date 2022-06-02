package dateutil


import (
    "io/fs"
    "path/filepath"
    "fmt"
)

func GetPathByDate(fileinfo fs.FileInfo) (out string, err error) {
	/* given info about a file, return the directory where the file
	   should reside */

    t := getDate(fileinfo)
	yearstr := fmt.Sprintf("%v", t.Year())
	monthname := swmonths[t.Month()-1]
	out = filepath.Join(".", yearstr, monthname, fileinfo.Name())
	out = filepath.ToSlash(out)
	return
}

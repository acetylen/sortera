//+build linux

package dateutil


import(
    "io/fs"
    "time"
)

func getDate(fileinfo fs.FileInfo) time.Time {
	return fileinfo.ModTime()
}

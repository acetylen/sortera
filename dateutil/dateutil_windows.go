//+build windows

package dateutil

import (
    "io/fs"
    "syscall"
    "time"
)

func getDate(fileinfo fs.FileInfo) time.Time {

	sysinfo := fileinfo.Sys()

	win := sysinfo.(*syscall.Win32FileAttributeData)
	return time.Unix(0, win.CreationTime.Nanoseconds())
}


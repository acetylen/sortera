package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

var this os.FileInfo
var dry_run bool
var swmonths = [12]string{
	"Januari",
	"Februari",
	"Mars",
	"April",
	"Maj",
	"Juni",
	"Juli",
	"Augusti",
	"September",
	"Oktober",
	"November",
	"December",
}

func init() {
	var err error
	this, err = os.Stat(os.Args[0])
	if err != nil {
		log.Fatalln("Failed to stat self:", err)
	}

	debug := flag.Bool("debug", false, "Print debugging information")
	flag.BoolVar(&dry_run, "dry-run", false, "Only list what would happen, don't actually do it")
	flag.Parse()

	if *debug {
		log.SetOutput(os.Stderr)
	} else {
		devnull, err := os.Open(os.DevNull)
		if err != nil {
			log.Fatalln("Failed to open null file")
		}
		log.SetOutput(devnull)
	}

}

func getAllFilesToMove(f fs.FS, out map[string]string) error {
	/* Given a file tree, and a mapping, populate with files that are in the
	   wrong location, to where they should be. */
	return fs.WalkDir(f, ".", func(path string, info fs.DirEntry, err error) error {
		log.Println(path)
		if err != nil {
			log.Println("Got an error:", err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		fileinfo, err := info.Info()
		if err != nil {
			log.Println("Failed to get file info:", err)
			return err
		}
		if os.SameFile(fileinfo, this) {
			return nil
		}
		destpath, err := getPathByDate(fileinfo)
		if err != nil {
			log.Println("Failed to get date path:", err)
			return err
		}

		if destpath != path {
			out[path] = destpath
		}
		return nil

	})
}

func getAllEmptyDirs(f fs.FS) ([]string, error) {
	/* Given a file tree, return a list of directories which only contain
	   directories or nothing */
	dirs := make(map[string]int, 65535)

	// count files in each directory
	err := fs.WalkDir(f, ".", func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == "." {
			return nil
		}

		normalized := filepath.FromSlash(path)
		if info.IsDir() {
			dirs[normalized] = 0
		}

		dirs[filepath.Dir(normalized)]++
		return nil
	})

	if err != nil {
		return nil, err
	}

	// mark directories that are empty (or only contain directories) for removal
	var empty []string
	for {
		var remove []string
		for dirname, numfiles := range dirs {
			if numfiles == 0 {
				remove = append(remove, dirname)
				dirs[filepath.Dir(dirname)] -= 1
			}
		}

		if len(remove) == 0 {
			break
		}

		for _, dir := range remove {
			delete(dirs, dir)
			empty = append(empty, dir)
		}
	}

	return empty, nil
}

func getPathByDate(fileinfo fs.FileInfo) (out string, err error) {
	/* given info about a file, return the directory where the file
	   should reside */

	/*
	   if windows
	   sysinfo := fileinfo.Sys()

	   win := sysinfo.(*syscall.Win32FileAttributeData)
	   t := time.Unix(win.CreationTime.Nanoseconds())
	*/
	t := fileinfo.ModTime()

	yearstr := fmt.Sprintf("%v", t.Year())
	monthname := swmonths[t.Month()-1]
	out = filepath.Join(".", yearstr, monthname, fileinfo.Name())
	out = filepath.ToSlash(out)
	return
}

func getFirstFreeNameFor(filename string, blocklist map[string]string) (string, error) {
	// Given a filename, make sure it's unique by
	// adding a number if the name is already taken.
	_, err := os.Stat(filename)
	fileExists := !errors.Is(err, fs.ErrNotExist)
	if !fileExists {
		return filename, nil
	}

    basename := filepath.Base(filename)
    suffix := filepath.Ext(filename)

	for i := 1; i < 64; i++ {
		newname := fmt.Sprintf("%v_%v%v", basename, i, suffix)
		_, err = os.Stat(newname)
		fileExists = !errors.Is(err, fs.ErrNotExist)
		if !fileExists {
			accepted := true
			for _, f := range blocklist {
				if newname == f {
					accepted = false
					break
				}
			}
			if accepted {
				return newname, nil
			}
		}
	}
	return "", errors.New("Failed to make a new filename")
}

func main() {
	root, err := os.Getwd()
	if err != nil {
		log.Fatalln("Unable to get source directory:", err)
	}
	log.Println("I'm", this.Name())
	log.Println("Starting in", root)

	fs := os.DirFS(root)

	badfiles := make(map[string]string, 65535)

	log.Println("Finding misplaced files...")
	err = getAllFilesToMove(fs, badfiles)
	if err != nil {
		log.Fatalln("Failed to find bad files:", err)
	}

	log.Printf("Found %v misplaced files\n", len(badfiles))
	for src, dst := range badfiles {
		parent := filepath.Dir(dst)
		if _, err := os.Stat(parent); errors.Is(err, os.ErrNotExist) {

			if !dry_run {
				err = os.MkdirAll(parent, os.ModePerm)
				if err != nil {
					log.Fatalln("Failed to create directory:", err)
				}
			}
		}
		uniqdst, err := getFirstFreeNameFor(dst, badfiles)
		if err != nil {
			log.Fatalln("Failed to create new filename:", err)
		}
		log.Printf("%v -> %v\n", src, uniqdst)
		if !dry_run {
			err = os.Rename(src, uniqdst)
			if err != nil {
				log.Fatalln("Failed to move file:", err)
			}
		}
	}

	log.Println("Finding empty directories...")
	emptydirs, err := getAllEmptyDirs(fs)
	if err != nil {
		log.Fatalln("Failed to find empty dirs:", err)
	}
	log.Printf("Found %v empty directories\n", len(emptydirs))
	for _, d := range emptydirs {
		log.Println("rm", d)
		if !dry_run {
			err = os.Remove(d)
			if err != nil {
				log.Fatalln("Failed to remove empty dir:", err)
			}
		}
	}
}

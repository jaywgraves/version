package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage:  %s -c|-m [-d -t] [-u=string] files... \n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nThis utility will copy all the files that are given as arguments and give each\n")
	fmt.Fprintf(os.Stderr, "new file a name that has a datetime string (either current or file modification\n")
	fmt.Fprintf(os.Stderr, "time) preceding the file extension.\n\n")
	fmt.Fprintf(os.Stderr, "If neither the -m or the -c flag is set, then -m will be defaulted.\n")
	fmt.Fprintf(os.Stderr, "Any combination of the -d and -t flags can be set. If neither are set,\n")
	fmt.Fprintf(os.Stderr, "then both will be defaulted.\n")
	fmt.Fprintf(os.Stderr, "The -u flag will use the given string in place of a date/time.\n\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "  files...: 1 or more file name specifications\n")
}

func CopyFile(dst, src string) (int64, error) {
	sf, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer sf.Close()
	df, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer df.Close()
	return io.Copy(df, sf)
}

func CreateFileVersion(dst, src string, fi os.FileInfo) {
	CopyFile(dst, src)
	os.Chtimes(dst, time.Now(), fi.ModTime())
}

func newName(oldname, versionstring string) (oldnameabs, newname string) {
	oldnameabs, _ = filepath.Abs(oldname)
	dir, fname := filepath.Split(oldname)
	fext := filepath.Ext(fname)
	fbase := fname[:len(fname)-len(fext)]
	var newfname string
	if len(fbase) == 0 {
		// hidden files get the version string appended
		newfname = fext + "." + versionstring
	} else {
		// regular files keep their extension and embed the version string
		newfname = fbase + "." + versionstring + fext
	}
	newname = filepath.Join(dir, newfname)
	return

}
func main() {
	modflg := flag.Bool("m", false, "modification: use file modification time in version string")
	currflg := flag.Bool("c", false, "current: use current date in version string")
	dateflg := flag.Bool("d", false, "date: add YYMMDD to version string")
	timeflg := flag.Bool("t", false, "time: add HHMMSS to version string")
	customflg := flag.String("u", "", "custom: provide a short string for naming")
	silentflg := flag.Bool("s", false, "silent: suppress output")

	flag.Usage = Usage
	flag.Parse()
	switch {
	case *modflg && *currflg:
		fmt.Fprintf(os.Stderr, "Can not set both the -m and -c flags.\n")
		Usage()
		return
	case !*modflg && !*currflg:
		*modflg = true
	}
	fmtstring := ""
	switch {
	case (*dateflg && *timeflg) || (!*dateflg && !*timeflg):
		fmtstring = "20060102.150405"
	case *dateflg && !*timeflg:
		fmtstring = "20060102"
	case !*dateflg && *timeflg:
		fmtstring = "150405"
	}
	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "No file specifications given as arguments.\n")
		Usage()
		return
	}
	nowstr := ""
	if *currflg {
		// all files should have same time stamp
		nowstr = time.Now().Format(fmtstring)
	}
	var versionstring string
	exitcode := 0
	for _, a := range flag.Args() {
		matches, _ := filepath.Glob(a)
		if len(matches) == 0 {
			if !*silentflg {
				fmt.Printf("no matches found for '%s'.\n", a)
			}
			exitcode = 42
		}
		for _, m := range matches {
			fi, _ := os.Stat(m)
			if fi.IsDir() && !*silentflg {
				fmt.Printf("cannot version '%s': it is a directory\n", m)
				continue
			}
			switch {
			case *modflg:
				versionstring = fi.ModTime().Format(fmtstring)
			case *currflg:
				versionstring = nowstr
			}
			if len(*customflg) > 0 {
				versionstring = *customflg
			}
			oldname, newname := newName(m, versionstring)
			if !*silentflg {
				fmt.Println(m, "->", newname)
			}
			CreateFileVersion(newname, oldname, fi)
		}
	}
	os.Exit(exitcode)
}

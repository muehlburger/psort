// Copy and rename file based on exif information.
//
// Herbert Muehlburger <mail@muehlburger.at>
package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
)

var verbose = flag.Bool("v", false, "show verbose progress messages")
var dst = flag.String("d", "sorted", "define destiation folder for sorted files")

var supportedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

func main() {
	// Determine the initial directories.
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	// Traverse the file tree.
	paths := make(chan string)
	go func() {
		for _, root := range roots {
			err := walkDir(root, paths)
			if err != nil {
				log.Printf("psort %s: %v\n", root, err)
			}
		}
		close(paths)
	}()

	// Print the results periodically.
	var tick <-chan time.Time
	if *verbose {
		tick = time.Tick(500 * time.Millisecond)
	}
	var nfiles int64
	var names []string
loop:
	for {
		select {
		case name, ok := <-paths:
			if !ok {
				break loop // paths was closed
			}
			nfiles++
			if _, ok := supportedExtensions[strings.ToLower(filepath.Ext(name))]; ok {
				names = append(names, name)
			}
		case <-tick:
			fmt.Printf("%d files\n", nfiles)
		}
	}

	if nfiles > 0 {
		if err := os.MkdirAll(*dst, 0755); err != nil {
			log.Fatal(err)
		}
	}

	for _, f := range names {
		if err := CopyFile(f, *dst); err != nil {
			log.Fatal(err)
		}
	}

}

// CopyFile copies files from src to destination.
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	dstFilename, err := dstFilename(src)
	if err != nil {
		return err
	}
	outfile := filepath.Join(dst, dstFilename)
	log.Printf("%s -> %s", src, outfile)
	out, err := os.OpenFile(outfile, os.O_RDWR|os.O_CREATE, 0744)
	if err != nil {
		return err
	}

	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()
	_, err = io.Copy(out, in)

	fi, err := os.Stat(in.Name())
	if err != nil {
		return err
	}

	if err = os.Chmod(out.Name(), fi.Mode()); err != nil {
		os.Remove(out.Name())
		return err
	}

	// TODO set corret mtime and ctime
	if err = os.Chtimes(out.Name(), fi.ModTime(), fi.ModTime()); err != nil {
		os.Remove(out.Name())
		return err
	}
	return nil
}

func dstFilename(file string) (string, error) {
	tm, err := getCreationDatetime(file)
	if err != nil {
		return "", err
	}
	ext := strings.ToLower(filepath.Ext(file)) // e.g., ".jpg", ".JPEG"
	outfile := fmt.Sprintf("%d-%02d-%02d_%02d%02d%02d%s", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second(), ext)
	outfile = filepath.Base(outfile)
	return outfile, nil
}

// CreationTime extracts the creation Date of the given File.
func getCreationDatetime(path string) (time.Time, error) {
	f, err := os.Open(path)
	if err != nil {
		return time.Time{}, err
	}
	defer f.Close()

	exif.RegisterParsers(mknote.All...)

	x, err := exif.Decode(f)
	if err != nil {
		return time.Time{}, err
	}

	tm, err := x.DateTime()
	if err != nil {
		return time.Time{}, err
	}

	return tm, nil
}

// walkDir recursively walks the file tree rooted at dir
// and sends the absolute path of each found file on paths.
func walkDir(dir string, paths chan<- string) error {
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("could not read dir %s: %v", dir, err)
	}
	for _, fi := range fis {
		if !fi.IsDir() {
			path, err := filepath.Abs(filepath.Join(dir, fi.Name()))
			if err != nil {
				log.Fatal(err)
			}
			paths <- path
		} else {
			walkDir(filepath.Join(dir, fi.Name()), paths)
		}
	}
	return nil
}

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
			walkDir(root, paths)
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
			names = append(names, name)
		case <-tick:
			fmt.Printf("%d files found.\n", nfiles)
		}
	}
	printPaths(nfiles, names)
	for _, f := range names {
		if err := CopyFile(f); err != nil {
			log.Fatal(err)
		}
	}

}

//func main() {
//	dst := "export"
//
//	if len(os.Args) > 1 {
//		src = os.Args[1]
//	}
//
//	src, err := filepath.Abs(src)
//	log.Printf("sourcefile: %v", src)
//	if err != nil {
//		log.Fatalf("absolute %s: %v", src, err)
//	}
//
//	filename := filepath.Base(src)
//	log.Printf("filename: %s", filename)
//
//	log.Printf("destination: %v", dst)
//	os.MkdirAll(dst, 0744)
//
//	dst, err = filepath.Abs(filepath.Join(dst, filename))
//	if err != nil {
//		log.Fatalf("absolute %s: %v", dst, err)
//	}
//	log.Printf("destination: %v", dst)
//
//	if err = CopyFile(src, dst); err != nil {
//		log.Fatal(err)
//	}
//}

// CopyFile copies files from src to destination.
func CopyFile(path string) error {
	//TODO copy only images
	in, err := os.Open(path)
	if err != nil {
		return err
	}
	defer in.Close()

	dst := createFile(path)

	log.Printf("creating file: %s", dst)
	out, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0644)
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
	return nil
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

func createFile(infile string) string {
	ext := strings.ToLower(filepath.Ext(infile)) // e.g., ".jpg", ".JPEG"
	outfile := strings.TrimSuffix(infile, ext) + ".renamed" + ext
	fmt.Printf("infile: %v\n", infile)
	fmt.Printf("ext: %v\n", ext)
	fmt.Printf("outfile: %v\n", outfile)
	return outfile
}

func printPaths(nfiles int64, paths []string) {
	fmt.Printf("%d files found:\n", nfiles)
	for i, p := range paths {
		fmt.Printf("%d\t%v\n", i, p)
	}
}

// walkDir recursively walks the file tree rooted at dir
// and sends the absolute path of each found file on paths.
func walkDir(dir string, paths chan<- string) {
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			walkDir(subdir, paths)
		} else {
			p, err := filepath.Abs(filepath.Join(dir, entry.Name()))
			if err != nil {
				log.Fatal(err)
			}
			paths <- p
		}
	}
}

// dirents returns the entries of directory dir.
func dirents(dir string) []os.FileInfo {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return nil
	}
	return entries
}

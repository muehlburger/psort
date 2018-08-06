package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
)

// psort -s source -t
func main() {
	src, _ := os.Getwd()
	dst := "./export"

	if len(os.Args) > 1 {
		src = os.Args[1]
	}

	src, err := filepath.Abs(src)
	log.Printf("source: %v", src)
	if err != nil {
		log.Fatalf("absolute %s: %v", src, err)
	}

	filename := filepath.Base(src)
	log.Printf("filename: %s", filename)

	dst, err = filepath.Abs(filepath.Join(dst, filename))
	if err != nil {
		log.Fatalf("absolute %s: %v", dst, err)
	}
	log.Printf("destination: %v", dst)

	if err = CopyFile(src, dst, filename); err != nil {
		log.Fatal(err)
	}
}

// CopyFile copies files from src to destination.
func CopyFile(src, dst, filename string) error {
	in, err := os.Open(src)

	if err != nil {
		return err
	}
	defer in.Close()

	tm, err := getCreationDatetime(src)
	if err != nil {
		return err
	}

	log.Printf("creating file: %s", dst)
	tmp, err := ioutil.TempFile(filepath.Dir(dst), "")
	if err != nil {
		return err
	}
	_, err = io.Copy(tmp, in)
	if err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return err
	}

	if err = tmp.Close(); err != nil {
		return err
	}

	fi, err := os.Stat(in.Name())
	if err != nil {
		return err
	}

	if err = os.Chmod(tmp.Name(), fi.Mode()); err != nil {
		os.Remove(tmp.Name())
		return err
	}

	if err = os.Rename(tmp.Name(), dst); err != nil {
		os.Remove(tmp.Name())
		return err
	}

	fmt.Printf("%s", tm)

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

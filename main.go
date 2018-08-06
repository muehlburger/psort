package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
)

func main() {
	src, _ := os.Getwd()
	dst := "export"

	if len(os.Args) > 1 {
		src = os.Args[1]
	}

	src, err := filepath.Abs(src)
	log.Printf("sourcefile: %v", src)
	if err != nil {
		log.Fatalf("absolute %s: %v", src, err)
	}

	filename := filepath.Base(src)
	log.Printf("filename: %s", filename)

	log.Printf("destination: %v", dst)
	os.MkdirAll(dst, 0744)

	dst, err = filepath.Abs(filepath.Join(dst, filename))
	if err != nil {
		log.Fatalf("absolute %s: %v", dst, err)
	}
	log.Printf("destination: %v", dst)

	if err = CopyFile(src, dst); err != nil {
		log.Fatal(err)
	}
}

// CopyFile copies files from src to destination.
func CopyFile(src, dst string) error {
	//TODO copy only images
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	//	tm, err := getCreationDatetime(src)
	//	if err != nil {
	//		return err
	//	}
	//
	//filename := fmt.Sprintf("%v-%v-%v-%v", tm.Year(), tm.Month(), tm.Day(), tm.Minute())
	//path := filepath.Join(dst, filename)

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

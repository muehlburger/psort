package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

func main() {
	filename := "20180805_115407.jpg"
	src := "./testfiles/" + filename
	dst := "./testfiles/" + time.Now().Format("2006-01-02_150405") + "_" + filename

	err := CopyFile(src, dst)
	if err != nil {
		log.Fatal(err)
	}
}

// CopyFile copies files from src to destination.
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

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

	return os.Rename(tmp.Name(), dst)
}

// CreationTime extracts the creation Date of the given File.
func CreationTime(in *os.File) time.Time {
	x, err := exif.Decode(in)
	if err != nil {
		log.Fatal(err)
	}

	tm, err := x.DateTime()
	if err != nil {
		log.Fatal(err)
	}
	return tm
}

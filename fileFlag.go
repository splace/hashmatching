package main

import "os"

type fileValue struct {
	*os.File
}

func (fsf *fileValue) Set(v string) error {
	f, err := os.Open(v)
	fsf.File = f
	return err
}

func (fsf *fileValue) String() string {
	if fsf == nil || fsf.File == nil {
		return "<nil>"
	}
	return fsf.File.Name()
}

type newFileValue struct {
	fileValue
}

func (fsf *newFileValue) Set(v string) error {
	f, err := os.Create(v)
	fsf.File = f
	return err
}


package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/vidhanio/fson/errors"
)

type FSONType int

const (
	FSONTypeFile FSONType = iota
	FSONTypeObject
	FSONTypeArray
)

type FSON struct {
	Name     string
	Index    int
	Value    string
	FSONType FSONType
	Children []*FSON
	File     *FSON
	Parent   *FSON
}

func main() {

	sampleFSON := &FSON{
		Name:     "sample",
		FSONType: FSONTypeObject,
	}
	sampleFSON.NewNamedChild("vidhan", FSONTypeArray, "test").NewIndexedChild(FSONTypeFile, "test")

	err := sampleFSON.Write(".")
	if err != nil {
		panic(err)
	}
}

func New(path string) (*FSON, error) {
	fson := new(FSON)

	err := filepath.WalkDir(path,
		func(path string, d fs.DirEntry, err error) error {

			if err != nil {
				return err
			}

			if d.IsDir() {
				fson.NewNamedChild(d.Name(), getFolderType(path), "")
			} else {
				fson.NewIndexedChild(FSONTypeFile, d.Name())
			}

			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	return fson, nil
}

func (fson *FSON) Write(path ...string) error {
	fmt.Printf("%+v\n", fson)
	if fson.Parent == nil {
		path = append(path, fson.Name)
	} else {
		if fson.Parent.FSONType == FSONTypeArray {
			path = append(path, strconv.Itoa(fson.Index))
		} else {
			name := fson.Name
			if fson.FSONType == FSONTypeArray {
				name += "_"
			}
			path = append(path, name)
		}
	}

	if fson.FSONType == FSONTypeFile {
		f, err := os.Create(filepath.Join(path...))
		if err != nil {
			return err
		}

		defer f.Close()

		_, err = f.WriteString(fson.Value)
		if err != nil {
			return err
		}
	} else if fson.FSONType == FSONTypeObject {
		err := os.Mkdir(filepath.Join(path...), fs.ModePerm)
		if err != nil {
			return err
		}

		for _, child := range fson.Children {
			err = child.Write(path...)
			if err != nil {
				return err
			}
		}
	} else if fson.FSONType == FSONTypeArray {
		err := os.Mkdir(filepath.Join(path...), fs.ModePerm)
		if err != nil {
			return err
		}

		for _, child := range fson.Children {
			err = child.Write(path...)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (fson *FSON) Get(path ...string) (*FSON, error) {
	nF := *fson
	newFSON := &nF

	var err error

	path = append(path, "")

	for i, p := range path {
		if i == len(path)-1 {
			return newFSON, nil

		} else if newFSON.FSONType == FSONTypeObject {
			newFSON, err = newFSON.GetNamedChild(p)
			if err != nil {
				return nil, err
			}

		} else if newFSON.FSONType == FSONTypeArray {
			index, err := strconv.Atoi(p)
			if err != nil {
				return nil, err
			}

			newFSON, err = newFSON.GetIndexedChild(index)
			if err != nil {
				return nil, err
			}

		} else if newFSON.FSONType == FSONTypeFile {
			return nil, errors.ErrCannotAccessFileChildren
		}
	}

	return nil, nil
}

func (fson *FSON) GetNamedChild(name string) (*FSON, error) {
	if fson.FSONType != FSONTypeObject {
		return nil, errors.ErrNotAnObject
	}

	for _, child := range fson.Children {
		if child.Name == name {
			return child, nil
		}
	}

	return nil, nil
}

func (fson *FSON) GetIndexedChild(index int) (*FSON, error) {
	if fson.FSONType != FSONTypeArray {
		return nil, errors.ErrNotAnArray
	}

	if index >= len(fson.Children) {
		return nil, errors.ErrIndexOutOfBounds
	}

	return fson.Children[index], nil
}

/*
Creates a child for an FSON Object.
Make sure to assert that fson.FSONType == FSONTypeObject

	if fson.FSONType == FSONTypeObject {
		fson.NewNamedChild("fson name", FSONTypeFile, "fson value")
	}

Value is ignored if fsonType != FSONTypeFile
*/
func (fson *FSON) NewNamedChild(name string, fsonType FSONType, value string) *FSON {

	newFSON := &FSON{
		Name:     name,
		FSONType: fsonType,
		Parent:   fson,
	}

	if fsonType == FSONTypeFile {
		newFSON.Value = value
	}

	fson.Children = append(fson.Children, newFSON)

	return newFSON
}

/*
Creates a child for an FSON Array.
Make sure to assert that fson.FSONType == FSONTypeArray

	if fson.FSONType == FSONTypeArray {
		fson.NewIndexedChild(FSONTypeFile, "fson value")
	}

Value is ignored if fsonType != FSONTypeFile
*/
func (fson *FSON) NewIndexedChild(fsonType FSONType, value string) *FSON {
	newFSON := &FSON{
		Index:    len(fson.Children),
		FSONType: fsonType,
		Parent:   fson,
	}

	if fsonType == FSONTypeFile {
		newFSON.Value = value
	}

	fson.Children = append(fson.Children, newFSON)

	return newFSON
}

// Get the FSONType from the folder name.
func getFolderType(path string) FSONType {
	path = filepath.Base(path)

	if strings.HasSuffix(path, "_") {
		return FSONTypeArray
	}

	return FSONTypeObject
}

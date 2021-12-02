package fson

import (
	"io/fs"
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

func (fson *FSON) NewNamedChild(name string, fsonType FSONType, value string) (*FSON, error) {
	if fson.FSONType != FSONTypeObject {
		return nil, errors.ErrNotAnObject
	}

	newFSON := &FSON{
		Name:     name,
		FSONType: fsonType,
		Parent:   fson,
	}

	if fsonType == FSONTypeFile {
		newFSON.Value = value
	}

	fson.Children = append(fson.Children, newFSON)

	return newFSON, nil
}

func (fson *FSON) NewIndexedChild(fsonType FSONType, value string) (*FSON, error) {
	if fson.FSONType != FSONTypeArray {
		return nil, errors.ErrNotAnArray
	}

	newFSON := &FSON{
		Index:    len(fson.Children),
		FSONType: fsonType,
		Parent:   fson,
	}

	if fsonType == FSONTypeFile {
		newFSON.Value = value
	}

	fson.Children = append(fson.Children, newFSON)

	return newFSON, nil
}

// Get the FSONType from the folder name.
func getFolderType(path string) FSONType {
	path = filepath.Base(path)

	if strings.HasSuffix(path, "_") {
		return FSONTypeArray
	}

	return FSONTypeObject
}

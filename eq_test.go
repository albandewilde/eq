package main

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestEq(t *testing.T) {
	cases := []struct {
		f0       string
		f1       string
		expected bool
		err      string
	}{
		{
			f0:       "files/file0",
			f1:       "./files/file0",
			expected: true,
			err: fmt.Sprintf(
				"Files `%s` and `%s` are considered not equal",
				"files/file0",
				"./files/file0",
			),
		},
		{
			f0:       "files/file1",
			f1:       "files/file2",
			expected: true,
			err: fmt.Sprintf(
				"Files `%s` and `%s` are considered not equal",
				"files/file1",
				"files/file2",
			),
		},
		{
			f0:       "files/file0",
			f1:       "./files/file1",
			expected: false,
			err: fmt.Sprintf(
				"Files `%s` and `%s` are considered equal",
				"files/file0",
				"./files/file1",
			),
		},
		{
			f0:       "files/file0",
			f1:       "files/file2",
			expected: false,
			err: fmt.Sprintf(
				"Files `%s` and `%s` are considered equal",
				"files/file0",
				"files/file2",
			),
		},
	}

	for _, c := range cases {
		isEq, err := Eq(c.f0, c.f1)
		if err != nil {
			t.Errorf("An error has been raised: %s", err.Error())
		}
		if isEq != c.expected {
			t.Errorf(c.err)
		}
	}
}

func TestEqWithMissingFile(t *testing.T) {
	cases := []struct {
		f0       string
		f1       string
		expected bool
		err      string
	}{
		{
			f0:       "files/file4",
			f1:       "./files/file0",
			expected: true,
			err:      "open files/file4: no such file or directory",
		},
		{
			f0:       "./files/file4",
			f1:       "files/file2",
			expected: true,
			err:      "open ./files/file4: no such file or directory",
		},
		{
			f0:       "files/file0",
			f1:       "./files/file4",
			expected: false,
			err:      "open ./files/file4: no such file or directory",
		},
		{
			f0:       "files/file2",
			f1:       "files/file4",
			expected: false,
			err:      "open files/file4: no such file or directory",
		},
		{
			f0:       "files/file1",
			f1:       "files/file5",
			expected: false,
			err:      "open files/file5: permission denied",
		},
	}

	for _, c := range cases {
		isEq, err := Eq(c.f0, c.f1)
		if err.Error() != c.err {
			t.Errorf(
				"The wrong error has been raised, expected `%s` but got `%s`",
				c.err,
				err.Error(),
			)
		}
		if isEq {
			t.Errorf(
				fmt.Sprintf(
					"Files `%s` and `%s` are considered equal",
					c.f0,
					c.f1,
				),
			)
		}
	}
}

func TestSameFileInDir(t *testing.T) {
	cases := []struct {
		dir      string
		expected [][]string
		err      error
	}{
		{
			dir:      "./files/",
			expected: [][]string{{"files/file1", "files/file2", "files/file6"}},
			err:      nil,
		},
		{
			dir:      "./fls/",
			expected: [][]string{},
			err:      errors.New("open ./fls/: no such file or directory"),
		},
		{
			dir:      "./files/other_files/",
			expected: [][]string{},
			err:      nil,
		},
	}

	for _, c := range cases {
		duplicates, err := Duplicates(c.dir)
		if !(c.err == nil && err == nil) && err.Error() != c.err.Error() {
			t.Errorf(
				"An error occured, expected `%s` but got `%s`",
				c.err,
				err.Error(),
			)
		}

		if !(len(duplicates) == 0 && len(c.expected) == 0) && !reflect.DeepEqual(c.expected, duplicates) {
			t.Errorf(
				"Failed finding equal files, expected: `%s` but got: `%s`",
				c.expected,
				duplicates,
			)
		}
	}
}

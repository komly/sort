package main

import (
	"testing"
	"os"
	"io/ioutil"
	"strings"
	"bytes"
	"io"
)

func TestWriteToTempFile(t *testing.T) {
	f, err := writeToTempFile([]string{"A", "B", "C"})
	if err != nil {
		t.Fatalf("can't write: %v", err)
	}
	defer os.Remove(f)

	data, err := ioutil.ReadFile(f)
	if err != nil {
		t.Fatalf("can't read file: %v", err)
	}
	expected := "A\nB\nC\n"
	if string(data) != expected {
		t.Fatalf("content shoud be equal, got: %q", string(data))
	}
}

func TestMergeFiles(t *testing.T)  {
	for _, test := range []struct {
		In []io.Reader
		Out string
	}{
		{
			In: []io.Reader{
				strings.NewReader(""),
			},
			Out: "",
		},

		{
			In: []io.Reader{
				strings.NewReader("A\n"),
			},
			Out: "A\n",
		},
		{
			In: []io.Reader{
				strings.NewReader("A\nB\nC\n"),
			},
			Out: "A\nB\nC\n",
		},
		{
			In: []io.Reader{
				strings.NewReader("A\nB\nC"),
			},
			Out: "A\nB\nC\n",
		},
		{
		In: []io.Reader{
			strings.NewReader("B"),
			strings.NewReader("C"),
			strings.NewReader("A"),
			},
		Out: "A\nB\nC\n",
		},

		{
			In: []io.Reader{
				strings.NewReader("B\nC"),
				strings.NewReader("A"),
			},
			Out: "A\nB\nC\n",
		},
	} {
		res := bytes.NewBufferString("")
		if err := mergeReaders(test.In, res); err != nil {
			t.Fatalf("mergeReader error: %v", err)
		}
		if res.String() != test.Out {
			t.Fatalf("output should be equal, want: %q, got: %q", test.Out, res.String())
		}
	}

}


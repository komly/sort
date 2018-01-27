package main

import (
	"testing"
	"os"
	"io/ioutil"
)

func TestWriteSorted(t *testing.T) {
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
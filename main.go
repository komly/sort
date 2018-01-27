package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

var limit = flag.Int("chunk_size", 0, "chunk size in bytes")

func createSortedFiles(r io.Reader) ([]*os.File, error) {
	s := bufio.NewScanner(r)

	buffer := make([]string, 0)
	LIMIT := 300 * 1024 * 1024 // 100m
	read := 0
	files := make([]*os.File, 0)

	for s.Scan() {
		line := s.Text()
		buffer = append(buffer, line)
		read += len(line)

		if read >= LIMIT {
			sort.Strings(buffer)
			tempFile, err := ioutil.TempFile("", "chunk")
			if err != nil {
				return nil, fmt.Errorf("can't create temp file: %v", err)
			}
			for _, l := range buffer {
				fmt.Fprintf(tempFile, "%s\n", l)
			}

			if _, err := tempFile.Seek(0, 0); err != nil {
				os.Remove(tempFile.Name())
				return nil, fmt.Errorf("can't seek file: %v", err)
			}

			files = append(files, tempFile)
			buffer = make([]string, 0)
			read = 0
		}
	}

	if len(buffer) > 0 {
		sort.Strings(buffer)
		tempFile, err := ioutil.TempFile("", "chunk")
		if err != nil {
			return nil, fmt.Errorf("can't create temp file: %v", err)
		}
		for _, l := range buffer {
			fmt.Fprintf(tempFile, "%s\n", l)
		}

		if _, err := tempFile.Seek(0, 0); err != nil {
			os.Remove(tempFile.Name())
			return nil, fmt.Errorf("can't seek file: %v", err)
		}
		files = append(files, tempFile)
	}

	return files, nil
}

func main() {
	flag.Parse()
	if *limit == 0 {
		*limit = 10 * 1024 * 1024 //TODO: heuristic from runtime.Memstats or /proc
	}

	files, err := createSortedFiles(os.Stdin)
	if err != nil {
		log.Printf("can't create sorted files")
	}

	scanners := make(map[*bufio.Scanner]struct{})

	for _, f := range files {
		defer os.Remove(f.Name())
		scanners[bufio.NewScanner(f)] = struct{}{}
	}

	if len(scanners) <= 0 {
		return
	}

	for s := range scanners {
		if !s.Scan() {
			delete(scanners, s)
		}
	}

	for {
		var minScanner *bufio.Scanner

		for s := range scanners {
			minScanner = s
			break
		}

		for s := range scanners {
			if strings.Compare(s.Text(), minScanner.Text()) <= 0 {
				minScanner = s
			}
		}
		fmt.Printf("%s\n", minScanner.Text())
		if !minScanner.Scan() {
			delete(scanners, minScanner)
		}

		if len(scanners) <= 0 {
			break
		}

	}

}

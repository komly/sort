package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sort"
	"io/ioutil"
)

var limit = flag.Int("chunk_size", 0, "chunk size in bytes")
var verbose = flag.Bool("verbose", false, "verbose mode")


func writeToTempFile(lines []string) (string, error) {
	tempFile, err := ioutil.TempFile("", "chunk")
	if err != nil {
		return "", fmt.Errorf("can't create temp file: %v", err)
	}
	for _, l := range lines {
		fmt.Fprintf(tempFile, "%s\n", l)
	}

	return tempFile.Name(), nil
}


func chunkReader(r io.Reader, limit int) (chan []string) {
	res := make(chan []string)
	go func() {
		s := bufio.NewScanner(r)

		buffer := make([]string, 0)
		read := 0

		for s.Scan() {
			line := s.Text()
			buffer = append(buffer, line)
			read += len(line)

			if read >= limit {
				res <- buffer
				buffer = make([]string, 0)
				read = 0
			}
		}

		if len(buffer) > 0 {
			res <- buffer
		}

		close(res)
	}()

	return res

}

func main() {
	flag.Parse()
	if *limit == 0 {
		*limit = 10 * 1024 * 1024 //TODO: heuristic from runtime.Memstats or /proc
	}

	files := make([]string, 0)

	for chunk := range chunkReader(os.Stdin, *limit) {
		if *verbose {
			log.Printf("read chunk: %+v", chunk)
		}

		sort.Strings(chunk)
		fileName, err := writeToTempFile(chunk)
		if err != nil {
			log.Printf("can't write chunk file: %v", err)
			return
		}

		if *verbose {
			log.Printf("chunk filename: %v", fileName)
		}
		defer os.Remove(fileName)

		files = append(files, fileName)
	}

	scanners := make(map[*bufio.Scanner]struct{})
	for _, fileName := range files {
		f, err := os.Open(fileName)
		if err != nil {
			log.Printf("can't open chunk file: %v", err)
			return
		}
		defer f.Close()

		scanners[bufio.NewScanner(f)] = struct{}{}
	}



	for s := range scanners {
		if !s.Scan() {
			delete(scanners, s)
		}
	}

	if len(scanners) <= 0 {
		return
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

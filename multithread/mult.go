package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

func readCsv(filePath string) ([][]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}

	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}
	return records, err
}

func parsingList(a [][]string, g chan []string) {
	for _, l := range a {
		g <- l
	}
	close(g)
}

func parsingRow(e chan []string, q chan []int, errCh chan error) {
	defer close(q)
	defer close(errCh)
	for w := range e {
		i := make([]int, 0, len(w))
		for _, f := range w {
			s, err := strconv.Atoi(f)
			if err != nil {
				errCh <- err
			}
			i = append(i, s)
		}
		q <- i
	}
}

func summing(o chan []int, b chan int) {
	for d := range o {
		u := 0
		for _, s := range d {
			u += s
		}
		b <- u
	}
	close(b)
}

func converting(j chan int, z chan string) {
	m := make([]string, 0)
	for d := range j {
		s := strconv.Itoa(d)
		m = append(m, s)
		z <- s
	}
	close(z)
}

func writing(p string, v chan string, n chan string) {
	file, err := os.Create(p)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}

	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	w := csv.NewWriter(file)
	defer w.Flush()
	for rec := range v {
		if err := w.Write([]string{rec}); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}
	close(n)
}

func main() {
	ch := make(chan []string)
	ch2 := make(chan []int)
	chErr := make(chan error)
	chOutR := make(chan []int)
	ch3 := make(chan int)
	ch4 := make(chan string)
	ch5 := make(chan string)
	p := "The_sums"

	number, err := readCsv("multithread/test_100.csv")
	if err != nil {
		fmt.Println("Error when reading the file", err)
	} else {
		fmt.Println("Processing started")
	}

	go parsingList(number, ch)
	go parsingRow(ch, ch2, chErr)

	go func() {
		defer close(chOutR)
		for {
			select {
			case v := <-ch2:
				if v == nil {
					return
				}
				chOutR <- v
			case errFromCh := <-chErr:
				if errFromCh == nil {
					return
				}
				log.Fatal("Error caught: ", errFromCh)
			}
		}
	}()
	go summing(chOutR, ch3)
	go converting(ch3, ch4)
	go writing(p, ch4, ch5)
	<-ch5

	fmt.Println("Done")
}

package main

import (
	//"bytes"
	//"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
)

func TestReadCsv(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		output   [][]string
	}{
		{
			name:     "reading file",
			filePath: "test_2.csv",
			output:   [][]string{{"10", "11", "12"}, {"20", "21", "22"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, _ := readCsv(tt.filePath)
			assert.Equal(t, tt.output, res, "readCsv returns unexpected value")
		})
	}
}
func TestParsingList(t *testing.T) {
	tests := []struct {
		name    string
		input   [][]string
		output  []string
		output2 []string
	}{
		{
			name:    "correct parsing",
			input:   [][]string{{"10", "11", "12"}, {"20", "21", "22"}},
			output:  []string{"10", "11", "12"},
			output2: []string{"20", "21", "22"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := make(chan []string)
			go parsingList(tt.input, ch)
			assert.Equal(t, tt.output, <-ch, "parsingList returns unexpected value")
			assert.Equal(t, tt.output2, <-ch, "parsingList returns unexpected value")
		})
	}
}
func TestParsingRow(t *testing.T) {
	tests := []struct {
		name      string
		input     []string
		output    []int
		outputErr string
	}{
		{
			name:   "parsing row",
			input:  []string{"10", "11", "12"},
			output: []int{10, 11, 12},
		},
		{
			name:      "error parsing row",
			input:     []string{"10", "11", "df"},
			output:    []int{10, 11, 0},
			outputErr: "strconv.Atoi: parsing \"df\": invalid syntax",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chIn    := make(chan []string)
			chOut   := make(chan []int)
			chOut1  := make(chan []int)
			chErr   := make(chan error)

			go parsingRow(chIn, chOut, chErr)
			chIn <- tt.input
			defer close(chIn)

			go func() {
				for {
					select {
					case v := <-chOut:
						if v == nil {
							return
						}
						chOut1 <- v
					case errFromCh := <-chErr:
						if errFromCh == nil {
							return
						}
						assert.Equal(t, tt.outputErr, errFromCh.Error(), "parsingRow returns unexpected error")
					}
				}
			}()
			assert.Equal(t, tt.output, <-chOut1, "parsingRow returns unexpected value")
		})
	}
}
func TestSumming(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		output int
	}{
		{
			name:   "correct summing",
			input:  []int{10, 11, 12},
			output: 33,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chIn1 := make(chan []int)
			chOut1 := make(chan int)

			go summing(chIn1, chOut1)
			chIn1 <- tt.input
			close(chIn1)
			assert.Equal(t, tt.output, <-chOut1, "summing returns unexpected value")
		})
	}
}
func TestConverting(t *testing.T) {
	tests := []struct {
		name   string
		input  int
		output string
	}{
		{
			name:   "correct converting",
			input:  33,
			output: "33",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chIn2  := make(chan int)
			chOut2 := make(chan string)

			go converting(chIn2, chOut2)
			chIn2 <- tt.input
			close(chIn2)

			assert.Equal(t, tt.output, <-chOut2, "converting returns unexpected value")
		})
	}
}
func TestWriting(t *testing.T) {
	tests := []struct {
		name     string
		actual   string
		input    string
		input2   string
		expected string
	}{
		{
			name:     "correct writing",
			actual:   "Test_writing_actual.csv",
			input:    "33",
			input2:   "63",
			expected: "/Users/vantihovich/work_lyft/tasks/multithread/Test_write_expected.csv",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chIn3  := make(chan string)
			chOut3 := make(chan string)

			go writing(tt.actual, chIn3, chOut3)
			chIn3 <- tt.input
			chIn3 <- tt.input2
			close(chIn3)
			<-chOut3

			sfExpected, err := ioutil.ReadFile(tt.expected)
			if err != nil {
				log.Fatal(err)
			}
			dfActual, err := ioutil.ReadFile(tt.actual)
			if err != nil {
				log.Fatal(err)
			}

			assert.Equal(t, sfExpected, dfActual, "writing returns unexpected value")
		})
	}
}

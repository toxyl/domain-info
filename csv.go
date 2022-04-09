package main

import (
	"fmt"
	"math"
	"strings"
)

type CSV struct {
	Headers []string
	Formats []string
	Rows    [][]string
}

type CSVHeader struct {
	Name   string
	Format string
}

func (csv *CSV) AddHeaders(header ...CSVHeader) {
	for _, h := range header {
		csv.Headers = append(csv.Headers, h.Name)
		csv.Formats = append(csv.Formats, h.Format)
	}
}

func (csv *CSV) AddRow(fields ...string) {
	if len(fields) != len(csv.Headers) {
		fmt.Printf("Could not add row, expected %d fields, but got %d fields. (%+v)\n", len(csv.Headers), len(fields), fields)
		return
	}
	csv.Rows = append(csv.Rows, fields)
}

func (csv *CSV) String() string {
	res := ""
	maxLengths := make([]int, len(csv.Headers))
	for _, r := range csv.Rows {
		for i, v := range r {
			maxLengths[i] = int(math.Max(float64(maxLengths[i]), float64(len(v))+1.0))
			maxLengths[i] = int(math.Max(float64(maxLengths[i]), float64(len(csv.Headers[i]))+1.0))
		}
	}

	res += "\033[2K"
	for i, v := range csv.Headers {
		f := csv.Formats[i]
		if strings.Contains(f, "%-") {
			f = strings.ReplaceAll(f, "%-", "%-"+fmt.Sprint(maxLengths[i]-1))
		} else {
			f = strings.ReplaceAll(f, "%", "%"+fmt.Sprint(maxLengths[i]-1))
		}
		res += fmt.Sprintf(f+" ", v)
	}
	res += "\n"

	res += "\033[2K"
	for i := range csv.Headers {
		res += strings.Repeat("-", maxLengths[i])
	}
	res += "\n"

	for _, r := range csv.Rows {
		res += "\033[2K"
		for i, v := range r {
			f := csv.Formats[i]
			if strings.Contains(f, "%-") {
				f = strings.ReplaceAll(f, "%-", "%-"+fmt.Sprint(maxLengths[i]-1))
			} else {
				f = strings.ReplaceAll(f, "%", "%"+fmt.Sprint(maxLengths[i]-1))
			}
			res += fmt.Sprintf(f+" ", v)
		}
		res += "\n"
	}
	return res
}

func NewCSV() *CSV {
	return &CSV{
		Headers: []string{},
		Formats: []string{},
		Rows:    [][]string{},
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func FlushAndInit() {
	_, err := os.Stat("Compressed_Circuits")
	if os.IsNotExist(err) {
		fmt.Println("Compressed circuits folder doesnt exists making one now")
		os.Mkdir("Compressed_Circuits", 0755)
	}
	os.RemoveAll("Compressed_Circuits")
	os.Mkdir("Compressed_Circuits", 0755)
	WriteCircuits([]Circuit{{1}}, 1)
}

func Generate_Combinations(length int) {
	abort_ch := make(chan struct{})
	known_resistances := make(map[Circuit]bool)
	for i := 1; i < length; i++ {
		j := length - i
		if j < i {
			break
		}
		first_file_ch := ReadCircuits(i, abort_ch)
		for first_buffer := range first_file_ch {
			circuits := []Circuit{}
			second_file_ch := ReadCircuits(j, abort_ch)
			for second_buffer := range second_file_ch {
				for _, f := range first_buffer {
					for _, g := range second_buffer {
						series := f.AddSeries(g)
						if !known_resistances[series] {
							known_resistances[series] = true
							circuits = append(circuits, series)
						}
						parallel := f.AddParallel(g)
						if !known_resistances[parallel] {
							known_resistances[parallel] = true
							circuits = append(circuits, parallel)
						}
					}
				}
			}
			WriteCircuits(circuits, length)
		}

	}
}

func Decompose(target float64, length int) string {
	if target == 1 && length == 1 {
		return "1"
	}
	abort_ch := make(chan struct{})

	for i := length / 2; i < length; i++ {
		j := length - i
		first_file_ch := ReadCircuits(i, abort_ch)
		for first_buffer := range first_file_ch {
			second_file_ch := ReadCircuits(j, abort_ch)
			for second_buffer := range second_file_ch {
				for _, f := range first_buffer {
					for _, g := range second_buffer {
						series := f.AddSeries(g)
						parallel := f.AddParallel(g)
						if series.Resistance == target {
							f_comp := Decompose(f.Resistance, i)
							g_comp := Decompose(g.Resistance, j)
							out := ""
							if f_comp != "1" {
								out += "(" + f_comp + ") "
							} else {
								out += "1 "
							}
							out += "+"
							if g_comp != "1" {
								out += " (" + g_comp + ")"
							} else {
								out += " 1"
							}
							return out
						}
						if parallel.Resistance == target {
							f_comp := Decompose(f.Resistance, i)
							g_comp := Decompose(g.Resistance, j)
							out := ""
							if f_comp != "1" {
								out += "(" + f_comp + ") "
							} else {
								out += "1 "
							}
							out += "*"
							if g_comp != "1" {
								out += " (" + g_comp + ")"
							} else {
								out += " 1"
							}
							return out
						}
					}
				}
			}
		}

	}
	return "No Valid Decomposition found"
}

func Search(value float64) (string, bool) {
	i := 1
	for {

		path := filepath.FromSlash("Compressed_Circuits/" + strconv.Itoa(i) + ".circ")
		_, err := os.Stat(path)

		if os.IsNotExist(err) {
			break
		}

		abort_ch := make(chan struct{})
		value_chanel := ReadCircuits(i, abort_ch)
		position := 0
		for buffer := range value_chanel {
			for _, circuit := range buffer {
				if circuit.Resistance == value {
					return strconv.Itoa(i) + "-" + strconv.Itoa(position), true
				}
				position++
			}
		}
		i++
	}
	return "", false
}

func main() {
	for {
		reader := bufio.NewScanner(os.Stdin)
		for reader.Scan() {
			fields := strings.Fields(string(reader.Text()))
			switch fields[0] {
			case "quit":
				return
			case "flush":
				FlushAndInit()
			case "generate":
				i, err := strconv.ParseInt(fields[1], 10, 64)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("GENERATING")
				now := time.Now()
				Generate_Combinations(int(i))
				taken := time.Now().Sub(now)
				fmt.Println("Done: ", taken)
			case "generateto":
				i, err := strconv.ParseInt(fields[1], 10, 64)
				if err != nil {
					fmt.Println(err)
					continue
				}
				FlushAndInit()
				fmt.Println("GENERATING")
				for to_gen := 2; int64(to_gen) <= i; to_gen++ {
					Generate_Combinations(int(to_gen))
				}
				fmt.Println("Done")

			case "search":
				f := 0.0
				if strings.Contains(fields[1], "/") {
					parts := strings.Split(fields[1], "/")
					f1, _ := strconv.ParseFloat(parts[0], 64)
					f2, _ := strconv.ParseFloat(parts[1], 64)
					f = f1 / f2
				} else {
					f, _ = strconv.ParseFloat(fields[1], 64)
				}
				fmt.Println("Searching")
				result, found := Search(f)
				if !found {
					fmt.Println("Not Found")
				} else {
					fmt.Println(result)
				}
			case "decomp":
				f := 0.0
				if strings.Contains(fields[1], "/") {
					parts := strings.Split(fields[1], "/")
					f1, _ := strconv.ParseFloat(parts[0], 64)
					f2, _ := strconv.ParseFloat(parts[1], 64)
					f = f1 / f2
				} else {
					f, _ = strconv.ParseFloat(fields[1], 64)
				}
				no_resistors, _ := strconv.ParseInt(fields[2], 10, 8)
				result := Decompose(f, int(no_resistors))
				fmt.Println(result)
			}
		}
	}
}

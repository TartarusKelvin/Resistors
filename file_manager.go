package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
)

func Float64ToByte(f float64) []byte {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], math.Float64bits(f))
	return buf[:]
}

func ByteToFloat64(bytes []byte) float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(bytes))
}

func WriteCircuits(circs []Circuit, length int) {
	path := filepath.FromSlash("Compressed_Circuits/" + strconv.Itoa(length) + ".circ")
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	bytes := []byte{}
	for _, a := range circs {
		bytes = append(bytes, Float64ToByte(a.Resistance)...)
	}
	_, err = f.Write(bytes)
	if err != nil {
		fmt.Println(err)
		return
	}
}

const max_read int = 1000000

func ReadCircuits(length int, abort_chanel <-chan struct{}) <-chan []Circuit {
	// Each length of circuit has its own file to save us from
	// searching each time. Returns a channel that will provide
	// a series of compiled circuits. To stop the reading process pass a value into
	// the abort channel
	ch := make(chan []Circuit)
	go func() {
		defer close(ch)
		// Start a subroutine
		path := filepath.FromSlash("Compressed_Circuits/" + strconv.Itoa(length) + ".circ")
		file, err := os.Open(path)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		buffer := make([]byte, 8)

		should_break := false
		for {
			if should_break {
				break
			}
			circs := []Circuit{}
			for i := 0; i < max_read; i++ {
				_, err := file.Read(buffer)
				if err != nil {
					if err != io.EOF {
						fmt.Println(err)
					}
					should_break = true
					break
				}
				next_circ := Circuit{Resistance: ByteToFloat64(buffer)}

				circs = append(circs, next_circ)
			}
			select {
			case ch <- circs:
			case <-abort_chanel:
				fmt.Println("Aborting read")
				return
			}
		}
	}()
	return ch
}

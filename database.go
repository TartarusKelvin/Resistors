package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
)

/*
We sadly cant keep results in memory past ~20 circuits, as such
we need to store results on the hardrive when not in use. Further
we also expect that the total results may get to large to load in
at once so we need to be able to read in as we need without the chance
that the database gets to large
*/

const max_read int = 1048576 /// try and keep this a power of 2

func read_as_float(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}

func float32ToByte(f float32) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, f)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}

func ReadCircuits(length int, abort_chanel <-chan struct{}) <-chan []CompiledCircuit {
	// Each length of circuit has its own file to save us from
	// searching each time. Returns a channel that will provide
	// a series of compiled circuits. To stop the reading process pass a value into
	// the abort channel
	ch := make(chan []CompiledCircuit)
	go func() {
		defer close(ch)
		// Start a subroutine
		file, err := os.Open("Compressed_Circuits\\" + strconv.Itoa(length) + ".circ")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		buffer := make([]byte, 4)

		should_break := false
		for {
			if should_break {
				break
			}
			circs := []CompiledCircuit{}
			for i := 0; i < max_read; i++ {
				_, err := file.Read(buffer)
				if err != nil {
					if err != io.EOF {
						fmt.Println(err)
					}
					should_break = true
					break
				}
				next_circ := CompiledCircuit{Resistance: read_as_float(buffer), TotalResistors: length}
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

func WriteCircuit(c []CompiledCircuit) bool {
	f, err := os.OpenFile("Compressed_Circuits\\"+strconv.Itoa(c[0].TotalResistors)+".circ", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		return false
	}
	bytes := []byte{}
	for _, a := range c {
		bytes = append(bytes, float32ToByte(a.Resistance)...)
	}
	_, err = f.Write(bytes)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

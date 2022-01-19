package main

import "fmt"

var CircuitDatabase map[int][]CompiledCircuit

func InitDatabase() {
	CircuitDatabase = make(map[int][]CompiledCircuit)
	CircuitDatabase[1] = []CompiledCircuit{{Resistance: 1, TotalResistors: 1}}
}

func AddCircuitsToDatabase(circs []CompiledCircuit, length int) {
	WriteCircuit(circs)
}

func GetCombinations(target int) <-chan [][]CompiledCircuit {
	// Generates all possible 2 circuit combinations that use the target number of resistors
	ch := make(chan [][]CompiledCircuit)

	go func() {
		defer close(ch)
		for i := target - 1; i >= 1; i-- {
			j := target - i
			if j > i {
				break
			}
			abort := make(chan struct{})
			first_circuits_ch := ReadCircuits(i, abort)
			for first_circuits := range first_circuits_ch {
				second_circuits_ch := ReadCircuits(j, abort)
				for second_circuits := range second_circuits_ch {
					combinations := [][]CompiledCircuit{}
					for _, f := range first_circuits {
						for _, s := range second_circuits {
							combinations = append(combinations, []CompiledCircuit{f, s})
						}
					}
					ch <- combinations
				}
			}
		}
	}()
	return ch
}

func GenerateCircuits(length int) {
	combinations_ch := GetCombinations(length)
	achieved_resistances := make(map[float32]bool)
	for combinations := range combinations_ch {
		circuits := []CompiledCircuit{}
		fmt.Print(".")
		for _, combination := range combinations {
			series_circuit := MakeSeriesCircuit(combination[0], combination[1])
			resistance := series_circuit.Resistance()
			if !achieved_resistances[resistance] {
				achieved_resistances[resistance] = true
				circuits = append(circuits, series_circuit.Compile())
			}
			parallel_circuit := MakeParallelCircuit(combination[0], combination[1])
			resistance = parallel_circuit.Resistance()
			if !achieved_resistances[resistance] {
				achieved_resistances[resistance] = true
				circuits = append(circuits, parallel_circuit.Compile())
			}
		}
		AddCircuitsToDatabase(circuits, length)
	}
}

func main() {
	InitDatabase()
	for i := 23; i < 30; i++ {
		fmt.Println(i, ")")
		GenerateCircuits(i)
		//fmt.Println("	- New_Circuits: ", len(new_circuits))
	}
	//fmt.Println(GenerateCircuits(2))
}


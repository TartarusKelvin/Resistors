package main

type CompiledCircuit struct {
	Resistance     float32
	TotalResistors int
}

type Circuit struct {
	Series [][]CompiledCircuit
}

func MakeParallelCircuit(top CompiledCircuit, bot CompiledCircuit) Circuit {
	// Should this just return a Compiled Circuit? We never actually need the normal circuit
	// at which point why do we even have the normal circuit if we arent using it?
	return Circuit{Series: [][]CompiledCircuit{{top, bot}}}
}

func MakeSeriesCircuit(first CompiledCircuit, second CompiledCircuit) Circuit {
	return Circuit{Series: [][]CompiledCircuit{{first}, {second}}}
}

func (circ Circuit) Resistance() float32 {
	resistance := float32(0.0)
	for _, block := range circ.Series {
		if len(block) == 1 {
			resistance += block[0].Resistance
		} else {
			recip_resistance := float32(0.0)
			for _, c := range block {
				recip_resistance += 1 / c.Resistance
			}
			resistance += 1 / recip_resistance
		}
	}
	return resistance
}

func (circ Circuit) Compile() CompiledCircuit {
	total_resistors := 0
	for _, block := range circ.Series {
		for _, c := range block {
			total_resistors += c.TotalResistors
		}
	}
	return CompiledCircuit{TotalResistors: total_resistors, Resistance: circ.Resistance()}
}

package main

type Circuit struct {
	Resistance float64
}

func (a Circuit) AddSeries(b Circuit) Circuit {
	return Circuit{a.Resistance + b.Resistance}
}

func (a Circuit) AddParallel(b Circuit) Circuit {
	return Circuit{a.Resistance * b.Resistance / (a.Resistance + b.Resistance)}
}

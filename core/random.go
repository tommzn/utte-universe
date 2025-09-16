package core

import "math/rand"

type BuiltInRand struct{}

func (r *BuiltInRand) Seek() float64 {
	return rand.Float64()
}

func (r *BuiltInRand) Of(n int) int {
	return rand.Intn(n)
}

func NewBuiltInRand() Random {
	return &BuiltInRand{}
}

func (r *BuiltInRand) OfRange(min, max int) int {
	return rand.Intn(max-min) + min
}

func (r *BuiltInRand) OfIntRange(rng intRange) int {
	return rand.Intn(rng.Max-rng.Min) + rng.Min
}

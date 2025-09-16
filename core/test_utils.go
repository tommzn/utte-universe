package core

import (
	"fmt"

	"github.com/tommzn/go-config"
)

type mockRand struct {
	seekVal float64
	ofVal   int
}

func (r *mockRand) Seek() float64    { return r.seekVal }
func (r *mockRand) Of(n int) int     { return r.ofVal }
func (r *mockRand) Float64() float64 { return r.seekVal }
func (r *mockRand) Intn(n int) int   { return r.ofVal % n }

func (r *mockRand) OfRange(min, max int) int {
	if max <= min {
		return min
	}
	return min + (r.ofVal % (max - min))
}

func (r *mockRand) OfIntRange(rng intRange) int {
	if rng.Max <= rng.Min {
		return rng.Min
	}
	return rng.Min + (r.ofVal % (rng.Max - rng.Min))
}

func asInt(f float64) int {
	return int(f)
}

type mockLog struct {
	infos  []string
	errors []string
	debugs []string
}

func (l *mockLog) Info(format string, args ...interface{}) {
	l.infos = append(l.infos, fmt.Sprintf(format, args...))
}

func (l *mockLog) Error(format string, args ...interface{}) {
	l.errors = append(l.errors, fmt.Sprintf(format, args...))
}

func (l *mockLog) Debug(format string, args ...interface{}) {
	l.debugs = append(l.debugs, fmt.Sprintf(format, args...))
}

func loadConfigForTest(fileName *string) config.Config {

	configFile := "fixtures/testconfig.yml"
	if fileName != nil {
		configFile = *fileName
	}
	configLoader := config.NewFileConfigSource(&configFile)
	config, _ := configLoader.Load()
	return config
}

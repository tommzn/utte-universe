package core

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type BuiltInRandSuite struct {
	suite.Suite
	r *BuiltInRand
}

func TestBuiltInRandSuite(t *testing.T) {
	suite.Run(t, new(BuiltInRandSuite))
}

func (s *BuiltInRandSuite) SetupTest() {
	s.r = &BuiltInRand{}
}

func (s *BuiltInRandSuite) TestSeek() {
	val := s.r.Seek()
	s.True(val >= 0.0 && val < 1.0)
}

func (s *BuiltInRandSuite) TestOf() {
	val := s.r.Of(10)
	s.True(val >= 0 && val < 10)
}

func (s *BuiltInRandSuite) TestOfRange() {
	val := s.r.OfRange(2, 6)
	s.True(val >= 2 && val < 6)
}

func (s *BuiltInRandSuite) TestOfIntRange() {
	rng := intRange{Min: 10, Max: 15}
	val := s.r.OfIntRange(rng)
	s.True(val >= 10 && val < 15)
}

func (s *BuiltInRandSuite) TestNewBuiltInRand() {
	r := NewBuiltInRand()
	s.NotNil(r)
	_, ok := r.(*BuiltInRand)
	s.True(ok)
}

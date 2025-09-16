package core

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type BuildingSuite struct {
	suite.Suite
	planet   *Planet
	building *Building
	log      *mockLog
}

func TestBuildingSuite(t *testing.T) {
	suite.Run(t, new(BuildingSuite))
}

func (s *BuildingSuite) SetupTest() {
	s.planet = &Planet{
		Resources: map[ResourceType]int{
			Iron: 100,
			Food: 50,
			Fuel: 30,
		},
	}
	s.building = &Building{
		Type:      Mine,
		Level:     1,
		BuildCost: map[ResourceType]int{Iron: 10, Food: 5},
	}
	s.log = &mockLog{}
}

func (s *BuildingSuite) TestUpgradeSuccess() {
	ok := s.building.Upgrade(s.planet, s.log)
	s.True(ok)
	s.Equal(2, s.building.Level)
	s.Equal(100-20, s.planet.Resources[Iron]) // 10 * (1+1)
	s.Equal(50-10, s.planet.Resources[Food])  // 5 * (1+1)
}

func (s *BuildingSuite) TestUpgradeInsufficientResources() {
	s.planet.Resources[Iron] = 5
	ok := s.building.Upgrade(s.planet, s.log)
	s.False(ok)
	s.Equal(1, s.building.Level)
	s.Equal(5, s.planet.Resources[Iron])
}

func (s *BuildingSuite) TestUpgradeMultipleLevels() {
	s.building.Upgrade(s.planet, s.log) // Level 2
	s.building.Upgrade(s.planet, s.log) // Level 3
	s.Equal(3, s.building.Level)
	// Iron cost: 10*2 + 10*3 = 20 + 30 = 50
	// Food cost: 5*2 + 5*3 = 10 + 15 = 25
	s.Equal(100-20-30, s.planet.Resources[Iron])
	s.Equal(50-10-15, s.planet.Resources[Food])
}

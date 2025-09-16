package core

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type PlanetSuite struct {
	suite.Suite
	planet *Planet
	log    *mockLog
}

func TestPlanetSuite(t *testing.T) {
	suite.Run(t, new(PlanetSuite))
}

func (s *PlanetSuite) SetupTest() {
	s.planet = &Planet{
		Type: TerraLike,
		Resources: map[ResourceType]int{
			Iron: 100,
			Food: 50,
			Fuel: 30,
		},
		Buildings: []*Building{},
	}
	s.log = &mockLog{}
}

func (s *PlanetSuite) TestCanBuildFarmOnTerraLike() {
	farm := &Building{
		Type:      Farm,
		BuildCost: map[ResourceType]int{Iron: 10, Food: 5},
	}
	s.True(s.planet.CanBuild(farm, s.log))
}

func (s *PlanetSuite) TestCannotBuildFarmOnGasGiant() {
	s.planet.Type = GasGiant
	farm := &Building{
		Type:      Farm,
		BuildCost: map[ResourceType]int{Iron: 10, Food: 5},
	}
	s.False(s.planet.CanBuild(farm, s.log))
}

func (s *PlanetSuite) TestCanBuildMineOnTerraLike() {
	mine := &Building{
		Type:      Mine,
		BuildCost: map[ResourceType]int{Iron: 10},
	}
	s.True(s.planet.CanBuild(mine, s.log))
}

func (s *PlanetSuite) TestCannotBuildMineOnGasGiant() {
	s.planet.Type = GasGiant
	mine := &Building{
		Type:      Mine,
		BuildCost: map[ResourceType]int{Iron: 10},
	}
	s.False(s.planet.CanBuild(mine, s.log))
}

func (s *PlanetSuite) TestCanBuildRefineryAnywhere() {
	s.planet.Type = GasGiant
	refinery := &Building{
		Type:      Refinery,
		BuildCost: map[ResourceType]int{Iron: 10},
	}
	s.True(s.planet.CanBuild(refinery, s.log))
}

func (s *PlanetSuite) TestCanBuildCityAnywhere() {
	s.planet.Type = Icy
	city := &Building{
		Type:      City,
		BuildCost: map[ResourceType]int{Iron: 10},
	}
	s.True(s.planet.CanBuild(city, s.log))
}

func (s *PlanetSuite) TestCannotBuildIfInsufficientResources() {
	farm := &Building{
		Type:      Farm,
		BuildCost: map[ResourceType]int{Iron: 200, Food: 5},
	}
	s.False(s.planet.CanBuild(farm, s.log))
}

func (s *PlanetSuite) TestBuildSuccess() {
	mine := &Building{
		Type:      Mine,
		BuildCost: map[ResourceType]int{Iron: 10},
	}
	ok := s.planet.Build(mine, s.log)
	s.True(ok)
	s.Equal(90, s.planet.Resources[Iron])
	s.Contains(s.planet.Buildings, mine)
}

func (s *PlanetSuite) TestBuildFailure() {
	farm := &Building{
		Type:      Farm,
		BuildCost: map[ResourceType]int{Iron: 200, Food: 5},
	}
	ok := s.planet.Build(farm, s.log)
	s.False(ok)
	s.NotContains(s.planet.Buildings, farm)
}

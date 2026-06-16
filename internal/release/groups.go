package release

import (
	"fmt"
	"time"
)

type BetaGroup struct {
	Name        string    `json:"name"`
	TesterCount int       `json:"tester_count"`
	BuildCount  int       `json:"build_count"`
	CreatedAt   time.Time `json:"created_at"`
	PublicLink  string    `json:"public_link,omitempty"`
}

type GroupsResult struct {
	Groups []BetaGroup `json:"groups"`
	Total  int         `json:"total"`
}

func (s *Session) ListGroups() (*GroupsResult, error) {
	s.logInfo("listing TestFlight groups")

	groups := []BetaGroup{
		{
			Name:        "Internal Testers",
			TesterCount: 5,
			BuildCount:  3,
			CreatedAt:   time.Now().AddDate(0, -1, 0),
		},
		{
			Name:        "External Testers",
			TesterCount: 20,
			BuildCount:  2,
			CreatedAt:   time.Now().AddDate(0, -2, 0),
		},
	}

	result := &GroupsResult{
		Groups: groups,
		Total:  len(groups),
	}

	s.logInfo("groups listed", "count", result.Total)
	return result, nil
}

func (s *Session) InspectGroup(name string) (*BetaGroup, error) {
	s.logInfo("inspecting group", "name", name)

	groups, err := s.ListGroups()
	if err != nil {
		return nil, err
	}

	for _, g := range groups.Groups {
		if g.Name == name {
			return &g, nil
		}
	}

	return nil, fmt.Errorf("group not found: %s", name)
}

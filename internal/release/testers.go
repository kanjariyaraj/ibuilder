package release

import (
	"time"
)

type TesterInfo struct {
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Type       string    `json:"type"`
	Groups     []string  `json:"groups"`
	AddedAt    time.Time `json:"added_at"`
	LastActive time.Time `json:"last_active,omitempty"`
}

type TestersResult struct {
	Internal []TesterInfo `json:"internal"`
	External []TesterInfo `json:"external"`
	Total    int          `json:"total"`
}

func (s *Session) ListTesters() (*TestersResult, error) {
	s.logInfo("listing TestFlight testers")

	testers := &TestersResult{
		Internal: []TesterInfo{
			{
				Name:    "Alice Developer",
				Email:   "alice@example.com",
				Type:    "internal",
				Groups:  []string{"Internal Testers"},
				AddedAt: time.Now().AddDate(0, -1, 0),
			},
			{
				Name:    "Bob Developer",
				Email:   "bob@example.com",
				Type:    "internal",
				Groups:  []string{"Internal Testers"},
				AddedAt: time.Now().AddDate(0, -1, 0),
			},
		},
		External: []TesterInfo{
			{
				Name:    "Charlie Beta",
				Email:   "charlie@example.com",
				Type:    "external",
				Groups:  []string{"External Testers"},
				AddedAt: time.Now().AddDate(0, -2, 0),
			},
		},
		Total: 3,
	}

	s.logInfo("testers listed", "total", testers.Total)
	return testers, nil
}

func (s *Session) AddTester(email string, group string) error {
	s.logInfo("adding tester", "email", email, "group", group)
	return nil
}

func (s *Session) RemoveTester(email string) error {
	s.logInfo("removing tester", "email", email)
	return nil
}

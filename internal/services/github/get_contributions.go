package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gabsdev-go/internal/domain"
)

// contributionWindowWeeks condenses the heatmap to roughly 4-5 months
// instead of GitHub's full 52-week year, so it fits the site's fixed
// content width without horizontal scroll.
const contributionWindowWeeks = 20

const contributionsQuery = `query($login: String!, $from: DateTime!, $to: DateTime!) {
  user(login: $login) {
    contributionsCollection(from: $from, to: $to) {
      contributionCalendar {
        totalContributions
        weeks {
          contributionDays {
            date
            contributionCount
          }
        }
      }
    }
  }
}`

type contributionsResponse struct {
	Data struct {
		User struct {
			ContributionsCollection struct {
				ContributionCalendar struct {
					TotalContributions int `json:"totalContributions"`
					Weeks              []struct {
						ContributionDays []struct {
							Date              string `json:"date"`
							ContributionCount int    `json:"contributionCount"`
						} `json:"contributionDays"`
					} `json:"weeks"`
				} `json:"contributionCalendar"`
			} `json:"contributionsCollection"`
		} `json:"user"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// fetchContributions requires GITHUB_TOKEN (scope read:user) — the
// contribution calendar has no unauthenticated REST equivalent.
func fetchContributions() ([]domain.ContributionWeek, int, error) {
	tok := token()
	if tok == "" {
		return nil, 0, fmt.Errorf("github: GITHUB_TOKEN not set")
	}

	to := time.Now().UTC()
	from := to.AddDate(0, 0, -contributionWindowWeeks*7)

	body, err := json.Marshal(map[string]any{
		"query": contributionsQuery,
		"variables": map[string]string{
			"login": username(),
			"from":  from.Format(time.RFC3339),
			"to":    to.Format(time.RFC3339),
		},
	})
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequest(http.MethodPost, graphqlURL, bytes.NewReader(body))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("github graphql: unexpected status %d", resp.StatusCode)
	}

	var parsed contributionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, 0, err
	}
	if len(parsed.Errors) > 0 {
		return nil, 0, fmt.Errorf("github graphql: %s", parsed.Errors[0].Message)
	}

	cal := parsed.Data.User.ContributionsCollection.ContributionCalendar

	maxCount := 0
	for _, w := range cal.Weeks {
		for _, d := range w.ContributionDays {
			if d.ContributionCount > maxCount {
				maxCount = d.ContributionCount
			}
		}
	}

	weeks := make([]domain.ContributionWeek, 0, len(cal.Weeks))
	for _, w := range cal.Weeks {
		week := domain.ContributionWeek{Days: make([]domain.ContributionDay, 0, len(w.ContributionDays))}
		for _, d := range w.ContributionDays {
			week.Days = append(week.Days, domain.ContributionDay{
				Date:  d.Date,
				Count: d.ContributionCount,
				Level: contributionLevel(d.ContributionCount, maxCount),
			})
		}
		weeks = append(weeks, week)
	}

	return weeks, cal.TotalContributions, nil
}

func contributionLevel(count, max int) int {
	if count == 0 {
		return 0
	}
	if max <= 0 {
		return 1
	}
	ratio := float64(count) / float64(max)
	switch {
	case ratio >= 0.75:
		return 4
	case ratio >= 0.5:
		return 3
	case ratio >= 0.25:
		return 2
	default:
		return 1
	}
}

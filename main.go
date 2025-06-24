package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const pullRequestAPI string = "https://api.github.com/repos/%s/%s/pulls?state=closed&per_page=100&page=%d"

// PullRequests holds information on a GitHub pull request.
type PullRequest struct {
	Number    int       `json:"number"`
	CreatedAt time.Time `json:"created_at"`
	MergedAt  time.Time `json:"merged_at"`
}

// fetchMergedPRs fetches all merged PRs for a specific repo from GitHub.
func fetchMergedPRs(owner, repo, token string, page int) ([]PullRequest, error) {
	url := fmt.Sprintf(pullRequestAPI, owner, repo, page)

	client := &http.Client{}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "token "+token)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var allPRs []PullRequest
	json.Unmarshal(body, &allPRs)

	var mergedPRs []PullRequest
	for _, pr := range allPRs {
		// TODO: input dates, can probably put this on the API request for fewer results and pagination
		rangeStart := time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC)
		rangeEnd := time.Date(2025, 6, 30, 23, 59, 59, 0, time.UTC)

		if !pr.MergedAt.IsZero() && pr.MergedAt.After(rangeStart) && pr.MergedAt.Before(rangeEnd) {
			mergedPRs = append(mergedPRs, pr)
		}
	}

	return mergedPRs, nil
}

// getWeekdaysBetween gets the number of weekdays between two dates.
func getWeekdaysBetween(start, end time.Time) int {
	days := 0

	for day := start; day.Before(end); day = day.AddDate(0, 0, 1) {
		if day.Weekday() != time.Saturday && day.Weekday() != time.Sunday {
			days++
		}
	}

	return days
}

func main() {
	// TODO: get this from CLI input.
	owner := "CompanyCam"
	repo := "Company-Cam-API"
	token := os.Getenv("GITHUB_TOKEN")
	page := 1
	totalDays := 0
	count := 0

	for {
		prs, err := fetchMergedPRs(owner, repo, token, page)

		if err != nil || len(prs) == 0 {
			break
		}

		for _, pr := range prs {
			days := getWeekdaysBetween(pr.CreatedAt, pr.MergedAt)
			fmt.Printf("PR #%d: %d business days\n", pr.Number, days)

			totalDays += days
			count++
		}

		page++
	}

	if count > 0 {
		fmt.Printf("Average business days cycle time: %.2f\n", float64(totalDays)/float64(count))
	} else {
		fmt.Printf("No merged PRs found.")
	}
}

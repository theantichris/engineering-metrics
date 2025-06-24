package main

import "time"

type PullRequest struct {
	Number    int       `json:"number"`
	CreatedAt time.Time `json:"created_at"`
	MergedAt  time.Time `json:"merged_at"`
}

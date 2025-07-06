// Copyright 2025 Jamf Software LLC.
// https://developer.jamf.com/platform-api/reference/getbenchmark
// https://developer.jamf.com/platform-api/reference/gettenantbenchmarks
// https://developer.jamf.com/platform-api/reference/postbenchmark
// https://developer.jamf.com/platform-api/reference/deletebenchmark

package client

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

// Source represents source information
type Source struct {
	Branch   string `json:"branch" validate:"required"`
	Revision string `json:"revision" validate:"required"`
}

// TargetV2 represents the target configuration
type TargetV2 struct {
	DeviceGroups []string `json:"deviceGroups" validate:"required"`
}

// BenchmarkRequestV2 represents the request body for creating/updating benchmarks
type BenchmarkRequestV2 struct {
	Title            string        `json:"title" validate:"required,max=100"`
	Description      string        `json:"description,omitempty" validate:"max=1000"`
	SourceBaselineID string        `json:"sourceBaselineId" validate:"required"`
	Sources          []Source      `json:"sources" validate:"required"`
	Rules            []RuleRequest `json:"rules" validate:"required"`
	Target           TargetV2      `json:"target" validate:"required"`
	EnforcementMode  string        `json:"enforcementMode" validate:"required,oneof=MONITOR MONITOR_AND_ENFORCE"`
}

// BenchmarkResponseV2 represents the response for benchmark operations
type BenchmarkResponseV2 struct {
	BenchmarkID     string     `json:"benchmarkId"`
	TenantID        string     `json:"tenantId"`
	Title           string     `json:"title"`
	Description     string     `json:"description,omitempty"`
	Sources         []Source   `json:"sources"`
	Rules           []RuleInfo `json:"rules"`
	Target          TargetV2   `json:"target"`
	EnforcementMode string     `json:"enforcementMode"`
	Deleted         bool       `json:"deleted"`
	UpdateAvailable bool       `json:"updateAvailable"`
	LastUpdatedAt   time.Time  `json:"lastUpdatedAt"`
}

// BenchmarksResponseV2 represents the response for listing benchmarks
type BenchmarksResponseV2 struct {
	Benchmarks []BenchmarkV2 `json:"benchmarks"`
}

// BenchmarkV2 represents a benchmark in the list response
type BenchmarkV2 struct {
	ID              string   `json:"id"`
	Title           string   `json:"title"`
	Description     string   `json:"description,omitempty"`
	UpdateAvailable bool     `json:"updateAvailable"`
	Target          TargetV2 `json:"target"`
	SyncState       string   `json:"syncState"`
}

// CreateBenchmark creates a new benchmark
func (c *Client) CreateBenchmark(ctx context.Context, request *BenchmarkRequestV2) (*BenchmarkResponseV2, error) {
	resp, err := c.makeRequest(ctx, "POST", "/v2/benchmarks", request)
	if err != nil {
		return nil, fmt.Errorf("failed to create benchmark: %w", err)
	}

	var result BenchmarkResponseV2
	if err := c.handleAPIResponse(resp, 202, &result); err != nil {
		return nil, err
	}

	const pollInterval = 5
	const maxRetries = 24
	retries := 0
	for retries < maxRetries {
		time.Sleep(pollInterval * time.Second)
		benchmarks, err := c.GetBenchmarks(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to poll benchmarks after create: %w", err)
		}
		var found *BenchmarkV2
		for _, b := range benchmarks.Benchmarks {
			if b.ID == result.BenchmarkID {
				found = &b
				break
			}
		}
		if found == nil {
			return nil, fmt.Errorf("benchmark %s not found during polling after create", result.BenchmarkID)
		}
		switch found.SyncState {
		case "PENDING":
			retries++
			continue
		case "SYNCED":
			return &result, nil
		case "FAILED":
			_ = c.DeleteBenchmark(ctx, found.ID)
			return c.CreateBenchmark(ctx, request)
		default:
			return nil, fmt.Errorf("unexpected syncState after create: %s", found.SyncState)
		}
	}
	return nil, fmt.Errorf("timed out waiting for benchmark creation to complete after %d retries", maxRetries)
}

// GetBenchmarks retrieves all benchmarks for the tenant
func (c *Client) GetBenchmarks(ctx context.Context) (*BenchmarksResponseV2, error) {
	resp, err := c.makeRequest(ctx, "GET", "/v2/benchmarks", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get benchmarks: %w", err)
	}

	var result BenchmarksResponseV2
	if err := c.handleAPIResponse(resp, 200, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetBenchmarkByID retrieves a specific benchmark by ID
func (c *Client) GetBenchmarkByID(ctx context.Context, id string) (*BenchmarkResponseV2, error) {
	endpoint := fmt.Sprintf("/v2/benchmarks/%s", url.PathEscape(id))

	resp, err := c.makeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get benchmark %s: %w", id, err)
	}

	var result BenchmarkResponseV2
	if err := c.handleAPIResponse(resp, 200, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteBenchmark removes a benchmark by ID
func (c *Client) DeleteBenchmark(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("/v1/benchmarks/%s", url.PathEscape(id))

	resp, err := c.makeRequest(ctx, "DELETE", endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to delete benchmark %s: %w", id, err)
	}

	if err := c.handleAPIResponse(resp, 204, nil); err != nil {
		return err
	}

	const pollInterval = 5
	const maxRetries = 24
	retries := 0
	for retries < maxRetries {
		time.Sleep(pollInterval * time.Second)

		benchmarks, err := c.GetBenchmarks(ctx)
		if err != nil {
			return fmt.Errorf("failed to poll benchmarks after delete: %w", err)
		}

		found := false
		for _, b := range benchmarks.Benchmarks {
			if b.ID == id {
				found = true
				if b.SyncState == "DELETING" {
					_ = c.DeleteBenchmark(ctx, id)
					retries++
					break
				} else if b.SyncState == "DELETE_FAILED" {
					return fmt.Errorf("benchmark %s deletion failed: syncState=DELETE_FAILED", id)
				} else {
					return fmt.Errorf("benchmark %s still present after delete, syncState=%s", id, b.SyncState)
				}
			}
		}
		if !found {
			return nil
		}
	}
	return fmt.Errorf("timed out waiting for benchmark deletion to complete after %d retries", maxRetries)
}

// GetBenchmarkByTitle retrieves a specific benchmark by title
func (c *Client) GetBenchmarkByTitle(ctx context.Context, title string) (*BenchmarkResponseV2, error) {
	benchmarks, err := c.GetBenchmarks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get benchmarks list: %w", err)
	}

	for _, benchmark := range benchmarks.Benchmarks {
		if benchmark.Title == title {
			return c.GetBenchmarkByID(ctx, benchmark.ID)
		}
	}

	return nil, fmt.Errorf("benchmark with title '%s' not found", title)
}

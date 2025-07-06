// Copyright 2025 Jamf Software LLC.
// https://developer.jamf.com/platform-api/reference/getbaselines

package client

import (
	"context"
	"fmt"
)

// BaselinesResponse represents the response for baselines listing
type BaselinesResponse struct {
	Baselines []BaselineInfo `json:"baselines,omitempty"`
}

// BaselineInfo represents information about a baseline
type BaselineInfo struct {
	ID          string `json:"id,omitempty"`
	BaselineID  string `json:"baselineId,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Version     string `json:"version,omitempty"`
	Title       string `json:"title,omitempty"`
	RuleCount   int64  `json:"ruleCount,omitempty"`
}

// GetBaselines returns list of available mSCP baselines
func (c *Client) GetBaselines(ctx context.Context) (*BaselinesResponse, error) {
	resp, err := c.makeRequest(ctx, "GET", "/v1/baselines", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get baselines: %w", err)
	}

	var result BaselinesResponse
	if err := c.handleAPIResponse(resp, 200, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

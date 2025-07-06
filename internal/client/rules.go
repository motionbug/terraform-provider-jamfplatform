// Copyright 2025 Jamf Software LLC.
// https://developer.jamf.com/platform-api/reference/getrules

package client

import (
	"context"
	"fmt"
	"net/url"
)

// RuleRequest represents a rule in the request
type RuleRequest struct {
	ID      string      `json:"id" validate:"required"`
	Enabled bool        `json:"enabled" validate:"required"`
	ODV     *ODVRequest `json:"odv,omitempty"`
}

// ODVRequest represents an organization-defined value in requests
type ODVRequest struct {
	Value string `json:"value" validate:"required"`
}

// RuleInfo represents detailed rule information in responses
type RuleInfo struct {
	ID                 string                        `json:"id"`
	SectionName        string                        `json:"sectionName"`
	Enabled            bool                          `json:"enabled"`
	Title              string                        `json:"title"`
	References         []string                      `json:"references,omitempty"`
	Description        string                        `json:"description"`
	ODV                *OrganizationDefinedValue     `json:"odv,omitempty"`
	SupportedOS        []OSInfo                      `json:"supportedOs"`
	OSSpecificDefaults map[string]OSSpecificRuleInfo `json:"osSpecificDefaults"`
	RuleRelation       *RuleRelation                 `json:"ruleRelation,omitempty"`
}

// OrganizationDefinedValue represents ODV with full details
type OrganizationDefinedValue struct {
	Value       string                 `json:"value" validate:"required"`
	Hint        string                 `json:"hint,omitempty"`
	Placeholder string                 `json:"placeholder,omitempty"`
	Type        string                 `json:"type,omitempty"`
	Validation  *ValidationConstraints `json:"validation,omitempty"`
}

// ValidationConstraints represents validation rules for ODV
type ValidationConstraints struct {
	Min        *int     `json:"min,omitempty"`
	Max        *int     `json:"max,omitempty"`
	EnumValues []string `json:"enumValues,omitempty"`
	Regex      string   `json:"regex,omitempty"`
}

// OSInfo represents operating system information
type OSInfo struct {
	OSType         string `json:"osType"`
	OSVersion      int    `json:"osVersion"`
	ManagementType string `json:"managementType"`
}

// OSSpecificRuleInfo represents OS-specific rule details
type OSSpecificRuleInfo struct {
	Title       string             `json:"title"`
	Description string             `json:"description"`
	ODV         *ODVRecommendation `json:"odv,omitempty"`
}

// ODVRecommendation represents ODV recommendation
type ODVRecommendation struct {
	Value string `json:"value,omitempty"`
	Hint  string `json:"hint,omitempty"`
}

// RuleRelation represents rule dependencies
type RuleRelation struct {
	DependsOn []string `json:"dependsOn,omitempty"`
}

// SourcedRules represents rules with their sources
type SourcedRules struct {
	Sources []Source   `json:"sources"`
	Rules   []RuleInfo `json:"rules"`
}

// GetRules returns list of rules for given baseline
func (c *Client) GetRules(ctx context.Context, baselineID string) (*SourcedRules, error) {
	endpoint := fmt.Sprintf("/v1/rules?baselineId=%s", url.QueryEscape(baselineID))

	resp, err := c.makeRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get rules for baseline %s: %w", baselineID, err)
	}

	var result SourcedRules
	if err := c.handleAPIResponse(resp, 200, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

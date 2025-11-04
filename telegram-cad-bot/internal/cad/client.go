package cad

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"
)

// Client handles communication with the CAD API
type Client struct {
	BaseURL    string
	AgencyID   int
	HTTPClient *http.Client
	VerifySSL  bool
	Timeout    time.Duration
}

// NewClient creates a new CAD API client
func NewClient(baseURL string, agencyID int, verifySSL bool, timeout time.Duration) *Client {
	jar, _ := cookiejar.New(nil)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !verifySSL},
	}

	httpClient := &http.Client{
		Timeout:   timeout,
		Jar:       jar,
		Transport: transport,
	}

	return &Client{
		BaseURL:    baseURL,
		AgencyID:   agencyID,
		HTTPClient: httpClient,
		VerifySSL:  verifySSL,
		Timeout:    timeout,
	}
}

// GetActiveCalls fetches active CAD calls from the API
func (c *Client) GetActiveCalls(take int) (*CADResponse, error) {
	// Step 1: Visit main site to get cookies
	req, err := http.NewRequest("GET", c.BaseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create main page request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to visit main page: %w", err)
	}
	resp.Body.Close()

	// Step 2: Visit CADCalls page
	cadURL := fmt.Sprintf("%s/CADCalls", c.BaseURL)
	req, err = http.NewRequest("GET", cadURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CADCalls page request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err = c.HTTPClient.Do(req)
	if err != nil {
		// Try alternative URL
		cadURL = fmt.Sprintf("%s/Home/CADCalls", c.BaseURL)
		req, _ = http.NewRequest("GET", cadURL, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		resp, err = c.HTTPClient.Do(req)
		if err != nil {
			// Continue anyway, some sites might not require this
		}
	}
	if resp != nil {
		resp.Body.Close()
	}

	// Step 3: Make API request
	apiURL := fmt.Sprintf("%s/api/CADCalls/%d", c.BaseURL, c.AgencyID)

	payload := APIRequest{
		IncludeOpenCalls:   true,
		IncludeClosedCalls: false,
		IncludeCount:       true,
		PagingOptions: PagingOptions{
			SortOptions: []SortOption{
				{
					Name:          "StartTime",
					SortDirection: "Descending",
					Sequence:      1,
				},
			},
			Take: take,
			Skip: 0,
		},
		FilterOptions: FilterOptions{
			IntersectionSearch: true,
			SearchText:         "",
			Parameters:         []string{},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err = http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create API request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", c.BaseURL)
	req.Header.Set("Referer", cadURL)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err = c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make API request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)
		if len(bodyStr) > 500 {
			bodyStr = bodyStr[:500] + "..."
		}
		return nil, fmt.Errorf("API returned status %d.\nURL: %s\nResponse: %s",
			resp.StatusCode, apiURL, bodyStr)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var cadResp CADResponse
	if err := json.Unmarshal(body, &cadResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &cadResp, nil
}

package cad

import "time"

// CADCall represents a single CAD call from the API
type CADCall struct {
	IncidentID  string    `json:"IncidentId"`
	CallType    string    `json:"CallType"`
	Nature      string    `json:"Nature"`
	Address     string    `json:"Address"`
	StartTime   time.Time `json:"StartTime"`
	Agency      string    `json:"Agency"`
	HasLocation bool      `json:"HasLocation"`
	Latitude    float64   `json:"Latitude"`
	Longitude   float64   `json:"Longitude"`
}

// CADResponse represents the API response structure
type CADResponse struct {
	CADCalls []CADCall `json:"CADCalls"`
	Total    int       `json:"Total"`
}

// APIRequest represents the request payload for the CAD API
type APIRequest struct {
	IncludeOpenCalls   bool              `json:"IncludeOpenCalls"`
	IncludeClosedCalls bool              `json:"IncludeClosedCalls"`
	IncludeCount       bool              `json:"IncludeCount"`
	PagingOptions      PagingOptions     `json:"PagingOptions"`
	FilterOptions      FilterOptions     `json:"FilterOptionsParameters"`
}

// PagingOptions defines pagination and sorting
type PagingOptions struct {
	SortOptions []SortOption `json:"SortOptions"`
	Take        int          `json:"Take"`
	Skip        int          `json:"Skip"`
}

// SortOption defines how to sort results
type SortOption struct {
	Name          string `json:"Name"`
	SortDirection string `json:"SortDirection"`
	Sequence      int    `json:"Sequence"`
}

// FilterOptions defines search and filter parameters
type FilterOptions struct {
	IntersectionSearch bool        `json:"IntersectionSearch"`
	SearchText         string      `json:"SearchText"`
	Parameters         []string    `json:"Parameters"`
}

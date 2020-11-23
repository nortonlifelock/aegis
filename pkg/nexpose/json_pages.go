package nexpose

import (
	"context"
)

// Page identifies the information specific to paging in the nexpose Session
type Page struct {

	// The index (zero-based) of the current page returned.
	Number int64 `json:"number,omitempty"`

	// The maximum size of the page returned.
	Size int64 `json:"size,omitempty"`

	// The total number of pages available.
	TotalPages int64 `json:"totalPages,omitempty"`

	// The total number of resources available across all pages.
	TotalResources int64 `json:"totalResources,omitempty"`
}

// PageOfAsset is a json struct representation of paging in nexpose
type PageOfAsset struct {

	// Hypermedia links to corresponding or related resources.
	Links []Link `json:"links,omitempty"`

	// The details of pagination indicating which page was returned, and how the remaining pages can be retrieved.
	Page *Page `json:"page,omitempty"`

	// The page of resources returned.
	Res []Asset `json:"resources,omitempty"`
}

// Resources returns the resources of the page
func (pa *PageOfAsset) Resources(ctx context.Context, out chan<- interface{}) {

	// Push each fo the data onto the channel back to the requester
	for _, r := range pa.Res {
		select {
		case <-ctx.Done():
			return
		case out <- r:
		}
	}
}

// TotalPages returns the total number of pages for this set of queries
func (pa *PageOfAsset) TotalPages() int {
	var pages = 0
	if pa.Page != nil {
		pages = int(pa.Page.TotalPages)
	}

	return pages
}

// PageOfSolution is a json struct representation of paging in nexpose
type PageOfSolution struct {

	// Hypermedia links to corresponding or related resources.
	Links []Link `json:"links,omitempty"`

	// The details of pagination indicating which page was returned, and how the remaining pages can be retrieved.
	Page *Page `json:"page,omitempty"`

	// The page of resources returned.
	Res []Solution `json:"resources,omitempty"`
}

// Resources returns the resources of the page
func (pa *PageOfSolution) Resources(ctx context.Context, out chan<- interface{}) {

	// Push each fo the data onto the channel back to the requester
	for _, r := range pa.Res {
		select {
		case <-ctx.Done():
			return
		case out <- r:
		}
	}
}

// TotalPages returns the total number of pages for this set of queries
func (pa *PageOfSolution) TotalPages() int {
	var pages = 0
	if pa.Page != nil {
		pages = int(pa.Page.TotalPages)
	}

	return pages
}

// PageOfScan is a json struct representation of paging in nexpose
type PageOfScan struct {

	// Hypermedia links to corresponding or related resources.
	Links []Link `json:"links,omitempty"`

	// The details of pagination indicating which page was returned, and how the remaining pages can be retrieved.
	Page *Page `json:"page,omitempty"`

	// The page of resources returned.
	Res []Scan `json:"resources,omitempty"`
}

// Resources returns the resources of the page
func (pa *PageOfScan) Resources(ctx context.Context, out chan<- interface{}) {

	// Push each fo the data onto the channel back to the requester
	for _, r := range pa.Res {

		select {
		case <-ctx.Done():
			return
		case out <- r:
		}
	}
}

// TotalPages returns the total number of pages for this set of queries
func (pa *PageOfScan) TotalPages() int {
	var pages = 0
	if pa.Page != nil {
		pages = int(pa.Page.TotalPages)
	}

	return pages
}

// PageOfVulnerability is a json struct representation of paging in nexpose
type PageOfVulnerability struct {

	// Hypermedia links to corresponding or related resources.
	Links []Link `json:"links,omitempty"`

	// The details of pagination indicating which page was returned, and how the remaining pages can be retrieved.
	Page *Page `json:"page,omitempty"`

	// The page of resources returned.
	Res []Vulnerability `json:"resources,omitempty"`
}

// Resources returns the resources of the page
func (pa *PageOfVulnerability) Resources(ctx context.Context, out chan<- interface{}) {

	// Push each fo the data onto the channel back to the requester
	for _, r := range pa.Res {

		select {
		case <-ctx.Done():
			return
		case out <- r:
		}
	}
}

// TotalPages returns the total number of pages for this set of queries
func (pa *PageOfVulnerability) TotalPages() int {
	var pages = 0
	if pa.Page != nil {
		pages = int(pa.Page.TotalPages)
	}

	return pages
}

// PageOfVulnerabilityCheck is a json struct representation of paging in nexpose
type PageOfVulnerabilityCheck struct {

	// Hypermedia links to corresponding or related resources.
	Links []Link `json:"links,omitempty"`

	// The details of pagination indicating which page was returned, and how the remaining pages can be retrieved.
	Page *Page `json:"page,omitempty"`

	// The page of resources returned.
	Res []Check `json:"resources,omitempty"`
}

// Resources returns the resources of the page
func (pa *PageOfVulnerabilityCheck) Resources(ctx context.Context, out chan<- interface{}) {

	// Push each fo the data onto the channel back to the requester
	for _, r := range pa.Res {

		select {
		case <-ctx.Done():
			return
		case out <- r:
		}
	}
}

// TotalPages returns the total number of pages for this set of queries
func (pa *PageOfVulnerabilityCheck) TotalPages() int {
	var pages = 0
	if pa.Page != nil {
		pages = int(pa.Page.TotalPages)
	}

	return pages
}

// PageOfFinding is a json struct representation of paging in nexpose
type PageOfFinding struct {

	// Hypermedia links to corresponding or related resources.
	Links []Link `json:"links,omitempty"`

	// The details of pagination indicating which page was returned, and how the remaining pages can be retrieved.
	Page *Page `json:"page,omitempty"`

	// The page of resources returned.
	Res []Finding `json:"resources,omitempty"`
}

// Resources returns the resources of the page
func (pa *PageOfFinding) Resources(ctx context.Context, out chan<- interface{}) {

	// Push each fo the data onto the channel back to the requester
	for _, r := range pa.Res {

		select {
		case <-ctx.Done():
			return
		case out <- r:
		}
	}
}

// TotalPages returns the total number of pages for this set of queries
func (pa *PageOfFinding) TotalPages() int {
	var pages = 0
	if pa.Page != nil {
		pages = int(pa.Page.TotalPages)
	}

	return pages
}

// PageOfVulnerabilityReference is a json struct representation of paging in nexpose
type PageOfVulnerabilityReference struct {

	// Hypermedia links to corresponding or related resources.
	Links []Link `json:"links,omitempty"`

	// The details of pagination indicating which page was returned, and how the remaining pages can be retrieved.
	Page *Page `json:"page,omitempty"`

	// The page of resources returned.
	Res []VulnerabilityReference `json:"resources,omitempty"`
}

// Resources returns the resources of the page
func (pa *PageOfVulnerabilityReference) Resources(ctx context.Context, out chan<- interface{}) {

	// Push each fo the data onto the channel back to the requester
	for _, r := range pa.Res {

		select {
		case <-ctx.Done():
			return
		case out <- r:
		}
	}
}

// TotalPages returns the total number of pages for this set of queries
func (pa *PageOfVulnerabilityReference) TotalPages() int {
	var pages = 0
	if pa.Page != nil {
		pages = int(pa.Page.TotalPages)
	}

	return pages
}

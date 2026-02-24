package models

// AptRepositoryFile holds all the repos parsed from a single sources file.
type AptRepositoryFile struct {
	Digest       string          `json:"digest"`
	FileType     string          `json:"file-type"`
	Filename     string          `json:"filename"`
	Repositories []AptRepository `json:"repositories"`
}

// AptRepository is a single APT repository configuration entry.
type AptRepository struct {
	Components  []string `json:"Components,omitempty"`
	Enabled     bool     `json:"Enabled"`
	FileType    string   `json:"FileType,omitempty"`
	PackageType string   `json:"Package-Type,omitempty"`
	Suites      []string `json:"Suites,omitempty"`
	Types       []string `json:"Types,omitempty"`
	URIs        []string `json:"URIs,omitempty"`
}

// AptRepositoriesResponse is the full APT repo list response, grouped by source file.
type AptRepositoriesResponse struct {
	Digest string              `json:"digest"`
	Files  []AptRepositoryFile `json:"files"`
}

// AptRepositoryAddRequest is sent when enabling a standard APT repository by handle.
type AptRepositoryAddRequest struct {
	Handle string `json:"handle"`
	Digest string `json:"digest,omitempty"`
}

// AptRepositoryChangeRequest is sent to enable or disable an existing APT repository.
type AptRepositoryChangeRequest struct {
	Path    string `json:"path"`
	Index   int    `json:"index"`
	Enabled int    `json:"enabled"`
	Digest  string `json:"digest,omitempty"`
}

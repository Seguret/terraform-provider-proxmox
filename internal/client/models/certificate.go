package models

// NodeCertificate holds info about a TLS certificate installed on a node.
type NodeCertificate struct {
	Filename    string   `json:"filename,omitempty"`
	Fingerprint string   `json:"fingerprint,omitempty"`
	Issuer      string   `json:"issuer,omitempty"`
	NotAfter    int64    `json:"notafter,omitempty"`
	NotBefore   int64    `json:"notbefore,omitempty"`
	PEM         string   `json:"pem,omitempty"`
	PublicKeyBits int    `json:"public-key-bits,omitempty"`
	PublicKeyType string  `json:"public-key-type,omitempty"`
	SAN         []string `json:"san,omitempty"`
	Subject     string   `json:"subject,omitempty"`
}

// NodeCertificateUploadRequest is sent to upload a custom TLS cert to a node.
type NodeCertificateUploadRequest struct {
	Certificates string `json:"certificates"`
	Key          string `json:"key,omitempty"`
	Force        *int   `json:"force,omitempty"`
	Restart      *int   `json:"restart,omitempty"`
}

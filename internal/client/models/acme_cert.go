package models

// ACMECertOrderRequest is sent to order a new ACME certificate for a node.
type ACMECertOrderRequest struct {
	Force bool `json:"force,omitempty"`
}

// ACMERenewRequest is sent to renew an existing ACME certificate.
type ACMERenewRequest struct {
	Force bool `json:"force,omitempty"`
}

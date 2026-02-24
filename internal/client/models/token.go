package models

// UserToken is an API access token belonging to a user account.
type UserToken struct {
	TokenID string `json:"tokenid"`
	Comment string `json:"comment,omitempty"`
	Expire  int64  `json:"expire,omitempty"`
	Privsep *int   `json:"privsep,omitempty"`
	FullID  string `json:"full-tokenid,omitempty"`
	Value   string `json:"value,omitempty"` // only returned on creation
}

// UserTokenCreateRequest is sent when creating a new API token for a user.
type UserTokenCreateRequest struct {
	Comment string `json:"comment,omitempty"`
	Expire  int64  `json:"expire,omitempty"`
	Privsep *int   `json:"privsep,omitempty"`
}

// UserTokenCreateResponse is the create response - includes the token secret (only returned once).
type UserTokenCreateResponse struct {
	FullTokenID string    `json:"full-tokenid"`
	Info        UserToken `json:"info"`
	Value       string    `json:"value"`
}

// UserTokenUpdateRequest is sent to update token metadata (comment, expiry, privilege separation).
type UserTokenUpdateRequest struct {
	Comment string `json:"comment,omitempty"`
	Expire  *int64 `json:"expire,omitempty"`
	Privsep *int   `json:"privsep,omitempty"`
}

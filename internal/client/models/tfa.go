package models

// TFAEntry is a registered 2FA device or secret for a user.
type TFAEntry struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
	Created     int64  `json:"created,omitempty"`
	Enable      bool   `json:"enable,omitempty"`
}

// TFACreateRequest is sent when registering a new 2FA entry for a user.
type TFACreateRequest struct {
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
	TOTP        string `json:"totp,omitempty"`
	Value       string `json:"value,omitempty"`
	Challenge   string `json:"challenge,omitempty"`
	Password    string `json:"password,omitempty"`
}

// TFACreateResponse is returned after registering 2FA - may include QR code or recovery codes.
type TFACreateResponse struct {
	ID       string   `json:"id"`
	QRCode   string   `json:"qrcode,omitempty"`
	URL      string   `json:"url,omitempty"`
	Secret   string   `json:"secret,omitempty"`
	Recovery []string `json:"recovery,omitempty"`
}

// TFAUpdateRequest is sent to rename or enable/disable a 2FA entry.
type TFAUpdateRequest struct {
	Description string `json:"description,omitempty"`
	Enable      *bool  `json:"enable,omitempty"`
	Password    string `json:"password,omitempty"`
}

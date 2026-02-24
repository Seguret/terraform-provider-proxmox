package models

// User is a Proxmox VE user account.
type User struct {
	UserID    string `json:"userid"`
	Email     string `json:"email,omitempty"`
	Enable    *int   `json:"enable,omitempty"`
	Expire    int64  `json:"expire,omitempty"`
	FirstName string `json:"firstname,omitempty"`
	LastName  string `json:"lastname,omitempty"`
	Comment   string `json:"comment,omitempty"`
	Groups    string `json:"groups,omitempty"`
	Keys      string `json:"keys,omitempty"`
	Realm     string `json:"realm,omitempty"`
}

// UserCreateRequest is sent when creating a new user account.
type UserCreateRequest struct {
	UserID    string `json:"userid"`
	Password  string `json:"password,omitempty"`
	Email     string `json:"email,omitempty"`
	Enable    *int   `json:"enable,omitempty"`
	Expire    int64  `json:"expire,omitempty"`
	FirstName string `json:"firstname,omitempty"`
	LastName  string `json:"lastname,omitempty"`
	Comment   string `json:"comment,omitempty"`
	Groups    string `json:"groups,omitempty"`
	Keys      string `json:"keys,omitempty"`
}

// UserUpdateRequest is sent to modify user account details.
type UserUpdateRequest struct {
	Email     *string `json:"email,omitempty"`
	Enable    *int    `json:"enable,omitempty"`
	Expire    *int64  `json:"expire,omitempty"`
	FirstName *string `json:"firstname,omitempty"`
	LastName  *string `json:"lastname,omitempty"`
	Comment   *string `json:"comment,omitempty"`
	Groups    *string `json:"groups,omitempty"`
	Keys      *string `json:"keys,omitempty"`
}

// UserPermissions holds the permission map for a user across all API paths.
type UserPermissions struct {
	Permissions map[string]map[string]int `json:"permissions,omitempty"`
}

// OpenIDConfig is the OIDC discovery document info for a realm.
type OpenIDConfig struct {
	Issuer                string   `json:"issuer"`
	AuthorizationEndpoint string   `json:"authorization_endpoint"`
	TokenEndpoint         string   `json:"token_endpoint"`
	UserinfoEndpoint      string   `json:"userinfo_endpoint,omitempty"`
	JwksURI               string   `json:"jwks_uri,omitempty"`
	ScopesSupported       []string `json:"scopes_supported,omitempty"`
	ResponseTypesSupported []string `json:"response_types_supported,omitempty"`
}

// OpenIDAuthResponse contains the redirect URL for starting an OIDC auth flow.
type OpenIDAuthResponse struct {
	URL string `json:"url"`
}

// OpenIDLoginResponse is returned after a sucessful OIDC login, containing the ticket.
type OpenIDLoginResponse struct {
	Username  string `json:"username"`
	Ticket    string `json:"ticket"`
	CSRFToken string `json:"CSRFPreventionToken"`
}

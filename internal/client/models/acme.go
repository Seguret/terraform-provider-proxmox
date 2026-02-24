package models

// ACMEAccount holds the details of a registered ACME account.
type ACMEAccount struct {
	Name    string `json:"name,omitempty"`
	Contact string `json:"contact,omitempty"`
	// Directory is the ACME directory URL, e.g. Let's Encrypt production or staging endpoint.
	Directory  string `json:"directory,omitempty"`
	TosURL     string `json:"tos_url,omitempty"`
	AccountURL string `json:"account,omitempty"`
}

// ACMEAccountCreateRequest is sent when registering a new ACME account.
type ACMEAccountCreateRequest struct {
	Name      string `json:"name,omitempty"`
	Contact   string `json:"contact"`
	Directory string `json:"directory,omitempty"`
	TosURL    string `json:"tos_url,omitempty"`
}

// ACMEAccountUpdateRequest is sent to update an existing ACME account (currently only contact).
type ACMEAccountUpdateRequest struct {
	Contact string `json:"contact,omitempty"`
}

// ACMEAccountListEntry is a single entry in the ACME accounts list.
type ACMEAccountListEntry struct {
	Name string `json:"name"`
}

// ACMEPlugin holds the config for an ACME DNS validation plugin.
type ACMEPlugin struct {
	ID                string `json:"id"`
	Type              string `json:"type"`
	API               string `json:"api,omitempty"`
	Data              string `json:"data,omitempty"`
	Nodes             string `json:"nodes,omitempty"`
	ValidationDelay   int    `json:"validation-delay,omitempty"`
}

// ACMEPluginCreateRequest is sent when adding a new ACME DNS plugin.
type ACMEPluginCreateRequest struct {
	ID              string `json:"id"`
	Type            string `json:"type"`
	API             string `json:"api,omitempty"`
	Data            string `json:"data,omitempty"`
	Nodes           string `json:"nodes,omitempty"`
	ValidationDelay int    `json:"validation-delay,omitempty"`
}

// ACMEPluginUpdateRequest is sent when updating an existing ACME plugin config.
type ACMEPluginUpdateRequest struct {
	API             string `json:"api,omitempty"`
	Data            string `json:"data,omitempty"`
	Nodes           string `json:"nodes,omitempty"`
	ValidationDelay int    `json:"validation-delay,omitempty"`
}

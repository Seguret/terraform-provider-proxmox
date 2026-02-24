package models

// AuthRealm is an authentication realm (pam, pve, ldap, ad, openid, etc).
type AuthRealm struct {
	Realm   string `json:"realm"`
	Type    string `json:"type"` // pam, pve, ad, ldap, openid
	Comment string `json:"comment,omitempty"`
	Default *int   `json:"default,omitempty"`
	// LDAP and AD specific fields
	Server1       string `json:"server1,omitempty"`
	Server2       string `json:"server2,omitempty"`
	Port          *int   `json:"port,omitempty"`
	BaseDN        string `json:"base_dn,omitempty"`
	BindDN        string `json:"bind_dn,omitempty"`
	Password      string `json:"password,omitempty"`
	UserAttr      string `json:"user_attr,omitempty"`
	UserClasses   string `json:"user_classes,omitempty"`
	GroupDN       string `json:"groupdn,omitempty"`
	GroupFilter   string `json:"group_filter,omitempty"`
	GroupNameAttr string `json:"group_name_attr,omitempty"`
	GroupClasses  string `json:"group_classes,omitempty"`
	CASLPath      string `json:"capath,omitempty"`
	CAPath        string `json:"certpath,omitempty"`
	Cert          string `json:"cert,omitempty"`
	CertKey       string `json:"certkey,omitempty"`
	TFAType       string `json:"tfa,omitempty"`
	Secure        *int   `json:"secure,omitempty"`
	SSLVersion    string `json:"sslversion,omitempty"`
	Verify        *int   `json:"verify,omitempty"`
	Sync          string `json:"sync_defaults_options,omitempty"`
	// OpenID
	IssuerURL  string `json:"issuer-url,omitempty"`
	ClientID   string `json:"client-id,omitempty"`
	ClientKey  string `json:"client-key,omitempty"`
	Username   string `json:"username-claim,omitempty"`
	AutoCreate *int   `json:"autocreate,omitempty"`
}

// AuthRealmCreateRequest is sent when registering a new authentication realm.
type AuthRealmCreateRequest struct {
	Realm         string `json:"realm"`
	Type          string `json:"type"`
	Comment       string `json:"comment,omitempty"`
	Default       *int   `json:"default,omitempty"`
	Server1       string `json:"server1,omitempty"`
	Server2       string `json:"server2,omitempty"`
	Port          *int   `json:"port,omitempty"`
	BaseDN        string `json:"base_dn,omitempty"`
	BindDN        string `json:"bind_dn,omitempty"`
	Password      string `json:"password,omitempty"`
	UserAttr      string `json:"user_attr,omitempty"`
	TFAType       string `json:"tfa,omitempty"`
	Secure        *int   `json:"secure,omitempty"`
	SSLVersion    string `json:"sslversion,omitempty"`
	Verify        *int   `json:"verify,omitempty"`
	IssuerURL     string `json:"issuer-url,omitempty"`
	ClientID      string `json:"client-id,omitempty"`
	ClientKey     string `json:"client-key,omitempty"`
	UsernameClaim string `json:"username-claim,omitempty"`
	AutoCreate    *int   `json:"autocreate,omitempty"`
}

// AuthRealmUpdateRequest is sent to modify an existing auth realm's config.
type AuthRealmUpdateRequest struct {
	Comment       string `json:"comment,omitempty"`
	Default       *int   `json:"default,omitempty"`
	Server1       string `json:"server1,omitempty"`
	Server2       string `json:"server2,omitempty"`
	Port          *int   `json:"port,omitempty"`
	BaseDN        string `json:"base_dn,omitempty"`
	BindDN        string `json:"bind_dn,omitempty"`
	Password      string `json:"password,omitempty"`
	UserAttr      string `json:"user_attr,omitempty"`
	TFAType       string `json:"tfa,omitempty"`
	Secure        *int   `json:"secure,omitempty"`
	SSLVersion    string `json:"sslversion,omitempty"`
	Verify        *int   `json:"verify,omitempty"`
	IssuerURL     string `json:"issuer-url,omitempty"`
	ClientID      string `json:"client-id,omitempty"`
	ClientKey     string `json:"client-key,omitempty"`
	UsernameClaim string `json:"username-claim,omitempty"`
	AutoCreate    *int   `json:"autocreate,omitempty"`
}

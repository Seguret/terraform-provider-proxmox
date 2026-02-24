package models

// NotificationEndpointSendmail represents a sendmail notification endpoint.
type NotificationEndpointSendmail struct {
	Name        string   `json:"name"`
	Mailto      []string `json:"mailto,omitempty"`
	MailtoUser  []string `json:"mailto-user,omitempty"`
	FromAddress string   `json:"from-address,omitempty"`
	Author      string   `json:"author,omitempty"`
	Comment     string   `json:"comment,omitempty"`
	Disable     bool     `json:"disable,omitempty"`
}

// NotificationEndpointSendmailCreateRequest is the POST body for creating a sendmail endpoint.
type NotificationEndpointSendmailCreateRequest struct {
	Name        string   `json:"name"`
	Mailto      []string `json:"mailto,omitempty"`
	MailtoUser  []string `json:"mailto-user,omitempty"`
	FromAddress string   `json:"from-address,omitempty"`
	Author      string   `json:"author,omitempty"`
	Comment     string   `json:"comment,omitempty"`
	Disable     bool     `json:"disable,omitempty"`
}

// NotificationEndpointSendmailUpdateRequest is the PUT body for updating a sendmail endpoint.
type NotificationEndpointSendmailUpdateRequest struct {
	Mailto      []string `json:"mailto,omitempty"`
	MailtoUser  []string `json:"mailto-user,omitempty"`
	FromAddress string   `json:"from-address,omitempty"`
	Author      string   `json:"author,omitempty"`
	Comment     string   `json:"comment,omitempty"`
	Disable     *bool    `json:"disable,omitempty"`
}

// NotificationEndpointGotify represents a Gotify notification endpoint.
type NotificationEndpointGotify struct {
	Name    string `json:"name"`
	Server  string `json:"server"`
	Token   string `json:"token,omitempty"`
	Comment string `json:"comment,omitempty"`
	Disable bool   `json:"disable,omitempty"`
}

// NotificationEndpointGotifyCreateRequest is the POST body for creating a Gotify endpoint.
type NotificationEndpointGotifyCreateRequest struct {
	Name    string `json:"name"`
	Server  string `json:"server"`
	Token   string `json:"token,omitempty"`
	Comment string `json:"comment,omitempty"`
	Disable bool   `json:"disable,omitempty"`
}

// NotificationEndpointGotifyUpdateRequest is the PUT body for updating a Gotify endpoint.
type NotificationEndpointGotifyUpdateRequest struct {
	Server  string `json:"server,omitempty"`
	Token   string `json:"token,omitempty"`
	Comment string `json:"comment,omitempty"`
	Disable *bool  `json:"disable,omitempty"`
}

// NotificationEndpointSmtp represents an SMTP notification endpoint.
type NotificationEndpointSmtp struct {
	Name       string   `json:"name"`
	Server     string   `json:"server"`
	Port       int      `json:"port,omitempty"`
	Username   string   `json:"username,omitempty"`
	Mode       string   `json:"mode,omitempty"`
	Mailto     []string `json:"mailto,omitempty"`
	MailtoUser []string `json:"mailto-user,omitempty"`
	From       string   `json:"from"`
	Comment    string   `json:"comment,omitempty"`
	Disable    bool     `json:"disable,omitempty"`
}

// NotificationEndpointSmtpCreateRequest is the POST body for creating an SMTP endpoint.
type NotificationEndpointSmtpCreateRequest struct {
	Name       string   `json:"name"`
	Server     string   `json:"server"`
	Port       int      `json:"port,omitempty"`
	Username   string   `json:"username,omitempty"`
	Password   string   `json:"password,omitempty"`
	Mode       string   `json:"mode,omitempty"`
	Mailto     []string `json:"mailto,omitempty"`
	MailtoUser []string `json:"mailto-user,omitempty"`
	From       string   `json:"from"`
	Comment    string   `json:"comment,omitempty"`
	Disable    bool     `json:"disable,omitempty"`
}

// NotificationEndpointSmtpUpdateRequest is the PUT body for updating an SMTP endpoint.
type NotificationEndpointSmtpUpdateRequest struct {
	Server     string   `json:"server,omitempty"`
	Port       int      `json:"port,omitempty"`
	Username   string   `json:"username,omitempty"`
	Password   string   `json:"password,omitempty"`
	Mode       string   `json:"mode,omitempty"`
	Mailto     []string `json:"mailto,omitempty"`
	MailtoUser []string `json:"mailto-user,omitempty"`
	From       string   `json:"from,omitempty"`
	Comment    string   `json:"comment,omitempty"`
	Disable    *bool    `json:"disable,omitempty"`
}

// NotificationEndpointWebhook represents a webhook notification endpoint.
type NotificationEndpointWebhook struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Method  string `json:"method,omitempty"`
	Comment string `json:"comment,omitempty"`
	Disable bool   `json:"disable,omitempty"`
}

// NotificationEndpointWebhookCreateRequest is the POST body for creating a webhook endpoint.
type NotificationEndpointWebhookCreateRequest struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Method  string `json:"method,omitempty"`
	Comment string `json:"comment,omitempty"`
	Disable bool   `json:"disable,omitempty"`
}

// NotificationEndpointWebhookUpdateRequest is the PUT body for updating a webhook endpoint.
type NotificationEndpointWebhookUpdateRequest struct {
	URL     string `json:"url,omitempty"`
	Method  string `json:"method,omitempty"`
	Comment string `json:"comment,omitempty"`
	Disable *bool  `json:"disable,omitempty"`
}

// NotificationFilter represents a notification filter.
type NotificationFilter struct {
	Name        string   `json:"name"`
	MinSeverity string   `json:"minSeverity,omitempty"`
	MaxSeverity string   `json:"maxSeverity,omitempty"`
	Mode        string   `json:"mode,omitempty"`
	Rules       []string `json:"rules,omitempty"`
	Comment     string   `json:"comment,omitempty"`
	Disable     bool     `json:"disable,omitempty"`
}

// NotificationFilterCreateRequest is the POST body for creating a notification filter.
type NotificationFilterCreateRequest struct {
	Name        string   `json:"name"`
	MinSeverity string   `json:"minSeverity,omitempty"`
	MaxSeverity string   `json:"maxSeverity,omitempty"`
	Mode        string   `json:"mode,omitempty"`
	Rules       []string `json:"rules,omitempty"`
	Comment     string   `json:"comment,omitempty"`
	Disable     bool     `json:"disable,omitempty"`
}

// NotificationFilterUpdateRequest is the PUT body for updating a notification filter.
type NotificationFilterUpdateRequest struct {
	MinSeverity string   `json:"minSeverity,omitempty"`
	MaxSeverity string   `json:"maxSeverity,omitempty"`
	Mode        string   `json:"mode,omitempty"`
	Rules       []string `json:"rules,omitempty"`
	Comment     string   `json:"comment,omitempty"`
	Disable     *bool    `json:"disable,omitempty"`
}

// NotificationMatcher represents a notification matcher.
type NotificationMatcher struct {
	Name          string   `json:"name"`
	MatchSeverity []string `json:"match-severity,omitempty"`
	MatchCalendar []string `json:"match-calendar,omitempty"`
	MatchField    []string `json:"match-field,omitempty"`
	Target        []string `json:"target,omitempty"`
	Mode          string   `json:"mode,omitempty"`
	Comment       string   `json:"comment,omitempty"`
	Disable       bool     `json:"disable,omitempty"`
}

// NotificationMatcherCreateRequest is the POST body for creating a notification matcher.
type NotificationMatcherCreateRequest struct {
	Name          string   `json:"name"`
	MatchSeverity []string `json:"match-severity,omitempty"`
	MatchCalendar []string `json:"match-calendar,omitempty"`
	MatchField    []string `json:"match-field,omitempty"`
	Target        []string `json:"target,omitempty"`
	Mode          string   `json:"mode,omitempty"`
	Comment       string   `json:"comment,omitempty"`
	Disable       bool     `json:"disable,omitempty"`
}

// NotificationMatcherUpdateRequest is the PUT body for updating a notification matcher.
type NotificationMatcherUpdateRequest struct {
	MatchSeverity []string `json:"match-severity,omitempty"`
	MatchCalendar []string `json:"match-calendar,omitempty"`
	MatchField    []string `json:"match-field,omitempty"`
	Target        []string `json:"target,omitempty"`
	Mode          string   `json:"mode,omitempty"`
	Comment       string   `json:"comment,omitempty"`
	Disable       *bool    `json:"disable,omitempty"`
}

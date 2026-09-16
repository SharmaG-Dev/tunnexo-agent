package agent

type Config struct {
	Environment string
	ServerURL   string

	Token string

	LocalTarget string

	AgentIP string

	AgentName string

	TunnelName string
}

type AuthenticationResponse struct {
	OK        bool   `json:"ok"`
	AgentID   int    `json:"agentId"`
	AgentName string `json:"agentName"`
}

type TunnelCreateRequest struct {
	Name string `json:"name"`
}

type TunnelCreateResponse struct {
	OK        bool   `json:"ok"`
	TunnelID  int    `json:"tunnelId"`
	Subdomain string `json:"subdomain"`
	URL       string `json:"url"`
	Error     string `json:"error,omitempty"`
}

type HTTPRequest struct {
	RequestID string
	Subdomain string
	Method    string
	Path      string
	Headers   map[string]string
	Body      string
}

type HTTPResponse struct {
	RequestID string `json:"requestId"`

	StatusCode int `json:"statusCode"`

	Headers map[string]string `json:"headers"`

	Body string `json:"body"`

	IsBase64Encoded bool `json:"isBase64Encoded"`
}

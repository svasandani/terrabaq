package db

// Client - application requesting user data
type Client struct {
	Name        string `json:"name,omitempty"`
	ID          string `json:"id,omitempty"`
	Secret      string `json:"secret,omitempty"`
	RedirectURI string `json:"redirect_uri"`
}

// User - export user struct for http
type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Roles    []Role `json:"roles,omitempty"`
}

// Role - export generic Role struct
type Role struct {
	Type       string `json:"type"`
	ResourceID string `json:"resource_id"`
}

// ClientAccessRequest - struct for clients requesting user data
type ClientAccessRequest struct {
	GrantType string `json:"grant_type"`
	AuthCode  string `json:"auth_code"`
	Client    Client `json:"client"`
}

// ClientAccessResponse - struct for responding to client access request
type ClientAccessResponse struct {
	User User `json:"user"`
}

// EnqueueRequest - struct for getting a request to enqueue query
type EnqueueRequest struct {
	SessionToken string `json:"session_token"`
}

// EnqueueResponse - struct for responding to enqueue query
type EnqueueResponse struct {
	User User `json:"user"`
}

// UpdateRequest - struct for getting a request to update user roles
type UpdateRequest struct {
	OldUser User `json:"old_user"`
	NewUser User `json:"new_user"`
}

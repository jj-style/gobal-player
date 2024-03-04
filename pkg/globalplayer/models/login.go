package models

import "time"

type LoginRequest struct {
	LoginId  string `json:"loginID"`
	Password string `json:"password"`
	ApiKey   string `json:"ApiKey"`
}

type LoginResponse struct {
	CallID                     string      `json:"callId"`
	ErrorCode                  int         `json:"errorCode"`
	APIVersion                 int         `json:"apiVersion"`
	StatusCode                 int         `json:"statusCode"`
	StatusReason               string      `json:"statusReason"`
	Time                       time.Time   `json:"time"`
	RegisteredTimestamp        int         `json:"registeredTimestamp"`
	UID                        string      `json:"UID"`
	UIDSignature               string      `json:"UIDSignature"`
	SignatureTimestamp         string      `json:"signatureTimestamp"`
	Created                    time.Time   `json:"created"`
	CreatedTimestamp           int         `json:"createdTimestamp"`
	IsActive                   bool        `json:"isActive"`
	IsRegistered               bool        `json:"isRegistered"`
	IsVerified                 bool        `json:"isVerified"`
	LastLogin                  time.Time   `json:"lastLogin"`
	LastLoginTimestamp         int         `json:"lastLoginTimestamp"`
	LastUpdated                time.Time   `json:"lastUpdated"`
	LastUpdatedTimestamp       int64       `json:"lastUpdatedTimestamp"`
	LoginProvider              string      `json:"loginProvider"`
	OldestDataUpdated          time.Time   `json:"oldestDataUpdated"`
	OldestDataUpdatedTimestamp int64       `json:"oldestDataUpdatedTimestamp"`
	Profile                    Profile     `json:"profile"`
	Registered                 time.Time   `json:"registered"`
	SocialProviders            string      `json:"socialProviders"`
	NewUser                    bool        `json:"newUser"`
	SessionInfo                SessionInfo `json:"sessionInfo"`
}

type Profile struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Age        int    `json:"age"`
	BirthDay   int    `json:"birthDay"`
	BirthMonth int    `json:"birthMonth"`
	BirthYear  int    `json:"birthYear"`
	Country    string `json:"country"`
	Email      string `json:"email"`
	Gender     string `json:"gender"`
	Zip        string `json:"zip"`
}

type SessionInfo struct {
	CookieName  string `json:"cookieName"`
	CookieValue string `json:"cookieValue"`
}

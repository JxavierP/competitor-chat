package models

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

func (r Role) isValid() bool {
	switch r {
	case RoleUser, RoleAssistant:
		return true
	default:
		return false
	}
}

type Message struct {
	Role    Role `json:"role"`
	Content string `json:"content"`
}

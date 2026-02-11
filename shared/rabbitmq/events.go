package rabbitmq

const (
	IdentityUserRegisteredEvent = "identity.user.registered"
)

type IdentityUserRegisteredEventPayload struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

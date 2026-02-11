package rabbitmq

// --------------------
// Identity Events
// --------------------

const (
	IdentityUserRegisteredEvent    = "identity.user.registered"
	IdentityUserEmailVerifiedEvent = "identity.user.email_verified"
	IdentityUserLoggedInEvent      = "identity.user.logged_in"
)

// --------------------
// Organization Events
// --------------------

const (
	OrganizationCreatedEvent       = "organization.organization.created"
	OrganizationMemberAddedEvent   = "organization.member.added"
	OrganizationDeactivatedEvent   = "organization.organization.deactivated"
	OrganizationMemberRemovedEvent = "organization.member.removed"
	OrganizationRoleUpdatedEvent   = "organization.member.role_updated"
)

// --------------------
// Identity Payloads
// --------------------

type IdentityUserRegisteredEventPayload struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

type IdentityUserEmailVerifiedEventPayload struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

type IdentityUserLoggedInEventPayload struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

// --------------------
// Organization Payloads
// --------------------

type OrganizationCreatedEventPayload struct {
	OrganizationID string `json:"organization_id"`
	Name           string `json:"name"`
	OwnerUserID    string `json:"owner_user_id"`
}

type OrganizationMemberAddedEventPayload struct {
	OrganizationID string `json:"organization_id"`
	UserID         string `json:"user_id"`
	Role           string `json:"role"`
}

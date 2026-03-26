package rabbitmq

import "time"

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
// Tracking Events
// --------------------

const (
	TrackingDriverLocationUpdatedEvent = "tracking.driver.location_updated"
)

// --------------------
// Identity Payloads
// --------------------

type IdentityUserRegisteredEventPayload struct {
	UserID string `json:"user_id"`
}

type IdentityUserEmailVerifiedEventPayload struct {
	UserID string `json:"user_id"`
}

type IdentityUserLoggedInEventPayload struct {
	UserID string `json:"user_id"`
	OrgID  string `json:"org_id"`
}

// --------------------
// Organization Payloads
// --------------------

type OrganizationCreatedEventPayload struct {
	OrganizationID string `json:"organization_id"`
	OwnerUserID    string `json:"owner_user_id"`
}

type OrganizationMemberAddedEventPayload struct {
	OrganizationID string `json:"organization_id"`
	UserID         string `json:"user_id"`
	Role           string `json:"role"`
}

// --------------------
// Tracking Payloads
// --------------------

type TrackingDriverLocationUpdatedEventPayload struct {
	TripID     string    `json:"trip_id"`
	DriverID   string    `json:"driver_id"`
	VehicleID  string    `json:"vehicle_id"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	RecordedAt time.Time `json:"recorded_at"`
	Sequence   int       `json:"sequence"`
}

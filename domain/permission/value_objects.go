package permission

import (
	"fmt"
	"strings"
)

// Permission represents a specific access right.
// Format: <scope>:<serviceName>:<action>
type Permission struct {
	Scope       string
	ServiceName string
	Action      string
}

// New creates a new Permission from individual components.
func New(scope, serviceName, action string) Permission {
	return Permission{
		Scope:       scope,
		ServiceName: serviceName,
		Action:      action,
	}
}

// Parse converts a permission string into a Permission value object.
func Parse(p string) (Permission, error) {
	parts := strings.Split(p, ":")
	if len(parts) != 3 {
		return Permission{}, fmt.Errorf("invalid permission format: %s", p)
	}
	return Permission{
		Scope:       parts[0],
		ServiceName: parts[1],
		Action:      parts[2],
	}, nil
}

// String returns the string representation of the permission.
func (p Permission) String() string {
	return fmt.Sprintf("%s:%s:%s", p.Scope, p.ServiceName, p.Action)
}

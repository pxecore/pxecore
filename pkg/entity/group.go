package entity

// Group entity
type Group struct {
	ID                string
	Vars              map[string]string
	ParentID          string
	DefaultTemplateID string
	HostsIDs          []string
	GroupIDs          []string
}

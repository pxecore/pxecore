package entity

import "github.com/pxecore/pxecore/pkg/util"

// Group entity
type Group struct {
	ID         string
	Vars       map[string]string
	ParentID   string
	TemplateID string
	HostsIDs   []string
	GroupIDs   []string
}

// AddHost add host to the entity list.
func (g *Group) AddHost(h string) {
	g.HostsIDs = util.AddUniqueStringToSlice(g.HostsIDs, h)
}

// RemoveHost remove h
func (g *Group) RemoveHost(h string) {
	g.HostsIDs = util.RemoveStringFromSlice(g.HostsIDs, h)
}

// AddGroup add host to the entity list.
func (g *Group) AddGroup(h string) {
	g.GroupIDs = util.AddUniqueStringToSlice(g.GroupIDs, h)
}

// RemoveGroup remove h
func (g *Group) RemoveGroup(h string) {
	g.GroupIDs = util.RemoveStringFromSlice(g.GroupIDs, h)
}

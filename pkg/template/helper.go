package template

import (
	"github.com/pxecore/pxecore/pkg/errors"
	"github.com/pxecore/pxecore/pkg/repository"
)

// Helper assist template creation providing data and related functionality.
type Helper struct {
	HostID       string
	TemplateID   string
	Vars         map[string]string
	TemplateBody string
	repository   repository.Repository
}

// NewHelper construct new Helper
func NewHelper(repository repository.Repository, hostID string, templateID string) *Helper {
	h := new(Helper)
	h.repository = repository
	h.HostID = hostID
	h.TemplateID = templateID
	return h
}

// GetVar retrieves a particular key from the vars or the default string is returned.
func (h *Helper) GetVar(key string, def string) string {
	if h.Vars == nil {
		return def
	}
	s, ok := h.Vars[key]
	if ok {
		return s
	}
	return def
}

// Init Initializes the helper.
func (h *Helper) Init() error {
	return h.repository.Read(func(session repository.Session) error {
		host, err := session.Host().Get(h.HostID)
		if err != nil {
			return &errors.Error{Code: errors.ETemplateError, Msg: "template.Helper host not found.", Err: err}
		}
		h.Vars, h.TemplateID, err = recursiveGroupMerge(session, make([]string, 0), host.GroupID)
		if err != nil {
			return err
		}
		h.Vars = mergeMaps(h.Vars, host.Vars)
		if host.TemplateID != "" {
			h.TemplateID = host.TemplateID
		}
		template, err := session.Template().Get(h.TemplateID)
		if err != nil {
			return &errors.Error{Code: errors.ETemplateError, Msg: "template.Helper template not found.", Err: err}
		}
		h.TemplateBody = template.Template
		return nil
	})
}

// recursiveGroupVarsMerge retrieves the groups and merges the results.
func recursiveGroupMerge(session repository.Session, groups []string, groupID string) (map[string]string, string, error) {
	if groupID == "" {
		return make(map[string]string), "", nil
	}
	for _, g := range groups {
		if g == groupID {
			return nil, "", &errors.Error{Code: errors.ETemplateError, Msg: "template.Helper recursive group error."}
		}
	}
	group, err := session.Group().Get(groupID)
	if err != nil {
		return nil, "", &errors.Error{Code: errors.ETemplateError, Msg: "template.Helper group not found.", Err: err}
	}
	groups = append(groups, group.ID)
	m1, tID, err := recursiveGroupMerge(session, groups, group.ParentID)
	if err != nil {
		return nil, "", err
	}
	if group.TemplateID != "" {
		return mergeMaps(m1, group.Vars), group.TemplateID, nil
	}
	return mergeMaps(m1, group.Vars), tID, nil
}

// mergeMaps merges maps. m2 overrides m1.
func mergeMaps(m1 map[string]string, m2 map[string]string) map[string]string {
	for k, v := range m2 {
		m1[k] = v
	}
	return m1
}

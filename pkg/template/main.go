package template

import (
	rep "github.com/pxecore/pxecore/pkg/repository"
	"io"
	"text/template"
)

// Compile executes the template body and returns the compiled body.
func Compile(w io.Writer, repository rep.Repository, hostID string, templateID string) error {
	h := NewHelper(repository, hostID, templateID)
	if err := h.Init(); err != nil {
		return err
	}
	tmpl, err := template.New(h.TemplateID).Parse(h.TemplateBody)
	if err != nil {
		return err
	}
	err = tmpl.Execute(w, h)
	return err
}

// CompileWithHardwareAddr executes the template body and returns the compiled body.
func CompileWithHardwareAddr(w io.Writer, repository rep.Repository, HardwareAddr string, templateID string) error {
	h := ""
	if err := repository.Read(func(session rep.Session) error {
		host, err := session.Host().FindByHardwareAddr(HardwareAddr)
		if err != nil {
			return err
		}
		h = host.ID
		return nil
	}); err != nil {
		return err
	}
	return Compile(w, repository, h, templateID)
}

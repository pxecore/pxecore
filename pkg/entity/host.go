package entity

// Host entity
type Host struct {
	ID            string
	HardwareAddr  []string
	TrapMode      bool
	TrapTriggered bool
	Vars          map[string]string
	GroupID       string
	TemplateID    string
}

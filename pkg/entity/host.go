package entity

// Host entity
type Host struct {
	ID            string
	HardwareAddr  []string
	Group         string
	TrapMode      bool
	TrapTriggered bool
	vars          map[string]string
	template      string
}

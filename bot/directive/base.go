package directive

type BaseDirective struct {
	Type string `json:"type"`
}

func (d *BaseDirective) GenToken() string {
	return uuid.Must(uuid.NewV4()).String()
}

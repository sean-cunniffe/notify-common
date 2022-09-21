package component

type Component struct {
	Name        string `json:"name"`
	Href        string `json:"href"`
	Command     string `json:"command"`
	Description string `json:"description"`
}

package konkfile

type File struct {
	Commands map[string]Command `json:"commands"`
}

type Command struct {
	Run       string   `json:"run"`
	Needs     []string `json:"needs"`
	Exclusive bool     `json:"exclusive"`
}

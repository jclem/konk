package konkfile

type File struct {
	Commands map[string]Command `json:"commands" yaml:"commands"`
}

type Command struct {
	Run       string   `json:"run"       yaml:"run"`
	Needs     []string `json:"needs"     yaml:"needs"`
	Exclusive bool     `json:"exclusive" yaml:"exclusive"`
}

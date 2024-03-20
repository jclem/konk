package konkfile

type File struct {
	Commands map[string]Command `json:"commands" toml:"commands" yaml:"commands"`
}

type Command struct {
	Run       string   `json:"run"       toml:"run"       yaml:"run"`
	Needs     []string `json:"needs"     toml:"needs"     yaml:"needs"`
	Exclusive bool     `json:"exclusive" toml:"exclusive" yaml:"exclusive"`
}

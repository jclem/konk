package konkfile

type File struct {
	Commands map[string]Command `json:"commands"`
}

type Command struct {
	Run          string   `json:"run"`
	Dependencies []string `json:"dependencies"`
	Exclusive    bool     `json:"exclusive"`
}

package index

// DirNode an abstract representation of a path in a file system.
type DirNode struct {
	DirName      string   `json:"dir_name"`       // dirname from config's RootDir
	Locked       bool     `json:"locked"`         // indicates a dir been locked
	SubFileNames []string `json:"sub_file_names"` // sub files name
}

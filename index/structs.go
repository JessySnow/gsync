package index

// DirNode an abstract representation of a path in a file system.
type DirNode struct {
	DirName      string   // dirname from config's RootDir
	Locked       bool     // indicates a dir been locked
	SubFileNames []string // sub files name
}

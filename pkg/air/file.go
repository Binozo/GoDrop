package air

// File Represents an AirDrop file.
// Is either a file or a directory.
type File struct {
	FileName string

	// Defines the folder structure.
	// Could look like this: "./myDirectory"
	// or "./myDirectory/data"
	FileBomPath string

	// Defines if this element is a directory
	IsDirectory bool
}

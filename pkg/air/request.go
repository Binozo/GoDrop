package air

type Request struct {
	// Defines the sender's IPv6 address
	SenderIP string

	// Defines the sender application.
	// Example: com.apple.finder
	SenderApplication string

	// Defines the sender name.
	SenderComputerName string

	// Unique identifier of the sender.
	SenderID string

	// Defines the sender product model.
	SenderModelName string

	// Preview icon for the to-be-sent file.
	// If imagick is installed on the system the .png format will be sent.
	// Otherwise, empty byte slice
	FileIcon []byte

	// The collection of the files from this request.
	Files []File
}

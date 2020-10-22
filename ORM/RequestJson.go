package ORM

type MessageBlock struct {
	MessageType string   `json:"message_type"`
	Username    string   `json:"username"`
	Uuid        string   `json:"uuid"`
	Plain       string   `json:"plain"`
	SendTo      string   `json:"send_to"`
	FileInfo    FileInfo `json:"file_info"`
	Event       Event    `json:"event"`
	Session     string   `json:"session"`
}

type FileInfo struct {
	Name string `json:"name"`
	Size string `json:"size"`
	MD5  string `json:"MD5"`
}

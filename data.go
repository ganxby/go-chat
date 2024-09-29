package chat

type ServiceMessage struct {
	MessageType string `json:"messageType"`
	Username    string `json:"username"`
	Message     string `json:"message"`
	Color       string `json:"color"`
}

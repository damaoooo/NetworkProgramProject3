package ORM

type CommonResponse struct {
	Result string `json:"result"`
	Uuid   string `json:"uuid"`
}

type EventResponse struct {
	Uuid        string `json:"uuid"`
	MessageType string `json:"message_type"`
	Event       Event  `json:"event"`
}

type LoginResponse struct {
	Result  string `json:"result"`
	Uuid    string `json:"uuid"`
	Session string `json:"session"`
}

type WrongSession struct {
	Info string `json:"info"`
}

package resources

type Note struct {
	Data struct {
		Attributes struct {
			Content   string `json:"content"`
			CreatedAt string `json:"created_at"`
		} `json:"attributes"`
	} `json:"data"`
}

type NoteIDResponse struct {
	Data struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"data"`
}

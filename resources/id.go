package resources

type IdResponse struct {
	Data Key `json:"data"`
}

func NewIdResponse(key Key) IdResponse {
	return IdResponse{
		Data: key,
	}
}

type Key struct {
	ID   string       `json:"id"`
	Type ResourceType `json:"type"`
}

type ResourceType string

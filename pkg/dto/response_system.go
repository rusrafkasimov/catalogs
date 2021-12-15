package dto

type PingResponse struct {
	Payload struct {
		Ping string `json:"pong"`
	} `json:"payload"`
	Meta ResponseMeta `json:"meta"`
} // @Name PingResponse


package telegram

// Here we define all the types that the client will work with

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID      int    `json:"update_id"`
	Message string `json:"message"`
}

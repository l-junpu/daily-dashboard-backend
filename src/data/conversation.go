package data

type Conversation struct {
	Title     string    `json:"title" bson:"title"`
	Tags      []string  `json:"tags" bson:"tags"`
	Documents []string  `json:"documents" bson:"documents"`
	Messages  []Message `json:"messages" bson:"messages"`
}

type CreateConversationRequest struct {
	Username  string    `json:"username" bson:"username"`
	Title     string    `json:"title" bson:"title"`
	Tags      []string  `json:"tags" bson:"tags"`
	Documents []string  `json:"documents" bson:"documents"`
	Messages  []Message `json:"messages" bson:"messages"`
}

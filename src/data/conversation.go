package data

type Conversation struct {
	Title    string    `json:"title" bson:"title"`
	Tags     []string  `json:"tags" bson:"tags"`
	Messages []Message `json:"messages" bson:"messages"`
}

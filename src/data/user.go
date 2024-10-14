package data

type User struct {
	Username string `json:"username"`
}

type UserDetails struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type MongoUserDetails struct {
	Username string   `bson:"username"`
	Password string   `bson:"password"`
	Titles   []string `bson:"titles"`
}
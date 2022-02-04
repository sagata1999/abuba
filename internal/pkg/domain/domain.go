package domain

type Joke struct {
	Categories []string `json:"categories"`
	Id         string   `json:"id"`
	Value      string   `json:"value"`
}

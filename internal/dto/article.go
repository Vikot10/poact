package dto

type Article struct {
	ID         int
	Title      string
	Body       string
	Url        string
	Categories []Category
}

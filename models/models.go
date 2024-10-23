package models

type Jersey struct {
	Name string
	Url  string
}

type Album struct {
	Title   string
	Jerseys []Jersey
}

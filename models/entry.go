package models

type Entry struct {
	Username string
	Date     string
	Body     string
}

func NewEntry() *Entry {
	e := Entry{}
	return &e
}

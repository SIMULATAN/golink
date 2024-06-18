package model

type Link struct {
	Code   string `db:"code"`
	Target string `db:"target"`
}

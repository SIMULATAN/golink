package service

import (
	"errors"
	"golink/model"
)

var LinkTargetInvalidError = errors.New("link target invalid")
var LinkTargetInvalidServiceError = ServiceError{
	error:  LinkTargetInvalidError,
	Status: 401,
}

type LinkService interface {
	GetLink(code string) *model.Link
	CreateLink(target string, short *string) (*model.Link, error)
	Init()
}

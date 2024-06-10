package service

import (
	"errors"
	"golink/model"
	"golink/utils"
)

var LinkTargetInvalidError = errors.New("link target invalid")
var LinkTargetInvalidServiceError = ServiceError{
	error:  LinkTargetInvalidError,
	Status: 401,
}

type LinkService struct {
	links map[string]model.Link
}

func (s *LinkService) Init() {
	s.links = make(map[string]model.Link)
}

func (s *LinkService) CreateLink(target string, short *string) (*model.Link, error) {
	if target == "" {
		return nil, &LinkTargetInvalidServiceError
	}

	var shortLink string
	if short == nil || *short == "" {
		shortLink = utils.RandSeq(32)
	} else {
		shortLink = *short
	}
	link := model.Link{
		Code:   shortLink,
		Target: target,
	}
	s.links[shortLink] = link
	return &link, nil
}

func (s *LinkService) GetLink(code string) *model.Link {
	link, ok := s.links[code]
	if !ok {
		return nil
	}
	return &link
}

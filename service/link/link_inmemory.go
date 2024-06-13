package link

import (
	"golink/model"
	"golink/service"
	"golink/utils"
)

type InMemoryLinkService struct {
	links map[string]model.Link
}

func (s *InMemoryLinkService) Init() {
	s.links = make(map[string]model.Link)
}

func (s *InMemoryLinkService) CreateLink(target string, short *string) (*model.Link, error) {
	if target == "" {
		return nil, &service.LinkTargetInvalidServiceError
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

func (s *InMemoryLinkService) GetLink(code string) *model.Link {
	link, ok := s.links[code]
	if !ok {
		return nil
	}
	return &link
}

package search

import (
	"github.com/short-d/short/app/entity"
	"github.com/short-d/short/app/usecase/repository"
)

type Search struct {
	userUrlRepo repository.UserURLRelation
	urlRepo     repository.URL
}

// SearchForURLs fetches all URLs for a given user
func (s Search) SearchForURLs(user entity.User) ([]entity.URL, error) {
	aliases, err := s.userUrlRepo.FindAliasesByUser(user)
	if err != nil {
		return nil, err
	}

	urls, err := s.urlRepo.GetByAliases(aliases)
	if err != nil {
		return nil, err
	}

	return urls, nil
}

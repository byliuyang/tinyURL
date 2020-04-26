package repository

import (
	"errors"
	"time"

	"github.com/short-d/short/app/entity"
)

var _ URL = (*URLFake)(nil)

// URLFake accesses URL information in url table through SQL.
type URLFake struct {
	urls map[string]entity.URL
}

// IsAliasExist checks whether a given alias exist in url table.
func (u URLFake) IsAliasExist(alias string) (bool, error) {
	_, ok := u.urls[alias]
	return ok, nil
}

// Create inserts a new URL into url table.
func (u *URLFake) Create(url entity.URL) error {
	isExist, err := u.IsAliasExist(url.Alias)
	if err != nil {
		return err
	}
	if isExist {
		return errors.New("alias exists")
	}
	u.urls[url.Alias] = url
	return nil
}

// GetByAlias finds an URL in url table given alias.
func (u URLFake) GetByAlias(alias string) (entity.URL, error) {
	isExist, err := u.IsAliasExist(alias)
	if err != nil {
		return entity.URL{}, err
	}
	if !isExist {
		return entity.URL{}, errors.New("alias not found")
	}
	url := u.urls[alias]
	return url, nil
}

// GetByAliases finds all URL for a list of aliases
func (u URLFake) GetByAliases(aliases []string) ([]entity.URL, error) {
	if len(aliases) == 0 {
		return []entity.URL{}, nil
	}

	var urls []entity.URL
	for _, alias := range aliases {
		url, err := u.GetByAlias(alias)

		if err != nil {
			return urls, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}

// UpdateURL updates a URL that exists within the URL table and returns the newly updated URL if there is no error found while
// committing the update.
func (u URLFake) UpdateURL(
	key string,
	newAlias string,
	newOriginalURL string,
	expireAt *time.Time,
) (entity.URL, error) {
	prevURL, ok := u.urls[key]
	if !ok {
		return entity.URL{}, errors.New("URL to update not found")
	}

	now := time.Now().UTC()
	createdBy := prevURL.CreatedBy
	createdAt := prevURL.CreatedAt
	newURL := entity.URL{
		Alias:       newAlias,
		OriginalURL: newOriginalURL,
		ExpireAt:    expireAt,
		CreatedBy:   createdBy,
		CreatedAt:   createdAt,
		UpdatedAt:   &now,
	}
	return newURL, nil
}

// NewURLFake creates in memory URL repository
func NewURLFake(urls map[string]entity.URL) URLFake {
	return URLFake{
		urls: urls,
	}
}

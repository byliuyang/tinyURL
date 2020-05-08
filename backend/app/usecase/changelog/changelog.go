package changelog

import (
	"time"

	"github.com/short-d/app/fw/timer"
	"github.com/short-d/short/backend/app/entity"
	"github.com/short-d/short/backend/app/usecase/keygen"
	"github.com/short-d/short/backend/app/usecase/repository"
)

var _ ChangeLog = (*Persist)(nil)

// ChangeLog retrieves change log and create changes.
type ChangeLog interface {
	CreateChange(title string, summaryMarkdown *string) (entity.Change, error)
	GetChangeLog() ([]entity.Change, error)
	GetLastViewedAt(user entity.User) (*time.Time, error)
}

// Persist retrieves change log from and saves changes to persistent data store.
type Persist struct {
	keyGen            keygen.KeyGenerator
	timer             timer.Timer
	changeLogRepo     repository.ChangeLog
	userChangeLogRepo repository.UserChangeLog
}

// CreateChange creates a new change in the data store.
func (p Persist) CreateChange(title string, summaryMarkdown *string) (entity.Change, error) {
	now := p.timer.Now()
	key, err := p.keyGen.NewKey()
	if err != nil {
		return entity.Change{}, err
	}
	newChange := entity.Change{
		ID:              string(key),
		Title:           title,
		SummaryMarkdown: summaryMarkdown,
		ReleasedAt:      now,
	}
	return p.changeLogRepo.CreateChange(newChange)
}

// GetChangeLog retrieves full ChangeLog from persistent data store.
func (p Persist) GetChangeLog() ([]entity.Change, error) {
	return p.changeLogRepo.GetChangeLog()
}

// GetLastViewedAt retrieves the last time the user viewed the change log
func (p Persist) GetLastViewedAt(user entity.User) (*time.Time, error) {
	lastViewedAt, err := p.userChangeLogRepo.GetLastViewedAt(user)
	if err == nil {
		return &lastViewedAt, nil
	}

	switch err.(type) {
	case repository.ErrEntryNotFound:
		return nil, nil
	}

	return nil, err
}

// NewPersist creates Persist
func NewPersist(
	keyGen keygen.KeyGenerator,
	timer timer.Timer,
	changeLog repository.ChangeLog,
	userChangeLog repository.UserChangeLog,
) Persist {
	return Persist{
		keyGen:            keyGen,
		timer:             timer,
		changeLogRepo:     changeLog,
		userChangeLogRepo: userChangeLog,
	}
}

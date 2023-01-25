package repository

import (
	"context"

	"github.com/kumparan/go-utils"
	log "github.com/sirupsen/logrus"

	"github.com/atjhoendz/notpushcation-service/internal/model"
	"gorm.io/gorm"
)

type threadRepository struct {
	db *gorm.DB
}

//  NewThreadRepository :nodoc:
func NewThreadRepository(db *gorm.DB) model.ThreadRepository {
	return &threadRepository{db: db}
}

// Create :nodoc:
func (r threadRepository) Create(ctx context.Context, t *model.Thread) (*model.Thread, error) {
	logger := log.WithFields(log.Fields{
		"ctx":    utils.DumpIncomingContext(ctx),
		"thread": utils.Dump(t),
	})

	t.ID = utils.GenerateID()
	err := r.db.WithContext(ctx).Create(t).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return t, nil
}

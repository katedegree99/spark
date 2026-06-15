package repository

import (
	"context"

	domainrepo "github.com/katedegree/spark/api/internal/domain/repository"
	"github.com/katedegree/spark/api/internal/domain/model"
	"gorm.io/gorm"
)

type thingRepository struct {
	db *gorm.DB
}

func NewThingRepository(db *gorm.DB) domainrepo.ThingRepository {
	return &thingRepository{db: db}
}

func (r *thingRepository) Search(ctx context.Context, q string) ([]*domainrepo.ThingRecord, error) {
	var things []model.Thing
	query := r.db.WithContext(ctx).Limit(50)
	if q != "" {
		query = query.Where("name LIKE ?", q+"%")
	}
	if err := query.Find(&things).Error; err != nil {
		return nil, err
	}

	records := make([]*domainrepo.ThingRecord, 0, len(things))
	for _, t := range things {
		aliases, err := r.FindAliasesByThingID(ctx, t.ID)
		if err != nil {
			return nil, err
		}
		records = append(records, &domainrepo.ThingRecord{
			ID:        t.ID,
			Name:      t.Name,
			Aliases:   aliases,
			CreatedAt: t.CreatedAt,
		})
	}
	return records, nil
}

func (r *thingRepository) FindByID(ctx context.Context, id uint) (*domainrepo.ThingRecord, error) {
	var thing model.Thing
	if err := r.db.WithContext(ctx).First(&thing, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	aliases, err := r.FindAliasesByThingID(ctx, thing.ID)
	if err != nil {
		return nil, err
	}

	return &domainrepo.ThingRecord{
		ID:        thing.ID,
		Name:      thing.Name,
		Aliases:   aliases,
		CreatedAt: thing.CreatedAt,
	}, nil
}

func (r *thingRepository) FindByName(ctx context.Context, name string) (*domainrepo.ThingRecord, error) {
	var thing model.Thing
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&thing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	aliases, err := r.FindAliasesByThingID(ctx, thing.ID)
	if err != nil {
		return nil, err
	}

	return &domainrepo.ThingRecord{
		ID:        thing.ID,
		Name:      thing.Name,
		Aliases:   aliases,
		CreatedAt: thing.CreatedAt,
	}, nil
}

func (r *thingRepository) FindByIDs(ctx context.Context, ids []uint) ([]*domainrepo.ThingRecord, error) {
	if len(ids) == 0 {
		return []*domainrepo.ThingRecord{}, nil
	}

	var things []model.Thing
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&things).Error; err != nil {
		return nil, err
	}

	records := make([]*domainrepo.ThingRecord, 0, len(things))
	for _, t := range things {
		aliases, err := r.FindAliasesByThingID(ctx, t.ID)
		if err != nil {
			return nil, err
		}
		records = append(records, &domainrepo.ThingRecord{
			ID:        t.ID,
			Name:      t.Name,
			Aliases:   aliases,
			CreatedAt: t.CreatedAt,
		})
	}
	return records, nil
}

func (r *thingRepository) Create(ctx context.Context, name string) (*domainrepo.ThingRecord, error) {
	thing := model.Thing{Name: name}
	if err := r.db.WithContext(ctx).Create(&thing).Error; err != nil {
		return nil, err
	}
	return &domainrepo.ThingRecord{
		ID:        thing.ID,
		Name:      thing.Name,
		Aliases:   []string{},
		CreatedAt: thing.CreatedAt,
	}, nil
}

func (r *thingRepository) FindAliasesByThingID(ctx context.Context, thingID uint) ([]string, error) {
	var aliases []model.ThingAlias
	if err := r.db.WithContext(ctx).Where("thing_id = ?", thingID).Find(&aliases).Error; err != nil {
		return nil, err
	}

	result := make([]string, 0, len(aliases))
	for _, a := range aliases {
		result = append(result, a.Alias)
	}
	return result, nil
}

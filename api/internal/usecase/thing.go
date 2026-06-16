package usecase

import (
	"context"
	"errors"
	"log"

	domainllm "github.com/katedegree/spark/api/internal/domain/llm"
	"github.com/katedegree/spark/api/internal/domain/repository"
)

var ErrThingAlreadyExists = errors.New("thing already exists")

type ThingUsecase interface {
	Search(ctx context.Context, q string) ([]*repository.ThingRecord, error)
	Create(ctx context.Context, name string) (*repository.ThingRecord, error)
}

type thingUsecase struct {
	thingRepo repository.ThingRepository
	aliasSvc  domainllm.AliasService
}

func NewThingUsecase(thingRepo repository.ThingRepository, aliasSvc domainllm.AliasService) ThingUsecase {
	return &thingUsecase{thingRepo: thingRepo, aliasSvc: aliasSvc}
}

func (u *thingUsecase) Search(ctx context.Context, q string) ([]*repository.ThingRecord, error) {
	return u.thingRepo.Search(ctx, q)
}

func (u *thingUsecase) Create(ctx context.Context, name string) (*repository.ThingRecord, error) {
	existing, err := u.thingRepo.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrThingAlreadyExists
	}

	record, err := u.thingRepo.Create(ctx, name)
	if err != nil {
		return nil, err
	}

	// name 自体を alias として即時登録（検索を alias テーブルで完結させるため）
	if err := u.thingRepo.AddAlias(ctx, record.ID, name); err != nil {
		return nil, err
	}

	// AI による英語 alias を非同期生成
	go u.generateAlias(context.Background(), record.ID, name)

	return record, nil
}

func (u *thingUsecase) generateAlias(ctx context.Context, thingID uint, name string) {
	// guard: すでに alias が設定されていれば LLM を叩かない
	hasAlias, err := u.thingRepo.HasAlias(ctx, thingID)
	if err != nil {
		log.Printf("thing alias guard error: %v", err)
		return
	}
	if hasAlias {
		return
	}

	existing, err := u.thingRepo.FindAllWithAliases(ctx)
	if err != nil {
		log.Printf("thing alias fetch existing error: %v", err)
		return
	}

	// 新しく作った thing 自身は除外して渡す
	infos := make([]domainllm.ExistingThingInfo, 0, len(existing))
	for _, t := range existing {
		if t.ID == thingID {
			continue
		}
		infos = append(infos, domainllm.ExistingThingInfo{
			ID:      t.ID,
			Name:    t.Name,
			Aliases: t.Aliases,
		})
	}

	decision, err := u.aliasSvc.SuggestAlias(ctx, name, infos)
	if err != nil {
		log.Printf("thing alias llm error: %v", err)
		return
	}

	if err := u.applyDecision(ctx, thingID, decision); err != nil {
		log.Printf("thing alias apply error: %v", err)
	}
}

func (u *thingUsecase) applyDecision(ctx context.Context, thingID uint, d *domainllm.AliasDecision) error {
	switch d.Action {
	case "new_alias", "use_existing_alias":
		return u.thingRepo.AddAlias(ctx, thingID, d.Alias)
	case "rename_existing":
		if err := u.thingRepo.UpdateName(ctx, d.ThingIDToRename, d.NewName); err != nil {
			return err
		}
		return u.thingRepo.AddAlias(ctx, thingID, d.Alias)
	}
	return nil
}

package llm

import "context"

type AliasDecision struct {
	// new_alias | use_existing_alias | rename_existing
	Action          string
	Alias           string
	ExistingThingID uint
	ThingIDToRename uint
	NewName         string
}

type ExistingThingInfo struct {
	ID      uint
	Name    string
	Aliases []string
}

type AliasService interface {
	SuggestAlias(ctx context.Context, newThingName string, existing []ExistingThingInfo) (*AliasDecision, error)
}

package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	domainllm "github.com/katedegree/spark/api/internal/domain/llm"
)

const systemPrompt = `あなたは「やっていること」「やってみたいこと」タグの正規化担当です。
新しいタグが追加されたとき、以下のルールに従って判断してください。

ルール:
1. 既存のエイリアスと意味的に同じ概念の場合 → そのエイリアスを使用する
2. 既存タグの名前より新しいタグの名前の方がより適切・一般的な場合 → 既存タグの名前を変更する
3. 全く新しい概念の場合 → 英語小文字スネークケースの新しいエイリアスを生成する

必ず以下のいずれか1つのJSON形式のみで回答してください（説明は不要）:
{"action":"new_alias","alias":"programming"}
{"action":"use_existing_alias","alias":"programming","existing_thing_id":1}
{"action":"rename_existing","thing_id":1,"new_name":"プログラミング","alias":"programming"}`

type claudeService struct {
	client anthropic.Client
}

func NewClaudeAliasService() domainllm.AliasService {
	client := anthropic.NewClient(option.WithAPIKey(os.Getenv("ANTHROPIC_TOKEN")))
	return &claudeService{client: client}
}

func (s *claudeService) SuggestAlias(ctx context.Context, newThingName string, existing []domainllm.ExistingThingInfo) (*domainllm.AliasDecision, error) {
	userMsg := buildUserMessage(newThingName, existing)

	msg, err := s.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeHaiku4_5,
		MaxTokens: 200,
		System: []anthropic.TextBlockParam{
			{
				Text:         systemPrompt,
				CacheControl: anthropic.NewCacheControlEphemeralParam(),
			},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(userMsg)),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("claude api: %w", err)
	}

	raw := msg.Content[0].AsText().Text
	return parseDecision(raw)
}

func buildUserMessage(name string, existing []domainllm.ExistingThingInfo) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("新しいタグ: %q\n\n", name))

	if len(existing) == 0 {
		sb.WriteString("既存タグ: なし\n")
	} else {
		sb.WriteString("既存タグ:\n")
		for _, e := range existing {
			aliases := "なし"
			if len(e.Aliases) > 0 {
				aliases = strings.Join(e.Aliases, ", ")
			}
			sb.WriteString(fmt.Sprintf("- ID:%d, name:%q, aliases:[%s]\n", e.ID, e.Name, aliases))
		}
	}

	return sb.String()
}

type rawDecision struct {
	Action          string `json:"action"`
	Alias           string `json:"alias"`
	ExistingThingID uint   `json:"existing_thing_id"`
	ThingID         uint   `json:"thing_id"`
	NewName         string `json:"new_name"`
}

func parseDecision(raw string) (*domainllm.AliasDecision, error) {
	raw = strings.TrimSpace(raw)
	var d rawDecision
	if err := json.Unmarshal([]byte(raw), &d); err != nil {
		return nil, fmt.Errorf("parse decision %q: %w", raw, err)
	}

	switch d.Action {
	case "new_alias", "use_existing_alias", "rename_existing":
	default:
		return nil, fmt.Errorf("unknown action: %q", d.Action)
	}

	return &domainllm.AliasDecision{
		Action:          d.Action,
		Alias:           d.Alias,
		ExistingThingID: d.ExistingThingID,
		ThingIDToRename: d.ThingID,
		NewName:         d.NewName,
	}, nil
}

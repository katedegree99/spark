---
name: schema-driven
description: スキーマ駆動開発のワークフロー。API 変更時は必ずこの順序で作業する。常時適用
---

## 基本原則

**API の変更は必ず schema から始める。** コードを先に書いてはいけない。

## 作業順序（強制）

```
1. schema/openapi/openapi.yaml を編集
2. bun run lint（schema/ で実行）でバリデーション
3. make generate でコード再生成
   - web/src/lib/api/ → orval が SWR フック + 型定義を生成
   - api/pkg/generated/api.gen.go → oapi-codegen が ServerInterface を生成
4. api 側で StrictServerInterface を実装
5. web 側で生成されたフックを使って UI を実装
```

## コード生成の仕組み

### web（orval + SWR）
- 設定: `web/orval.config.ts`
- 出力: `web/src/lib/api/generated.ts`（SWR フック）、`web/src/lib/api/model/`（型定義）
- スキーマを変更したら `bun run generate`（web/ で実行）

### api（oapi-codegen）
- 設定: `api/oapi-codegen.yaml`
- 出力: `api/pkg/generated/api.gen.go`
- `StrictServerInterface` を `internal/adapter/handler/` で実装する
- スキーマを変更したら api/ で `oapi-codegen --config=oapi-codegen.yaml ../schema/openapi/openapi.yaml`

## 禁止事項

- `api/pkg/generated/api.gen.go` を手動編集してはいけない（`DO NOT EDIT` コメントあり）
- `web/src/lib/api/generated.ts` および `web/src/lib/api/model/` を手動編集してはいけない
- OpenAPI 3.1 の構文は使わない（oapi-codegen が 3.0.3 まで対応）

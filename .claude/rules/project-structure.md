---
name: project-structure
description: spark モノレポの構成・ポート割り当て・主要コマンド。常時適用
---

## ディレクトリ構成

```
spark/
├── web/          # Next.js 16 + TypeScript + Tailwind + Biome (bun)
├── api/          # Go 1.25 + Echo + GORM + oapi-codegen (clean arch)
├── schema/       # OpenAPI 3.0.3 仕様書 (Redocly)
├── Makefile      # 統合起動スクリプト
└── .worktrees/   # git worktree 作業ディレクトリ（.gitignore 済み）
```

## ポート割り当て（固定）

| ポート | サービス |
|---|---|
| 3000 | Next.js (web) |
| 3001 | Redocly mock サーバー（モック API） |
| 8080 | Go API サーバー（本番 API） |
| 8081 | Redocly ドキュメントプレビュー |

## 主要コマンド

| コマンド | 内容 |
|---|---|
| `make up` | 全サービス並列起動（down してから起動） |
| `make down` | 全サービス停止（ポートベース kill） |
| `make generate` | web（orval）と api（oapi-codegen）のコード一括再生成 |
| `make mock` | モック API サーバーのみ起動 |
| `make schema` | ドキュメントプレビューのみ起動 |
| `bun run lint` | schema バリデーション（schema/ ディレクトリで実行） |

# spark CLAUDE.md

spark は web / api / schema の 3 サービスから成るモノレポ。
スキーマ駆動開発を基本とし、OpenAPI スペックを唯一の真実とする。

## .claude/ ファイル索引

| ファイル | 役割 |
|---|---|
| `project-structure.md` | **プロジェクト構成・ポート割り当て・主要コマンド**。ディレクトリ構成、各サービスのポート（3000/3001/8080/8081）、make / bun コマンド一覧。常時適用 |
| `schema-driven.md` | **スキーマ駆動開発ワークフロー（強制）**。API 変更は必ず `schema/openapi/openapi.yaml` から始める原則、`make generate` の実行タイミング、orval / oapi-codegen の生成物の禁止事項。常時適用 |
| `architecture.md` | **api/ クリーンアーキテクチャ規約**。4 層構成（domain / usecase / infrastructure / adapter）の依存方向、dig による DI 規約、マイグレーション管理、ハンドラー実装規約。api/ を編集するとき適用 |
| `git-workflow.md` | **ブランチ戦略・worktree 運用・PR 規約**。feature ブランチ命名、`.worktrees/` の使い方、コミットメッセージ規約、PR タイトル規約。git 操作時に適用 |

## 技術スタック早見表

| サービス | 言語 / FW | パッケージマネージャ |
|---|---|---|
| web | Next.js 16 + TypeScript | bun |
| api | Go 1.25 + Echo v4 + GORM | go modules |
| schema | OpenAPI 3.0.3 + Redocly CLI | bun |

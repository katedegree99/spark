# .claude/ ファイル索引

### rules/

| ファイル | 役割 |
|---|---|
| `tech-stack.md` | **技術スタックと環境変数**。各サービスの言語・FW・主要ライブラリ・`.env` 変数一覧。常時適用 |
| `project-structure.md` | **プロジェクト構成・ポート割り当て・主要コマンド**。ディレクトリ構成、各サービスのポート（3000/3001/8080/8081）、make / bun コマンド一覧。常時適用 |
| `schema-driven.md` | **スキーマ駆動開発ワークフロー（強制）**。API 変更は必ず `schema/openapi/openapi.yaml` から始める原則、`make generate` の実行タイミング、生成物への直接編集禁止。常時適用 |
| `architecture.md` | **api/ クリーンアーキテクチャ規約**。4 層構成（domain / usecase / infrastructure / adapter）の依存方向、dig による DI 規約、マイグレーション管理、ハンドラー実装規約。`api/` を編集するとき適用 |
| `git-workflow.md` | **ブランチ戦略・worktree 運用・PR 規約**。feature ブランチ命名、`.worktrees/` の使い方、コミットメッセージ規約、PR タイトル規約。git 操作時に適用 |

### skills/

| ファイル | 役割 |
|---|---|
| `gh-login/SKILL.md` | **GitHub CLI 認証確認 skill**。`gh auth status` で認証状態を確認し、未ログインなら `! gh auth login` の実行を案内する。PR 作成など `gh` を使う操作の前、または `/gh-login` 明示起動 |
| `google-login/SKILL.md` | **Google OAuth 認証情報セットアップ skill**。Secret Manager への client_id / secret の格納状況を確認し、未設定なら gcloud コマンドで設定する手順を案内する。Google ログイン実装前、または `/google-login` 明示起動 |

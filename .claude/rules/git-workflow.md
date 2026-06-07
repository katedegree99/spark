---
name: git-workflow
description: ブランチ戦略・worktree 運用・PR 規約。git 操作時に適用
---

## ブランチ戦略

- `main` → 常時デプロイ可能な状態を保つ
- `feature/<scope>-<name>` → 機能実装（例: `feature/auth-register`）
- エンドポイント単位で 1 ブランチ / 1 PR を原則とする

## worktree 運用

並列実装が必要なときは `.worktrees/` 配下に worktree を作成する。

```bash
# 作成
git worktree add .worktrees/<name> -b feature/<name>

# 削除（ブランチ削除も含む）
git worktree remove .worktrees/<name>
git branch -d feature/<name>
```

- `.worktrees/` は `.gitignore` 済み
- worktree 内のファイルを編集してもメインの作業ディレクトリには影響しない
- 各 worktree は独立した `go build` / `bun install` が必要な場合がある

## コミット規約

```
<type>: <subject>

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>
```

| type | 用途 |
|---|---|
| `feat` | 新機能 |
| `fix` | バグ修正 |
| `refactor` | リファクタリング |
| `chore` | ビルド・設定変更 |
| `docs` | ドキュメント |

## PR 規約

- タイトル: `<type>: <HTTP メソッド> <パス>`（例: `feat: POST /auth/register`）
- 1 PR = 1 エンドポイントを基本とする
- `make generate` の差分はコミットに含める

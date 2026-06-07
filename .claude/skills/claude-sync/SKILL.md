---
name: claude-sync
description: .claude/ や CLAUDE.md の変更を claude-code ブランチに同期し、現在のブランチへローカルマージする skill
---

## 概要

`/claude-sync` で起動する。`.claude/` または `CLAUDE.md` に変更がある場合に、
`claude-code` ブランチの worktree を作成し、変更を適用・プッシュした後、
現在のブランチにローカルマージする。

## 起動条件

- `/claude-sync` 明示起動のみ（自動起動しない）

## 実行フロー

1. `git status` で変更ファイルを取得し、`.claude/` 配下または `CLAUDE.md` に該当するものを抽出する
2. 該当ファイルがなければ「対象ファイルなし」と伝えて終了する
3. 現在のブランチ名を記録する（`git branch --show-current`）
4. `claude-code` ブランチの worktree を `.worktrees/claude-code` に作成する
   - すでに `.worktrees/claude-code` が存在する場合は `git worktree remove .worktrees/claude-code --force` で削除してからやり直す
   - リモートに `claude-code` が存在するか確認: `git ls-remote --heads origin claude-code`
   - 存在する場合: `git worktree add .worktrees/claude-code claude-code`
   - 存在しない場合: `git worktree add .worktrees/claude-code -b claude-code`
5. リモートに `claude-code` が存在する場合のみ最新を取得する
   ```
   git -C .worktrees/claude-code pull origin claude-code --rebase
   ```
6. 対象ファイルをリポジトリルートから worktree へコピーする（削除ファイルは worktree 側でも `git rm`）
7. worktree 内でコミットする
   ```
   git -C .worktrees/claude-code add -A
   git -C .worktrees/claude-code commit -m "chore: <変更内容の日本語要約>

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>"
   ```
8. worktree をプッシュする
   ```
   git -C .worktrees/claude-code push origin claude-code
   ```
9. worktree を削除する
   ```
   git worktree remove .worktrees/claude-code
   ```
10. 現在のブランチに `claude-code` をマージする
    ```
    git merge claude-code --no-edit
    ```

## 注意事項

- `.env` や秘密情報を含むファイルは絶対にコピー・コミットしない
- force push は行わない
- マージ後にコンフリクトが発生した場合はユーザーに報告して手動解決を促す

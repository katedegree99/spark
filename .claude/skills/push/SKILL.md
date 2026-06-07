---
name: push
description: 変更差分を確認し、日本語のコミットメッセージを生成して git commit → git push まで行う skill
---

## 概要

`/push` で起動する。`git diff` と `git status` から変更内容を把握し、
日本語で要約したコミットメッセージを生成して、確認なしに即座に
`git commit` と `git push` を実行する。

## 起動条件

- **人間が `/push` と入力した場合にのみ実行する**
- Claude が `Skill` ツール経由で自律的に呼び出すことは禁止
- タスクの流れの中で「push しておこう」と判断しても、このスキルを自動呼び出ししてはいけない

## 実行フロー

1. `git status` で未追跡・変更済みファイルを確認する
2. `git diff` および `git diff --staged` で差分を確認する
3. 変更内容を日本語で要約し、以下の形式でコミットメッセージを生成する

```
<type>: <日本語の要約>

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>
```

type は以下から選ぶ：

| type | 用途 |
|---|---|
| `feat` | 新機能 |
| `fix` | バグ修正 |
| `refactor` | リファクタリング |
| `chore` | ビルド・設定変更 |
| `docs` | ドキュメント |

4. ステージされていないファイルはすべて `git add` してからコミットする
   - `.env` や秘密情報を含むファイルは除外する
5. 確認なしに即座に `git commit` を実行する
6. 続けて `git push origin <現在のブランチ>` を実行する

## 注意事項

- ユーザーへの確認は不要。メッセージを提示したらそのまま実行する
- `api/.env` など `.gitignore` 対象ファイルは絶対にステージしない
- `main` ブランチへの force push は行わない

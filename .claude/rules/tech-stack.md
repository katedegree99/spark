---
name: tech-stack
description: 各サービスの技術スタックと環境情報。常時適用
---

## サービス構成

| サービス | 言語 / FW | パッケージマネージャ | 役割 |
|---|---|---|---|
| web | Next.js 16 + TypeScript + Tailwind CSS | bun | フロントエンド |
| api | Go 1.25 + Echo v4 + GORM | go modules | バックエンド API |
| schema | OpenAPI 3.0.3 + Redocly CLI v2 | bun | API 仕様管理 |

## web の主要ライブラリ

- **orval** — OpenAPI → SWR フック + 型定義のコード生成
- **swr** — データフェッチ（orval SWR モードで生成されたフックを利用）
- **Biome** — Linter / Formatter（ESLint + Prettier の代替）

## api の主要ライブラリ

- **oapi-codegen v2** — OpenAPI → Echo `StrictServerInterface` + 型定義のコード生成
- **uber-go/dig** — 依存注入コンテナ
- **golang-jwt/jwt v5** — JWT 生成・検証
- **golang.org/x/crypto/bcrypt** — パスワードハッシュ

## 環境変数

### web（`.env`）

| 変数 | 説明 |
|---|---|
| `NEXT_PUBLIC_API_URL` | モック API URL（開発時は `http://localhost:3001`） |

### api（`.env` または OS 環境変数）

| 変数 | 説明 |
|---|---|
| `DB_USER` | MySQL ユーザー |
| `DB_PASSWORD` | MySQL パスワード |
| `DB_HOST` | MySQL ホスト |
| `DB_PORT` | MySQL ポート |
| `DB_NAME` | MySQL データベース名 |
| `JWT_SECRET` | JWT 署名シークレット |
| `GOOGLE_CLIENT_ID` | Google OAuth クライアント ID |

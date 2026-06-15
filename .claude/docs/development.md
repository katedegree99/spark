# 開発手順

## よく使うコマンド

### ルート (Makefile) — 全サービスを同時起動

```bash
make up       # web + schema preview + mock を並列起動
make down     # next dev / redocly preview-docs / prism mock を停止
make restart  # down → up
```

> `make` の schema/mock ターゲットは nvm の Node (`v24.16.0`) を PATH に通して実行する(`Makefile` の `NVM_NODE`)。

### web/ (Next.js)

```bash
cd web
bun run dev       # 開発サーバ (http://localhost:3000)
bun run build     # 本番ビルド
bun run generate  # OpenAPI から API クライアントを再生成 (orval)
bun run check     # Biome チェック (lint + format 検証)
bun run lint      # Biome lint --write (自動修正)
bun run format    # Biome format --write
```

### schema/ (Redocly)

```bash
cd schema
bun run preview   # API ドキュメントプレビュー (Redocly CLI v1, http://localhost:8081)
bun run mock      # モックサーバ (Stoplight Prism, http://localhost:3001)
bun run lint      # OpenAPI 仕様の lint (Redocly)
bun run bundle    # 仕様を1ファイルに bundle → dist/openapi.yaml
```

## ポート割り当て

| サービス | ポート |
|---|---|
| web (Next.js dev) | 3000 |
| schema mock (Prism) | 3001 |
| schema preview (Redocly) | 8081 |
| API サーバ (spec の `servers` 定義) | 8080 |

## ツールチェーンの規約

- **Bun を Node.js の代わりに使う**(両ワークスペース共通)。詳細な Bun API の使い分けは `schema/CLAUDE.md` を参照。
- **`web/` の lint/format は Biome**(ESLint/Prettier ではない)。インデントはタブ、JS の文字列はダブルクォート(`web/biome.json`)。
- **この Next.js は破壊的変更を含むバージョン**。`web/` でコードを書く前に `node_modules/next/dist/docs/` の該当ガイドを読むこと(詳細は `web/AGENTS.md`)。

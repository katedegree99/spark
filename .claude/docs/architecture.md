# アーキテクチャ

Spark は2つのワークスペースで構成される。**OpenAPI 仕様 (`schema/`) が API 契約の単一の真実源** で、フロントエンド (`web/`) の API クライアントはそこから自動生成される。

```
schema/openapi/openapi.yaml  ──(orval)──►  web/src/lib/api/{generated.ts, model/}
```

- **`schema/`** — OpenAPI 3.1 の API 仕様。ドキュメントプレビューは Redocly、モックサーバは Stoplight Prism、lint は Redocly が担う。ランタイムは Bun。
- **`web/`** — Next.js 16 / React 19 のフロントエンド。`bun` をパッケージマネージャに使う。

## 契約駆動の重要ポイント

- API のリクエスト/レスポンス型を変えたいときは、**まず `schema/openapi/openapi.yaml` を編集** し、その後 `web/` で `bun run generate` (orval) を実行してクライアントを再生成する。
- `web/src/lib/api/generated.ts` と `web/src/lib/api/model/` は**生成物。手で編集しない**(orval が `clean: true` で毎回上書きする)。
- 生成クライアントは orval の `swr` モードで、SWR フックとして出力される。
- 現状の API は `auth` タグの認証エンドポイント群のみ(register / login / otp verify / google / refresh / logout)。

# アーキテクチャ

Spark は2つのワークスペースで構成される。**OpenAPI 仕様 (`schema/`) が API 契約の単一の真実源** で、フロントエンド (`web/`) の API クライアントはそこから自動生成される。

```
schema/openapi/openapi.yaml  ──(orval)──►  web/src/lib/api/generated/{client.ts, model/}
```

- **`schema/`** — OpenAPI 3.1 の API 仕様。ドキュメントプレビューは Redocly、モックサーバは Stoplight Prism、lint は Redocly が担う。ランタイムは Bun。
- **`web/`** — Next.js 16 / React 19 のフロントエンド。`bun` をパッケージマネージャに使う。

## 契約駆動の重要ポイント

- API のリクエスト/レスポンス型を変えたいときは、**まず `schema/openapi/openapi.yaml` を編集** し、その後 `web/` で `bun run generate` (orval) を実行してクライアントを再生成する。
- `web/src/lib/api/generated/`(`client.ts` / `model/`)は**生成物。手で編集しない**(orval が `clean: true` で毎回上書きする)。手書きの `fetcher.ts` は `generated/` の外に置く(同居すると `clean` で消えるため)。
- 生成クライアントは orval の `swr` モード(`generated/client.ts`)。ただし**認証の変更系は SWR フックを使わず、Server Action から `apiFetch`(`lib/api/fetcher.ts`)で server-to-server に通す**(トークンを httpOnly Cookie に隔離するため)。生成 SWR フックは将来の認証後・参照系(GET)取得用に温存する。
- 現状の API は `auth` タグの認証エンドポイント群のみ(register / login / otp verify / google / refresh / logout)。
- **TODO(次PR): 自動リフレッシュ / ログアウトが未配線。** フロントは register / login / otp / google までを実装済み。`/auth/refresh`・`/auth/logout` は backend にあるが未使用。access(15分)失効後に refresh(7日)で更新する経路が無く、現状は失効すると `/login` に戻る。App Router は Server Component 描画中に Cookie を書けないため、**Middleware で先回りリフレッシュ**(access 無 & refresh 有 → `/auth/refresh` → response に Cookie セット)+ `logoutAction`(`clearAuthCookies`)を次PRで実装する。refresh Cookie の寿命は backend DB の有効期限(7日)に一致させること。

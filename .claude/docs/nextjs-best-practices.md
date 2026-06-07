# Next.js ベストプラクティス (App Router / v16)

> このプロジェクトの Next.js は **16.2.7**。バージョン固有の破壊的変更があるため、**コードを書く前に必ず `web/node_modules/next/dist/docs/` の該当ガイドを読む**(`web/AGENTS.md` の指示)。以下は同梱ドキュメントから本リポジトリ向けに要約したもので、出典パスを併記する。記憶や訓練データではなく、必ず原典で裏取りする。

## 大前提

- App Router (`web/src/app/`) を使用。**layout / page はデフォルトで Server Component**。
- データ取得クライアントは `schema/` の OpenAPI から **orval が SWR フックとして生成**(`web/src/lib/api/generated.ts`)。`useSWR` はクライアント実行なので、これを使う UI は Client Component になる。

## Server / Client Component の使い分け

出典: `01-getting-started/05-server-and-client-components.md`

| 使う場面 | コンポーネント |
|---|---|
| state / イベントハンドラ (`onClick`, `onChange`)、`useEffect`、ブラウザAPI (`localStorage`/`window`)、カスタムフック(= `useSWR`) | **Client** (`"use client"`) |
| DB/API からのデータ取得、秘密情報(APIキー/トークン)の保持、JS バンドル削減、FCP 改善 | **Server**(デフォルト) |

- `"use client"` は**境界**。付けたファイルの import とそれが直接 render する子はすべてクライアントバンドルに入る。**できるだけ末端の対話的コンポーネントだけに付け、境界をツリーの下に下げる**(バンドル肥大を防ぐ)。
- Server Component は Client Component に **props または `children` 経由で合成**できる。Client の中に Server を「子」として差し込むパターン(モーダル等)を活用する。
- **Context Provider は Client Component に切り出し**、`layout` で `{children}` だけを包む(`<html>` 全体を包まない=静的部分の最適化を妨げない)。

## データ取得(本プロジェクトの実態に即して)

出典: `02-guides/single-page-applications.md`(SWR 節), `01-getting-started/06-fetching-data.md`

- 生成された `useSWR` フックは Client Component で使う。SWR のポーリング/再検証/キャッシュはクライアント側でのみ動く。
- 可能なデータは Server Component で取得して**ウォーターフォールを避ける**。SWR 2.3 + React 19 では、Server Component で取得した値を `<SWRConfig fallback={...}>` に渡し(`await` しない)、子の `useSWR(key)` で受け取るハイブリッドが可能 — フック側のコードは変更不要。
- 並列取得でネットワーク往復を減らす。

## フォームと認証(auth 画面)

出典: `02-guides/authentication.md`, `02-guides/forms.md`, `01-getting-started/07-mutating-data.md`

- フォーム送信は **Server Actions + `useActionState`** が推奨。Server Action は常にサーバ実行なので認証ロジックの安全な置き場になる。
- **サーバ側でバリデーション**する(Zod 等のスキーマ)。失敗時は早期 return して外部APIへの無駄な呼び出しを防ぐ。
- **トークン(access/refresh JWT)は `httpOnly` Cookie に置き、クライアント JS に露出させない**のが原則。`localStorage` 保存は XSS リスク。
  - ⚠️ 設計上の論点: 現状は OpenAPI の `servers` が `http://localhost:8080`(外部 API)で、orval の SWR クライアントはブラウザから直接叩く形。トークンの保管/付与方式(Cookie + Next の Route Handler / Proxy 経由にするか、クライアント保持にするか)は**実装前に方針を確定する**。安易に localStorage に入れない。出典: `02-guides/backend-for-frontend.md`(Route Handler / Proxy)。

## Request-time API と動的レンダリング

出典: `02-guides/production-checklist.md`, `01-app/04-glossary.md`

- `cookies()`, `headers()`, および `params` / `searchParams` は **Promise**。`await` してから使う(v15+ の変更)。
- これらを使うとそのルートは **Dynamic Rendering にオプトイン**する。**Root Layout で使うとアプリ全体が動的化**するため意図的に。必要箇所は `<Suspense>` で包む。

## 秘密情報と環境変数

出典: `01-getting-started/05-server-and-client-components.md`(環境汚染防止節), `02-guides/environment-variables.md`

- クライアントに渡るのは **`NEXT_PUBLIC_` 接頭辞の変数のみ**。接頭辞なしは空文字に置換される。
- サーバ専用モジュールは **`server-only`** パッケージを import して、クライアントへの誤 import をビルド時エラーにする。

## ルーティング / UX / 最適化

出典: `02-guides/production-checklist.md`

- ナビゲーションは **`<Link>`**(自動 prefetch・クライアント遷移)。共有 UI は **layout** に置き部分レンダリングを活かす。
- エラーは `error.tsx` / `app/global-error.tsx`、404 は `not-found.tsx` で扱う。
- 画像は **`next/image`**、フォントは **`next/font`**(レイアウトシフト防止・外部リクエスト削減)。

# ディレクトリ設計 (`web/src/`)

`web/` フロントエンドの内部ディレクトリ構成と配置ルール。**ハイブリッド構成** — 汎用部品は layer 別、機能固有のものは `features/<機能名>/` に閉じ込める。

> 前提となる規約は別ドキュメントを正とする:
> - リポジトリ全体(`schema/` ↔ `web/`)とスキーマ駆動: `architecture.md`
> - `cn` / `zustand` / `motion` / shadcn 互換方針: `library-usage.md`
> - Server/Client 分離・認証・Server Actions: `nextjs-best-practices.md`

## 全体ツリー

```
web/src/
  app/                  # ルーティング (App Router)。下記「app/ の設計」参照
  components/
    ui/                 # 汎用 UI 部品 (shadcn 互換)。cn + cva で自作
    <共通複合部品>       # ヘッダー等、複数機能で使う非 UI プリミティブ
  features/
    auth/               # 認証機能一式(下記「features/ の設計」参照)
      components/        #   この機能専用のコンポーネント
      hooks/             #   この機能専用のフック(任意)
      actions.ts         #   Server Actions
      schema.ts          #   Zod 等のバリデーションスキーマ
      store.ts           #   この機能専用の zustand store(任意)
  hooks/                # 複数機能で共有するカスタムフック
  lib/
    api/                # ★ orval 生成物。手で編集しない (generated.ts / model/)
    <その他共有ライブラリ> # api クライアント設定・サーバ専用ヘルパ等
  stores/               # 横断的(グローバル)な zustand store
  utils/                # 純粋関数ユーティリティ (cn.ts など)
```

## 配置の判断基準 — features か共有か

**1つの機能だけで使うか、複数機能で使うか**で置き場所を決める。

| 対象 | 1機能だけで使う | 複数機能で共有する |
|---|---|---|
| コンポーネント | `features/<機能>/components/` | UI プリミティブ → `components/ui/`<br>複合部品 → `components/` |
| カスタムフック | `features/<機能>/hooks/` | `hooks/` |
| zustand store | `features/<機能>/store.ts` | `stores/` |
| 純粋関数 | `features/<機能>/` 内に同居 | `utils/` |

- **迷ったら features 側に置く**。共有が必要になった時点で `components/` `hooks/` `utils/` へ引き上げる(早すぎる共通化を避ける)。
- `features/` 間の相互 import は避ける。共有したくなったら共有層へ昇格させる合図。
- 1ファイルで足りる要素(`actions.ts` / `schema.ts` / `store.ts`)はフォルダにせず直置き。増えてきたら同名ディレクトリに展開する。

## `features/<機能>/` の設計

機能に固有のものを 1 ディレクトリに凝集させる。認証 (`auth`) を例にすると:

| ファイル/ディレクトリ | 役割 |
|---|---|
| `components/` | その機能の画面・フォーム部品(`LoginForm` 等) |
| `actions.ts` | Server Actions。**フォーム送信・認証ロジックの安全な置き場**(`nextjs-best-practices.md` 参照) |
| `schema.ts` | 入力バリデーションスキーマ。Server Action とクライアント双方で共有 |
| `hooks/` | その機能専用フック(任意) |
| `store.ts` | その機能のクライアント状態(任意・zustand) |

- API 呼び出しは `lib/api/`(orval 生成の SWR フック)を使う。**`features/` 内で API 型やフックを再定義しない**。
- 認証のトークン保管方式(httpOnly Cookie + Route Handler/Proxy か、クライアント保持か)は **`nextjs-best-practices.md` の論点が未確定**。確定するまで `localStorage` に入れない。

## `app/` の設計(ルーティング)

App Router。**Route Group** で「認証前 (`(auth)`)」と「認証後 (`(app)`)」を分け、それぞれに専用 layout を持たせる。`( )` 付きディレクトリは **URL に出ない**グルーピング用。

```
app/
  layout.tsx            # ルート layout(<html>/<body>、Provider は Client に切り出して children だけ包む)
  globals.css           # デザイントークン(CSS 変数)。shadcn テーマ方式に合わせる
  page.tsx              # ルート("/")
  (auth)/               # 認証前の画面群(URL には出ない)
    layout.tsx          #   中央寄せ等、認証画面共通レイアウト
    login/page.tsx      #   /login
    register/page.tsx   #   /register
    otp/page.tsx        #   /otp(OTP 検証)
  (app)/                # 認証後の画面群(将来)
    layout.tsx          #   サイドバー等のアプリ共通シェル
  api/                  # Route Handler(必要になれば。BFF/Proxy・Cookie 付与の置き場)
```

> 上記ルート名(`login` / `register` / `otp`)は現状の OpenAPI の `auth` タグ(register / login / otp verify / google / refresh / logout)に対応する**想定**。実際のパスは実装時に確定する。

設計上の原則(`nextjs-best-practices.md` より):

- **layout / page はデフォルトで Server Component**。`"use client"` は末端の対話的コンポーネントだけに付け、境界をツリーの下へ下げる。
- **Context Provider は Client Component に切り出し**、layout では `{children}` だけを包む(`<html>` 全体を包まない)。
- `cookies()` / `headers()` / `params` / `searchParams` は **Promise**(`await` 必須)。使うとそのルートは Dynamic Rendering にオプトインするため、必要箇所は `<Suspense>` で包む。
- エラーは `error.tsx` / `app/global-error.tsx`、404 は `not-found.tsx`。

## 命名・import 規約

- パスエイリアスは `@/`(= `web/src/`)。例: `import { cn } from "@/utils/cn";`
- ディレクトリ名・ファイル名は **kebab-case**(例: `ui-store.ts`、`login-form.tsx`)。React コンポーネントの export 名は PascalCase。
- `lib/api/generated.ts` と `lib/api/model/` は **生成物 — 手で編集しない**(orval が `clean: true` で毎回上書き)。型を変えるときは `schema/openapi/openapi.yaml` を編集 → `bun run generate`。

# ホーム画面 UI 実装プラン(`feature/home-ui`)

## Context

マッチングアプリ「Spark」の認証後ホーム画面(`/home`)を実装する。現状の `home/page.tsx` は「ログインしました(仮)」だけで、`(app)/layout.tsx` も認証ガードのみ。Figma デザイン(SP / PC の 2 レイアウト)に沿って、ナビゲーションシェルと 3 セクション(今日のピックアップ / 新着 / おすすめのユーザー)を作る。

**方針(ユーザー合意済み):**
- **フロント先行**。バックエンドは別担当。OpenAPI / 生成物には一切手を入れない(スキーマ駆動・生成物編集禁止)。
- 3 セクションのうち **「今日のピックアップ」だけ実 API `GET /users/pickup` が既存**。**新着・おすすめは API 未実装なので一時モック**で作り、将来 API 化時に最小 diff で差し替えられるよう隔離する。
- 作業環境: worktree `/Users/ryota/Developer/personal/spark-worktrees/feature-home-ui`(origin/main 基準、ブランチ `feature/home-ui`)。

**前提として確認済みの事実:**
- データ取得は **Server Component + `apiFetch<T>(path, undefined, { auth: true })`** が確立した方針(`src/lib/api/fetcher.ts`、server-only、戻り値 `ApiResult<T> = {ok:true,status,data} | {ok:false,status,error}`)。orval の SWR フックは Bearer 付与なし・プロキシ無しのため**使わない**。
- 認証ガードは `(app)/layout.tsx` の `await requireSession()`(`src/lib/auth/session.ts`)。
- Client ナビの手本は `src/features/auth/components/auth-tabs.tsx`(`"use client"` + `usePathname` + `next/link` + `cn`)。クラス結合は `@/utils/cn`。
- デザイントークンは Tailwind utility に露出済み(`src/app/globals.css`、v4 CSS-first・config 無し):`text-ink`(#2f2f2f)/`text-secondary`(#a7a7a7)/`border-border`(#ddd)/`text-error`、グラデ `bg-brand-gradient`・`bg-brand-gradient-top`・`bg-brand-panel`・`bg-auth-header`、`--brand-yellow/orange/red`。
- アイコングラデは root layout 常設の `IconGradientDefs` により lucide に `stroke="url(#icon-gradient)"` を付けるだけで適用可。
- 生成型: `PickupUsersResponse { users?: PickupUserResponse[] }` / `PickupUserResponse { userId?, name?, bio?:string|null, iconUrl?:string|null, matchedTags?:ThingResponse[], unmatchedTags?:ThingResponse[] }` / `ThingResponse { id?, name?, aliases?, createdAt? }` / `ProfileResponse`。`matchedTags.length` =「共通のタグが N 個」。

---

## 実装内容

### 1. 共通 UI 部品(`src/components/ui/`、presentational・Server 可)
- `avatar.tsx` — 円形アバター。`<img>` + `object-cover` + size variant(sm/md/lg)、画像欠落時はイニシャル fallback。**`next/image` は使わない**(`next.config.ts` に `remotePatterns` 未設定。設定を触らず `<img>` で実装)。
- `tag-chip.tsx` — タグチップ。`variant: filled`(共通タグ=塗り)/ `outline`(非共通=枠線 `border-border`)。
- `section-header.tsx` — 見出し + 右上「もっとみる >」(`next/link` + lucide `ChevronRight`)。

### 2. ホーム feature(`src/features/home/`)
- `types.ts` — UI 用ローカル View Model(生成型の optional 地獄を UI から隠す):
  ```ts
  type TagVM = { id: number; name: string; matched: boolean };
  type PickupCardVM = { userId: number; name: string; bio: string; iconUrl: string | null; matchedCount: number; tags: TagVM[] };
  type NewArrivalVM = { userId: number; name: string; iconUrl: string | null };
  type RecommendedUserVM = { userId: number; name: string; bio: string; iconUrl: string | null; matchedCount: number; tags: TagVM[] };
  ```
- `fixtures.ts` — 新着・おすすめのモック実体(**将来この 1 ファイル削除で済むよう隔離**)。
- `data.ts` — `import "server-only"`。データアクセスを集約し、**シグネチャを将来も不変に保つ**:
  - `getPickupUsers(): Promise<PickupCardVM[] | { error: true }>` → `apiFetch<PickupUsersResponse>("/users/pickup", undefined, { auth:true })` を `toPickupCardVM` で変換。`ok:false` は `{error:true}` を返しセクション単位で degrade。
  - `getNewArrivals(): Promise<NewArrivalVM[]>` → 当面 `fixtures` を返す(`// TODO(api): GET /users/new-arrivals 実装後に apiFetch へ`)。
  - `getRecommendedUsers(): Promise<RecommendedUserVM[]>` → 同上(`// TODO(api): GET /users/recommended`)。
  - マッピング規約: `name ?? ""`、`bio ?? ""`、`matchedCount = matchedTags?.length ?? 0`、`tags = [...matched(matched:true), ...unmatched(matched:false)]`、`id`/`name` 欠落要素は除外。
- `components/`:
  - `pickup-card.tsx`(Server)— グラデヘッダ(アバター+名前+「共通のタグが N 個あります」)+ bio + タグ群(matched=filled / unmatched=outline)。
  - `pickup-carousel.tsx`(**Client**)— CSS scroll-snap(`snap-x snap-mandatory overflow-x-auto`、各スライド `snap-center shrink-0 w-[85%] md:w-[420px]`)+ ドット追従(`IntersectionObserver`、`threshold:0.6`)。ドットクリックで `scrollIntoView({behavior:"smooth", inline:"center"})`。1 件以下はドット非表示。`role="region" aria-label`。`cards: PickupCardVM[]` を受けて内部で `<PickupCard>` を map。
  - `pickup-section.tsx`(Server)— 見出し + carousel。`{error:true}` / 空 / 正常の出し分け(セクション単位 degrade、ホーム全体は落とさない)。
  - `new-arrivals-section.tsx`(Server)— SectionHeader + 横スクロールの丸アバター+名前列。
  - `recommended-card.tsx`(Server)— アバター+名前+bio+タグ+右側の大きな「N 共通」バッジ(`bg-brand-gradient`)。
  - `recommended-section.tsx`(Server)— SectionHeader + `grid grid-cols-1 md:grid-cols-2 gap-4`。

### 3. アプリシェル(`src/features/app-shell/`)
ナビのアクティブ判定だけ Client、骨格・取得・装飾は Server。
- `nav-items.ts` — ナビ項目定数 `{ href, label, icon, enabled }`。ホーム以外(検索/気になる/通知/設定/マイページ)は **`enabled:false`** → `aria-disabled` + `pointer-events-none`(今回は遷移先を作らない)。
- `components/sidebar-nav.tsx`(**Client**)— PC 左ナビ本体。`usePathname` でハイライト(auth-tabs パターン踏襲)。
- `components/bottom-tab-bar.tsx`(**Client**)— SP 下部固定タブ。`fixed inset-x-0 bottom-0 z-20`、中央 MATCH 円ボタン(`-translate-y` + `bg-brand-gradient` + `rounded-full`)、`pb-[env(safe-area-inset-bottom)]`。
- `components/mobile-header.tsx`(Server)— SP 上部グラデヘッダ(「スパーク / いい出会いがあるかも」+ 通知ベル)。`sticky top-0 z-10`。
- `components/sidebar.tsx`(Server)— PC サイドの器。`w-[260px] shrink-0 sticky top-0 h-dvh flex flex-col`。ロゴ / 「マッチ一覧」見出し / `<SidebarNav/>` / `mt-auto` で最下部にミニプロフィール + 設定。
- `components/mini-profile.tsx`(Server)— 自分のミニプロフィール(props で受ける)。
- `data.ts` — `getMyMiniProfile()`:`apiFetch<ProfileResponse>("/profiles/me", undefined, {auth:true})`。**404(未作成)は正常系**としてプレースホルダ(名前「ゲスト」・アバター fallback・ハンドル非表示)にフォールバック。

### 4. `(app)/layout.tsx` 改修
`requireSession()` の後にミニプロフィール取得し、シェルを組む:
```tsx
const me = await getMyMiniProfile();
return (
  <div className="flex min-h-dvh bg-white">
    <Sidebar me={me} className="hidden md:flex" />
    <div className="flex min-w-0 flex-1 flex-col">
      <MobileHeader className="md:hidden" />
      {children}
      <BottomTabBar className="md:hidden" />
    </div>
  </div>
);
```

### 5. `home/page.tsx`
Server Component。`data.ts` の 3 関数で取得し、3 セクションを縦に配置。コンテンツ幅 `mx-auto w-full max-w-[760px] px-4 md:px-8`、SP は下部タブに隠れないよう `pb-24 md:pb-0`。PC のみ `<h1 className="hidden md:block">ホーム</h1>`。

---

## 主要ファイル(新規・改修)
- 改修: `web/src/app/(app)/layout.tsx`、`web/src/app/(app)/home/page.tsx`
- 新規: `web/src/components/ui/{avatar,tag-chip,section-header}.tsx`
- 新規: `web/src/features/home/{types,fixtures,data}.ts` + `web/src/features/home/components/{pickup-card,pickup-carousel,pickup-section,new-arrivals-section,recommended-card,recommended-section}.tsx`
- 新規: `web/src/features/app-shell/{nav-items.ts,data.ts}` + `web/src/features/app-shell/components/{sidebar,sidebar-nav,bottom-tab-bar,mobile-header,mini-profile}.tsx`
- **不変更**: `web/src/lib/api/generated/**`、`schema/openapi/**`、`web/next.config.ts`

## 再利用する既存資産
- 取得: `apiFetch`(`web/src/lib/api/fetcher.ts`)、`requireSession`(`web/src/lib/auth/session.ts`)
- Client ナビ手本: `web/src/features/auth/components/auth-tabs.tsx`、`cn`(`web/src/utils/cn.ts`)
- レスポンシブ出し分け: `web/src/app/(auth)/layout.tsx`(`hidden md:flex` / `md:hidden` / `flex-1` + 固定幅)
- トークン/グラデ: `web/src/app/globals.css`、アイコングラデ: `IconGradientDefs`
- 生成型: `@/lib/api/generated/model`(import のみ)

## 実装順序
1. 共通 UI(avatar → tag-chip → section-header)
2. home の types → fixtures → data
3. home セクション(pickup-card → pickup-carousel → pickup-section、new-arrivals-section、recommended-card → recommended-section)
4. home/page.tsx
5. app-shell(nav-items → sidebar-nav/bottom-tab-bar → mobile-header/mini-profile/sidebar → data)
6. (app)/layout.tsx 改修
7. レスポンシブ微調整(固定要素余白 `pb-24`、`min-w-0`、`sticky`/`fixed`、safe-area)

## 検証
- 型: `cd web && bunx tsc --noEmit`(exit 0)
- Lint/Format: `bun run check`(Biome)。必要に応じ `bun run lint` / `bun run format`
- 目視: リポジトリ root で `make up`(web:3000 / mock:3001)→ ログイン後 `/home`
  - SP(<768px): 上部グラデヘッダ・下部タブバー・1 カラム・カルーセルのスナップ/ドット追従
  - PC(≥768px): 左サイドバー・「ホーム」見出し・おすすめ 2 カラム・最下部ミニプロフィール
  - ピックアップ実 API: `apiFetch("/users/pickup", …, {auth:true})` が Bearer 付きで到達、`matchedTags.length` =「共通のタグが N 個」、filled/outline がそれぞれ matched/unmatched と一致
  - degrade 確認: `/profiles/me` が 404 でもミニプロフィールがフォールバックして落ちない / ピックアップ 0 件・エラー時もセクション単位で degrade し新着・おすすめは表示継続
- 完了後にコミット(`/push` 等。push は確認後)。`feature/home-ui` は worktree 上で作業

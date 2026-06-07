# ライブラリ使用ガイド

このプロジェクトで追加した UI 系ライブラリの使い方と規約。バージョンは `web/package.json` を正とする(本書作成時: `clsx@2.1.1` / `tailwind-merge@3.6.0` / `zustand@5.0.14` / `motion@12.40.0`)。API を使う前に `web/node_modules/<pkg>` の型定義で裏取りすること。

## cn() — クラス結合 (`@/utils/cn`)

`clsx` + `tailwind-merge` の合成。実体は `web/src/utils/cn.ts`。

- **役割**: 条件付きで Tailwind クラスを結合し、競合するユーティリティを**後勝ち**で解決する。
- **使い方**:

  ```tsx
  import { cn } from "@/utils/cn";

  <div className={cn("px-2 py-1", isActive && "bg-blue-500", className)} />
  // px が競合しても後の指定が優先される
  ```

- **規約**: 複数クラスを動的に組み立てる箇所、特に `className` を props で受け取って上書き合成する再利用コンポーネントでは必ず `cn()` を通す(素の文字列連結やテンプレートリテラルで Tailwind クラスを混ぜない)。

## zustand — クライアント状態管理 (v5)

グローバル/横断的なクライアント状態に使う。`useSWR`(サーバ状態のキャッシュ)とは役割を分ける — **サーバから取得したデータは SWR、UI/セッション等のクライアント状態は zustand**。

- **import**: `import { create } from "zustand";`
- **store の定義**(v5 は TS 推論のため `create<T>()(...)` のカリー形を使う):

  ```ts
  // src/stores/ui-store.ts
  import { create } from "zustand";

  type UiState = {
    sidebarOpen: boolean;
    toggleSidebar: () => void;
  };

  export const useUiStore = create<UiState>()((set) => ({
    sidebarOpen: false,
    toggleSidebar: () => set((s) => ({ sidebarOpen: !s.sidebarOpen })),
  }));
  ```

- **読み出しはセレクタで絞る**(不要な再レンダリングを防ぐ):

  ```tsx
  const sidebarOpen = useUiStore((s) => s.sidebarOpen);
  ```

- **注意**:
  - フックなので **Client Component (`"use client"`) でのみ**使える。
  - store はモジュールスコープのシングルトン。**秘密情報(トークン等)を安易に保持しない**(`nextjs-best-practices.md` の認証方針を参照)。
  - 置き場所は `src/stores/` に集約する。

## motion — アニメーション (`motion/react`, v12)

旧 framer-motion の後継。`motion/react` から import する。

- **import**: `import { motion, AnimatePresence } from "motion/react";`
- **必ず `"use client"`**: `motion/react` はクライアント専用。アニメーションを使うコンポーネントの先頭に `"use client"` を付ける(付けずに Server Component で使うとエラー)。

  ```tsx
  "use client";

  import { motion, AnimatePresence } from "motion/react";

  export function FadeIn({ children }: { children: React.ReactNode }) {
    return (
      <motion.div
        initial={{ opacity: 0, y: 8 }}
        animate={{ opacity: 1, y: 0 }}
        exit={{ opacity: 0 }}
        transition={{ duration: 0.2 }}
      >
        {children}
      </motion.div>
    );
  }
  ```

- **入退場アニメーション**は `<AnimatePresence>` で要素を包む(exit を効かせるため)。
- **アクセシビリティ**: モーション過敏のユーザー向けに `useReducedMotion()` で分岐するか、上位を `<MotionConfig reducedMotion="user">` で包む。
- **`motion/react-client`** は `"use client"` 済みの motion DOM 要素のみを提供するサブパス。Server Component から `motion.div` 相当だけ使いたい特殊ケース用で、`AnimatePresence` 等は含まれない。基本は **`"use client"` + `motion/react`** を使う。

## 使い分けの早見

| やりたいこと | 使うもの |
|---|---|
| Tailwind クラスの条件結合・上書き合成 | `cn()` |
| サーバ取得データのキャッシュ/再検証 | `useSWR`(orval 生成) |
| 横断的なクライアント UI 状態 | `zustand` |
| アニメーション/トランジション | `motion/react`(要 `"use client"`) |

## UI コンポーネント方針(shadcn/ui は「将来導入」前提)

現状は **shadcn/ui を入れない**。認証画面など初期 UI は静的要素が中心で旨味が薄いため、`cn` + `cva`(class-variance-authority)で**自作**する。ただし**将来 shadcn を継ぎ足せるよう、最初から shadcn 互換の作法で書く**(モーダル/ドロップダウン/Select 等の複雑なインタラクティブ要素が出てきたら `npx shadcn add ...` で導入する想定)。

互換のための規約:

1. **配置**: 汎用 UI 部品は `src/components/ui/` に置く(shadcn がコンポーネントをコピーする先と同じ)。
2. **実装パターン**: `cn` + `cva` ベースで書く(shadcn 内部と同じ構造)。バリアントは `cva`、最終的なクラス合成は `cn`。
3. **デザイントークン**: 色・角丸・スペーシング等は `src/app/globals.css` に **CSS 変数**として定義し、Tailwind から参照する(shadcn のテーマ方式に合わせる)。Figma のデザイントークンはここへ集約する。

→ これにより「今は軽量・自作、将来は shadcn を継ぎ足すだけ」という低コストな移行パスを確保する。複雑なインタラクティブ要素が必要になった時点で、Radix UI ベースの shadcn コンポーネントを導入する。

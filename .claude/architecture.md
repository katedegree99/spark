---
name: architecture
description: api/ のクリーンアーキテクチャルールと依存注入規約。api/ を編集するとき適用
---

## レイヤー構成

```
api/internal/
├── domain/          # エンティティ・リポジトリインターフェース（外部依存なし）
│   ├── model/       # GORM モデル（DB 構造体）
│   └── repository/  # リポジトリインターフェース定義
├── usecase/         # ビジネスロジック（domain にのみ依存）
├── infrastructure/  # 外部技術の実装（GORM, 外部 API など）
│   ├── db/          # GORM DB 接続
│   └── repository/  # リポジトリインターフェースの GORM 実装
└── adapter/         # フレームワーク適合層（Echo）
    ├── handler/     # StrictServerInterface 実装
    └── router/      # Echo ルーティング
```

## 依存方向（厳守）

```
adapter → usecase → domain ← infrastructure
```

- `domain` は他のどの層にも依存しない
- `usecase` は `domain` のみ参照する（infrastructure を直接使わない）
- `adapter` は `usecase` のみ参照する（repository を直接使わない）
- `infrastructure` は `domain` のインターフェースを実装する

## 依存注入（dig）

- 全ての依存は `api/container.go` の `NewContainer()` でのみ登録する
- 各コンストラクタ（`NewXxx`）は引数でのみ依存を受け取る（グローバル変数禁止）
- 新しい依存を追加したら必ず `container.go` に `c.Provide(...)` を追記する

## DB マイグレーション

- `gorm.AutoMigrate` は使わない
- マイグレーションは `api/migrations/` に連番 SQL ファイルで管理する（例: `001_create_users.sql`）
- `domain/model/` の GORM モデルと `migrations/` の SQL は常に同期を保つ

## ハンドラー実装規約

- `api/pkg/generated/api.gen.go` の `StrictServerInterface` を必ず実装する
- panic の代わりに sentinel error を定義して usecase から返す
- エラーマッピング（sentinel → HTTP ステータス）は handler 層で行う

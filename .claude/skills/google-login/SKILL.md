---
name: google-login
description: Google OAuth の認証情報（Client ID / Secret）が Secret Manager に格納済みかを確認し、未設定であれば設定手順を案内する skill
---

## 概要

Google OAuth を使う前に、Secret Manager への認証情報の格納状況を確認し、
未設定の場合はユーザーに手順を案内する。
`feature/auth-google` の実装前や `GOOGLE_CLIENT_ID` を使う操作の前に自動起動する。

## 起動条件

- Google OAuth / Google ログイン関連の実装を行う前
- `api/.env` の `GOOGLE_CLIENT_ID` をダミー値から更新したいとき
- `/google-login` 明示起動

## 実行フロー

1. `gcloud auth list` を実行して gcloud 認証状態を確認する
2. **未認証の場合** → 以下を案内する

```
gcloud が未認証です。ターミナルで以下を実行してください：

  ! gcloud auth login

ブラウザが開くので Google アカウントでログインしてください。
完了したら教えてください。
```

3. **認証済みの場合** → Secret Manager のシークレット一覧を確認する

```bash
gcloud secrets versions list google-oauth-client-id --project=<project_id>
gcloud secrets versions list google-oauth-client-secret --project=<project_id>
```

4. **シークレットが未設定（バージョンなし）の場合** → 以下の手順を案内する

```
Google OAuth 認証情報がまだ Secret Manager に格納されていません。
以下の手順で設定してください：

1. OAuth コンセント画面を作成:
   ! gcloud alpha iap oauth-brands create \
       --application_title="Spark" \
       --support_email="<メールアドレス>" \
       --project=<project_id>

2. OAuth クライアントを作成:
   ! gcloud alpha iap oauth-clients create \
       projects/<project_id>/brands/<brand_id> \
       --display_name="Spark"

3. 発行された client_id / secret を格納:
   ! gcloud secrets versions add google-oauth-client-id \
       --data-file=- <<< "<client_id>"
   ! gcloud secrets versions add google-oauth-client-secret \
       --data-file=- <<< "<client_secret>"

4. ローカル開発用に api/.env を更新:
   GOOGLE_CLIENT_ID=<取得した client_id>
```

5. ユーザーから完了の返答を受けたら、シークレットのバージョンを再確認する
6. 確認できたら中断していた操作を再開する

## ローカル開発での補足

- `api/.env` の `GOOGLE_CLIENT_ID` はローカル動作確認用に直接値を設定する
- 本番では Secret Manager から取得する運用を想定（`api/terraform/gcp/main.tf` 参照）
- `GOOGLE_CLIENT_SECRET` はサーバーサイドのみで使用するため `.env` には含めない

## 注意事項

- `gcloud auth login` はインタラクティブな操作が必要なため Claude が直接実行することはできない
- ユーザーに `! <command>` と入力してもらうことで、Claude Code のターミナルセッションに出力が流れる
- Terraform は `terraform apply` 済みであること（`google_secret_manager_secret` リソースが作成済み）

# 必要な API を有効化
resource "google_project_service" "iap" {
  service            = "iap.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "secretmanager" {
  service            = "secretmanager.googleapis.com"
  disable_on_destroy = false
}

# ---------------------------------------------------------
# OAuth クライアント認証情報を Secret Manager で管理
#
# OAuth クライアント ID / Secret は Google Cloud Console または
# 以下の gcloud コマンドで作成し、シークレットに格納する。
#
# 1. OAuthコンセント画面を作成:
#   gcloud alpha iap oauth-brands create \
#     --application_title="Spark" \
#     --support_email="<email>" \
#     --project=<project_id>
#
# 2. OAuthクライアントを作成:
#   gcloud alpha iap oauth-clients create \
#     projects/<project_id>/brands/<brand_id> \
#     --display_name="Spark"
#
# 3. 発行された client_id / secret を以下のシークレットに格納:
#   gcloud secrets versions add google-oauth-client-id \
#     --data-file=- <<< "<client_id>"
#   gcloud secrets versions add google-oauth-client-secret \
#     --data-file=- <<< "<client_secret>"
# ---------------------------------------------------------

resource "google_secret_manager_secret" "google_client_id" {
  secret_id = "google-oauth-client-id"

  replication {
    auto {}
  }

  depends_on = [google_project_service.secretmanager]
}

resource "google_secret_manager_secret" "google_client_secret" {
  secret_id = "google-oauth-client-secret"

  replication {
    auto {}
  }

  depends_on = [google_project_service.secretmanager]
}

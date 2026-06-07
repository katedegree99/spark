output "secret_manager_client_id_name" {
  description = "Secret Manager secret name for Google OAuth Client ID"
  value       = google_secret_manager_secret.google_client_id.name
}

output "secret_manager_client_secret_name" {
  description = "Secret Manager secret name for Google OAuth Client Secret"
  value       = google_secret_manager_secret.google_client_secret.name
}

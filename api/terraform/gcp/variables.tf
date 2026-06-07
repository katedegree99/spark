variable "project_id" {
  description = "Google Cloud Project ID"
  type        = string
}

variable "region" {
  description = "Google Cloud region"
  type        = string
  default     = "asia-northeast1"
}

variable "app_name" {
  description = "Application name displayed on OAuth consent screen"
  type        = string
  default     = "Spark"
}

variable "support_email" {
  description = "Support email shown on OAuth consent screen"
  type        = string
}

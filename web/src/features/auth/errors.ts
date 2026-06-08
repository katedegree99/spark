import type {
	ErrorResponse,
	ValidationErrorResponse,
} from "@/lib/api/generated/model";

/** apiFetch がエラー時に返す body(ErrorResponse | ValidationErrorResponse | null)。 */
type AuthApiError = ErrorResponse | ValidationErrorResponse | null;

/**
 * backend の `error.code`(真実源)→ ユーザー向け日本語メッセージの対応表。
 * コードは `schema/openapi/openapi.yaml` と `api/.../handler/auth.go` に準拠。
 * status 番号ではなく code で分岐するため、HTTP ステータス体系に依存しない。
 */
const MESSAGES_BY_CODE: Record<string, string> = {
	EMAIL_ALREADY_EXISTS: "このメールアドレスは既に登録されています",
	INVALID_CREDENTIALS: "メールアドレスまたはパスワードが正しくありません",
	INVALID_OTP: "認証コードが正しくありません",
	INVALID_GOOGLE_TOKEN: "Google アカウントの認証に失敗しました",
	INVALID_REFRESH_TOKEN:
		"セッションの有効期限が切れました。再度ログインしてください",
};

function hasValidationErrors(
	error: AuthApiError,
): error is ValidationErrorResponse {
	return !!error && "errors" in error && Array.isArray(error.errors);
}

/**
 * API エラー → 表示メッセージへの写像を 1 箇所に集約する。
 * 各 Server Action はこの関数に委譲し、status 番号の散在をなくす。
 *
 * 優先順位:
 * 1. 既知の `error.code` があればコード基準で写像(status 非依存)
 * 2. VALIDATION_ERROR(422)は先頭のフィールドエラー文言を採用(握り潰さない)
 * 3. いずれも該当しなければ action ごとの `fallback`
 */
export function authErrorMessage(
	error: AuthApiError,
	fallback: string,
): string {
	if (error?.code && MESSAGES_BY_CODE[error.code]) {
		return MESSAGES_BY_CODE[error.code];
	}
	if (hasValidationErrors(error) && error.errors.length > 0) {
		return error.errors[0].message;
	}
	return fallback;
}

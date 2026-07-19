import type {
	ErrorResponse,
	ValidationErrorResponse,
} from "@/lib/api/generated/model";

/** apiFetch がエラー時に返す body(ErrorResponse | ValidationErrorResponse | null)。 */
type ProfileApiError = ErrorResponse | ValidationErrorResponse | null;

/**
 * backend の `error.code` → ユーザー向け日本語メッセージの対応表。
 * コードは `schema/openapi/openapi.yaml` の profile / things に準拠。
 * status 番号ではなく code で分岐する(HTTP ステータス体系に依存しない)。
 */
const MESSAGES_BY_CODE: Record<string, string> = {
	THING_ALREADY_EXISTS: "同じ名前の事柄が既に存在します",
};

function hasValidationErrors(
	error: ProfileApiError,
): error is ValidationErrorResponse {
	return !!error && "errors" in error && Array.isArray(error.errors);
}

/**
 * API エラー → 表示メッセージへの写像を 1 箇所に集約する。
 * `auth/errors.ts` と同じ優先順位(code → validation 先頭 → fallback)。
 */
export function profileErrorMessage(
	error: ProfileApiError,
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

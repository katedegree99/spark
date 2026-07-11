/**
 * Server Action 内の `redirect()` 由来の reject かを判定する。
 *
 * Next はアクションが redirect すると、クライアント側の action promise を
 * digest が "NEXT_REDIRECT" で始まるエラーで reject する(ナビゲーション自体は
 * ルーターが別途実行する)。イベントハンドラから action を直接 await する場合、
 * これを捕捉しないと成功時の後続処理が実行されず、未処理 rejection にもなる。
 * ナビゲーションは reject とは独立に実行済みのため、捕捉して握りつぶしてよい。
 */
export function isNextRedirectError(err: unknown): boolean {
	const digest = (err as { digest?: unknown } | null)?.digest;
	return typeof digest === "string" && digest.startsWith("NEXT_REDIRECT");
}

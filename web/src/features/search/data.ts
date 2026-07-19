/**
 * 探すページのデータアクセス層。Server Component からのみ呼ぶ。
 *
 * GET /users を `tagIds`(AND 条件)で叩き、ホームと同じ
 * `RecommendedUserVM` に変換して返す(カードも `RecommendedCard` を再利用)。
 */
import "server-only";
import { toRecommendedUserVM } from "@/features/home/data";
import type { RecommendedUserVM } from "@/features/home/types";
import { apiFetch } from "@/lib/api/fetcher";
import type { RecommendUsersResponse } from "@/lib/api/generated/model";

/** 検索結果の取得失敗を表す(home の PickupError パターン踏襲)。 */
export type SearchError = { error: true };

/**
 * タグ ID(AND 条件)でユーザーを検索する。タグ 0 個ならパラメータなしで
 * 全ユーザーを取得する。取得失敗(`ok:false`)は例外にせず `{ error: true }` を返す。
 *
 * TODO: ページネーション(offset/limit)は今回スコープ外。API デフォルトの
 * limit 20 のまま先頭ページのみ表示する。
 */
export async function searchUsers(
	tagIds: number[],
): Promise<RecommendedUserVM[] | SearchError> {
	const params = new URLSearchParams();
	for (const id of tagIds) {
		params.append("tagIds", String(id));
	}
	const query = params.toString();

	const res = await apiFetch<RecommendUsersResponse>(
		`/users${query ? `?${query}` : ""}`,
		undefined,
		{ auth: true },
	);
	if (!res.ok) return { error: true };
	return (res.data?.users ?? [])
		.map(toRecommendedUserVM)
		.filter((u): u is RecommendedUserVM => u !== null);
}

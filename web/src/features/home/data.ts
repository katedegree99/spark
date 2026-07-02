/**
 * ホーム画面のデータアクセス層。Server Component からのみ呼ぶ。
 *
 * 公開関数のシグネチャは将来も不変に保つ(API 化時に中身だけ差し替える)。
 * 「今日のピックアップ」は実 API `GET /users/pickup` を叩く。新着・おすすめは
 * API 未実装のため当面 `fixtures` を返す。
 */
import "server-only";
import { apiFetch } from "@/lib/api/fetcher";
import type {
	PickupUserResponse,
	PickupUsersResponse,
	ThingResponse,
} from "@/lib/api/generated/model";
import {
	newArrivalsFixture,
	pickupCardsFixture,
	recommendedUsersFixture,
} from "./fixtures";
import type {
	NewArrivalVM,
	PickupCardVM,
	RecommendedUserVM,
	TagVM,
} from "./types";

/** ピックアップ取得失敗を表す。セクション単位で degrade するために使う。 */
export type PickupError = { error: true };

/** `matched`/`unmatched` のタグ群を VM 配列へ。id/name 欠落要素は除外する。 */
function toTagVMs(
	matchedTags?: ThingResponse[],
	unmatchedTags?: ThingResponse[],
): TagVM[] {
	const toVMs = (
		tags: ThingResponse[] | undefined,
		matched: boolean,
	): TagVM[] =>
		(tags ?? [])
			.filter((t): t is ThingResponse & { id: number; name: string } => {
				return t.id != null && t.name != null;
			})
			.map((t) => ({ id: t.id, name: t.name, matched }));

	return [...toVMs(matchedTags, true), ...toVMs(unmatchedTags, false)];
}

/** ピックアップユーザーを VM に変換する。userId 欠落要素は除外(null を返す)。 */
function toPickupCardVM(user: PickupUserResponse): PickupCardVM | null {
	if (user.userId == null) return null;
	return {
		userId: user.userId,
		name: user.name ?? "",
		bio: user.bio ?? "",
		iconUrl: user.iconUrl ?? null,
		matchedCount: user.matchedTags?.length ?? 0,
		tags: toTagVMs(user.matchedTags, user.unmatchedTags),
	};
}

/**
 * 今日のピックアップを取得する。
 * 取得失敗(`ok:false`)は例外にせず `{ error: true }` を返し、ホーム全体は
 * 落とさずこのセクションだけ degrade させる。
 */
export async function getPickupUsers(): Promise<PickupCardVM[] | PickupError> {
	const res = await apiFetch<PickupUsersResponse>("/users/pickup", undefined, {
		auth: true,
	});
	// TODO(api): backend が本番データを返せるようになったら、この dummy フォールバックを外す。
	if (!res.ok) return pickupCardsFixture;

	const cards = (res.data.users ?? [])
		.map(toPickupCardVM)
		.filter((c): c is PickupCardVM => c !== null);

	// mock サーバーはアバター URL が無効(pub-xxx.r2.dev)な単一の例題データしか
	// 返さないため、実データが揃うまでは Unsplash 画像入りのダミーで見た目を確認する。
	const looksLikeMock =
		cards.length <= 1 &&
		cards.every((c) => c.iconUrl?.includes("pub-xxx") ?? false);
	return cards.length > 0 && !looksLikeMock ? cards : pickupCardsFixture;
}

/** 新着ユーザーを取得する。 */
// TODO(api): GET /users/new-arrivals 実装後に apiFetch へ差し替える。
export async function getNewArrivals(): Promise<NewArrivalVM[]> {
	return newArrivalsFixture;
}

/** おすすめユーザーを取得する。 */
// TODO(api): GET /users/recommended 実装後に apiFetch へ差し替える。
export async function getRecommendedUsers(): Promise<RecommendedUserVM[]> {
	return recommendedUsersFixture;
}

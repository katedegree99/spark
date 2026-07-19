/**
 * ホーム画面のデータアクセス層。Server Component からのみ呼ぶ。
 *
 * 3 セクションとも mock/実 API から取得する:
 *   今日のピックアップ → GET /users/pickup
 *   新着             → GET /users/new
 *   おすすめ          → GET /users/recommend
 * アバター URL が無効(mock の pub-xxx placeholder)/未設定のときは
 * `resolveAvatarUrl` でダミー画像にフォールバックする。
 */
import "server-only";
import { apiFetch } from "@/lib/api/fetcher";
import type {
	NewUserResponse,
	NewUsersResponse,
	PickupUserResponse,
	PickupUsersResponse,
	RecommendUserResponse,
	RecommendUsersResponse,
	TagResponse,
} from "@/lib/api/generated/model";
import { resolveAvatarUrl } from "@/lib/dummy-avatar";
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
	matchedTags?: TagResponse[],
	unmatchedTags?: TagResponse[],
): TagVM[] {
	const toVMs = (tags: TagResponse[] | undefined, matched: boolean): TagVM[] =>
		(tags ?? [])
			.filter((t): t is TagResponse & { id: number; name: string } => {
				return t.id != null && t.name != null;
			})
			.map((t) => ({ id: t.id, name: t.name, matched }));

	return [...toVMs(matchedTags, true), ...toVMs(unmatchedTags, false)];
}

/** ピックアップユーザーを VM に変換する。userId 欠落要素は除外(null を返す)。 */
function toPickupCardVM(user: PickupUserResponse): PickupCardVM | null {
	if (user.userId == null) return null;
	const tags = toTagVMs(user.matchedTags, user.unmatchedTags);
	return {
		userId: user.userId,
		name: user.name ?? "",
		bio: user.bio ?? "",
		iconUrl: resolveAvatarUrl(user.iconUrl, user.userId),
		// 「共通のタグが N 個」の N が表示する filled タグ数とズレないよう、
		// id/name 欠落を除外したフィルタ後の VM から数える。
		matchedCount: tags.filter((t) => t.matched).length,
		tags,
	};
}

/** 新着ユーザーを VM に変換する。userId 欠落要素は除外(null を返す)。 */
function toNewArrivalVM(user: NewUserResponse): NewArrivalVM | null {
	if (user.userId == null) return null;
	return {
		userId: user.userId,
		name: user.name ?? "",
		iconUrl: resolveAvatarUrl(user.iconUrl, user.userId),
	};
}

/** おすすめユーザーを VM に変換する。userId 欠落要素は除外(null を返す)。 */
export function toRecommendedUserVM(
	user: RecommendUserResponse,
): RecommendedUserVM | null {
	if (user.userId == null) return null;
	const tags = toTagVMs(user.matchedTags, user.unmatchedTags);
	return {
		userId: user.userId,
		name: user.name ?? "",
		bio: user.bio ?? "",
		iconUrl: resolveAvatarUrl(user.iconUrl, user.userId),
		// commonCount 優先。無ければフィルタ後の matched タグ数で代替する
		// (未フィルタの length だと表示 filled タグ数とズレうるため)。
		matchedCount: user.commonCount ?? tags.filter((t) => t.matched).length,
		tags,
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
	if (!res.ok) return { error: true };
	// res.data は 2xx でも JSON パース失敗時に null になりうるため optional chain で守る。
	return (res.data?.users ?? [])
		.map(toPickupCardVM)
		.filter((c): c is PickupCardVM => c !== null);
}

/**
 * 新着ユーザーを取得する。取得失敗時は空配列(セクション非表示)にして degrade する。
 */
export async function getNewArrivals(): Promise<NewArrivalVM[]> {
	const res = await apiFetch<NewUsersResponse>("/users/new", undefined, {
		auth: true,
	});
	if (!res.ok) return [];
	return (res.data?.users ?? [])
		.map(toNewArrivalVM)
		.filter((u): u is NewArrivalVM => u !== null);
}

/**
 * おすすめユーザーを取得する。取得失敗時は空配列(セクション非表示)にして degrade する。
 */
export async function getRecommendedUsers(): Promise<RecommendedUserVM[]> {
	const res = await apiFetch<RecommendUsersResponse>(
		"/users/recommend",
		undefined,
		{ auth: true },
	);
	if (!res.ok) return [];
	return (res.data?.users ?? [])
		.map(toRecommendedUserVM)
		.filter((u): u is RecommendedUserVM => u !== null);
}

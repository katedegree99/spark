/**
 * アプリシェル(サイドバー等)のデータアクセス層。Server Component 専用。
 */
import "server-only";
import { apiFetch } from "@/lib/api/fetcher";
import type { ProfileResponse } from "@/lib/api/generated/model";
import { resolveAvatarUrl } from "@/lib/dummy-avatar";

export type MiniProfileVM = {
	name: string;
	handle: string | null;
	iconUrl: string | null;
};

/** プロフィール未作成(404)等のフォールバック。アイコンはデフォルト表示に委ねる。 */
const GUEST_PROFILE: MiniProfileVM = {
	name: "ゲスト",
	handle: null,
	iconUrl: null,
};

/**
 * サイドバー最下部に出す自分のミニプロフィールを取得する。
 * `/profiles/me` が 404(プロフィール未作成)でも正常系として扱い、
 * プレースホルダ(ゲスト)へフォールバックして画面を落とさない。
 */
export async function getMyMiniProfile(): Promise<MiniProfileVM> {
	const res = await apiFetch<ProfileResponse>("/profiles/me", undefined, {
		auth: true,
	});
	// res.data は 2xx でも JSON パース失敗時に null になりうる。
	if (!res.ok || res.data == null) return GUEST_PROFILE;

	const profile = res.data;
	return {
		name: profile.name ?? GUEST_PROFILE.name,
		// スキーマ上ハンドル文字列は無く ID は userId(数値)のみ。Figma の「＠〜」表記に
		// 合わせて ＠userId を表示する(未設定時は非表示)。
		handle: profile.userId != null ? `＠${profile.userId}` : null,
		iconUrl: resolveAvatarUrl(profile.iconImage?.url, profile.userId ?? 0),
	};
}

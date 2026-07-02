/**
 * アプリシェル(サイドバー等)のデータアクセス層。Server Component 専用。
 */
import "server-only";
import { apiFetch } from "@/lib/api/fetcher";
import type { ProfileResponse } from "@/lib/api/generated/model";

export type MiniProfileVM = {
	name: string;
	handle: string | null;
	iconUrl: string | null;
};

// TODO(mock): 開発用のダミーアバター。API がプロフィール画像を返すようになったら外す。
const DUMMY_AVATAR =
	"https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=200&h=200&fit=crop&crop=faces&auto=format&q=80";

/** プロフィール未作成(404)等のフォールバック。 */
const GUEST_PROFILE: MiniProfileVM = {
	name: "ゲスト",
	handle: null,
	iconUrl: DUMMY_AVATAR,
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
	if (!res.ok) return GUEST_PROFILE;

	const profile = res.data;
	return {
		name: profile.name ?? GUEST_PROFILE.name,
		// スキーマ上ハンドル文字列は無く ID は userId(数値)のみ。Figma の「＠〜」表記に
		// 合わせて ＠userId を表示する(未設定時は非表示)。
		handle: profile.userId != null ? `＠${profile.userId}` : null,
		iconUrl: profile.iconImage?.url ?? DUMMY_AVATAR,
	};
}

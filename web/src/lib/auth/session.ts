/**
 * 認証セッション(トークン)の Cookie 管理。Server Action / Route Handler 専用。
 *
 * access/refresh JWT は httpOnly Cookie に隔離し、ブラウザ JS へは一切露出させない
 * (XSS 対策。`nextjs-best-practices.md` の方針)。`next/headers` の `cookies()` は
 * サーバ実行限定なので、このモジュールはクライアントから import しない。
 * `server-only` でクライアントからの誤 import をビルド時に弾く。
 */
import "server-only";
import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import type { AuthTokensResponse } from "@/lib/api/generated/model";

const ACCESS_TOKEN = "access_token";
const REFRESH_TOKEN = "refresh_token";

const ACCESS_FALLBACK_MAX_AGE = 60 * 15; // 15分(API が expires_in を返さない場合)
// backend の refresh token DB 有効期限(7日)に合わせる。長くすると Cookie だけ
// 生き残り「Cookie はあるが refresh は失効」状態になるため一致させること。
const REFRESH_MAX_AGE = 60 * 60 * 24 * 7; // 7日

// TODO(次PR): access 失効時の自動リフレッシュが未実装。
// refresh token を保存しているが /auth/refresh を叩く経路が無く、access が
// 15分で切れると refresh が生きていても requireSession が /login へ飛ばす。
// App Router は Server Component 描画中に Cookie を書けないため、Middleware で
// 「access 無 & refresh 有 → /auth/refresh で更新して response に Cookie セット」
// を先回り実装する(あわせて logoutAction + clearAuthCookies の配線も行う)。

/**
 * トークンを httpOnly Cookie に保存する。
 * access は API の `expires_in`、refresh は固定 30 日を既定の寿命とする。
 */
export async function setAuthCookies(
	tokens: AuthTokensResponse,
): Promise<void> {
	const store = await cookies();
	const base = {
		httpOnly: true,
		secure: process.env.NODE_ENV === "production",
		sameSite: "lax" as const,
		path: "/",
	};

	if (tokens.access_token) {
		store.set(ACCESS_TOKEN, tokens.access_token, {
			...base,
			maxAge: tokens.expires_in ?? ACCESS_FALLBACK_MAX_AGE,
		});
	}
	if (tokens.refresh_token) {
		store.set(REFRESH_TOKEN, tokens.refresh_token, {
			...base,
			maxAge: REFRESH_MAX_AGE,
		});
	}
}

/** ログアウト等でトークン Cookie を破棄する。 */
export async function clearAuthCookies(): Promise<void> {
	const store = await cookies();
	store.delete(ACCESS_TOKEN);
	store.delete(REFRESH_TOKEN);
}

/** アクセストークンを Cookie から読む(無ければ null)。サーバ実行限定。 */
export async function getAccessToken(): Promise<string | null> {
	const store = await cookies();
	return store.get(ACCESS_TOKEN)?.value ?? null;
}

/**
 * 認証ガード。アクセストークンが無ければ `/login` へリダイレクトする。
 * 認証後 (`(app)`) の layout / page 冒頭で `await requireSession()` して使う。
 * 戻り値はトークン文字列(リダイレクト時は到達しない)。
 */
export async function requireSession(): Promise<string> {
	const token = await getAccessToken();
	if (!token) {
		redirect("/login");
	}
	return token;
}

import type { ReactNode } from "react";
import { requireSession } from "@/lib/auth/session";

/**
 * 初回セットアップ (`(setup)`) 共通レイアウト。Server Component。
 *
 * 認証ガードのみで、`(app)` と違いアプリシェル(サイドバー・ヘッダ・下部タブ)は
 * 描画しない。プロフィール未設定ユーザーが通る画面(/profile/register)は
 * 独自のヘッダ・サイドパネルを持つため。
 * `cookies()` を読むためこの配下は Dynamic Rendering になる。
 */
export default async function SetupLayout({
	children,
}: {
	children: ReactNode;
}) {
	await requireSession();

	return <div className="flex min-h-dvh flex-col bg-white">{children}</div>;
}

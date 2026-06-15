import type { ReactNode } from "react";
import { requireSession } from "@/lib/auth/session";

/**
 * 認証後 (`(app)`) 共通レイアウト。Server Component。
 *
 * httpOnly Cookie のアクセストークンが無ければ `/login` へ送る認証ガード。
 * 将来サイドバー等のアプリシェルをここに置く(`directory-structure.md`)。
 * `cookies()` を読むためこの配下は Dynamic Rendering になる。
 */
export default async function AppLayout({ children }: { children: ReactNode }) {
	await requireSession();

	return <div className="flex min-h-dvh flex-col bg-white">{children}</div>;
}

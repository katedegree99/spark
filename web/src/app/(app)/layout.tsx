import type { ReactNode } from "react";
import { BottomTabBar } from "@/features/app-shell/components/bottom-tab-bar";
import { MobileHeader } from "@/features/app-shell/components/mobile-header";
import { Sidebar } from "@/features/app-shell/components/sidebar";
import { getMyMiniProfile } from "@/features/app-shell/data";
import { requireSession } from "@/lib/auth/session";

/**
 * 認証後 (`(app)`) 共通レイアウト。Server Component。
 *
 * httpOnly Cookie のアクセストークンが無ければ `/login` へ送る認証ガードを通し、
 * その後アプリシェル(PC 左サイドバー / SP 上部ヘッダ + 下部タブバー)を組む。
 * `cookies()` を読むためこの配下は Dynamic Rendering になる。
 */
export default async function AppLayout({ children }: { children: ReactNode }) {
	await requireSession();
	const me = await getMyMiniProfile();

	return (
		<div className="flex min-h-dvh bg-white">
			<Sidebar me={me} className="hidden md:flex" />
			<div className="flex min-w-0 flex-1 flex-col">
				<MobileHeader className="md:hidden" />
				{children}
				<BottomTabBar className="md:hidden" />
			</div>
		</div>
	);
}

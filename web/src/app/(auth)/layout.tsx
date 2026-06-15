import type { ReactNode } from "react";
import { AuthSidePanel } from "@/features/auth/components/auth-side-panel";

/**
 * 認証前の画面群(login / register / otp)共通レイアウト。
 * - SP (~md): 上部装飾 + フォーム(各ページ側で縦積み)
 * - PC (md 以上): 左サイドパネル + 右にフォームの 2 カラム
 *
 * フォーム側を 508px(460 + 左右余白)固定にし、サイドパネルを残り幅で可変
 * (`flex-1`)にする。旧実装はサイドパネルを 562px 固定にしていたため、md
 * 付近でフォームが潰れていた。可変にすることでフォーム幅を常に保ちつつ、
 * 広い画面ほどサイドパネルが広がる。
 */
export default function AuthLayout({ children }: { children: ReactNode }) {
	return (
		<div className="flex min-h-dvh bg-white">
			<aside className="hidden min-w-0 flex-1 p-6 md:flex">
				<AuthSidePanel />
			</aside>
			<main className="flex w-full flex-1 flex-col md:w-[508px] md:flex-none md:items-center md:justify-center md:p-6">
				<div className="flex w-full flex-1 flex-col md:max-w-[460px] md:flex-none">
					{children}
				</div>
			</main>
		</div>
	);
}

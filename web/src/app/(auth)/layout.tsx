import type { ReactNode } from "react";
import { AuthSidePanel } from "@/features/auth/components/auth-side-panel";

/**
 * 認証前の画面群(login / register / otp)共通レイアウト。
 * - SP: 子要素のみ(各ページが上部装飾 + フォームを縦に並べる)
 * - PC (md 以上): 左サイドパネル + 右に中央寄せした子要素の 2 カラム
 */
export default function AuthLayout({ children }: { children: ReactNode }) {
	return (
		<div className="flex min-h-dvh bg-white">
			<aside className="hidden shrink-0 p-6 md:flex md:w-[562px]">
				<AuthSidePanel />
			</aside>
			<main className="flex flex-1 flex-col md:items-center md:justify-center md:p-6">
				<div className="flex w-full flex-1 flex-col md:max-w-[460px] md:flex-none">
					{children}
				</div>
			</main>
		</div>
	);
}

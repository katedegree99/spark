import type { ReactNode } from "react";
import { AuthSidePanel } from "@/features/auth/components/auth-side-panel";

/**
 * 認証前の画面群(login / register / otp)共通レイアウト。
 * - SP (~md): 上部装飾 + フォーム(各ページ側で縦積み)
 * - PC (md 以上): 左サイドパネル + 右にフォームの 2 カラム
 *
 * サイドパネルを最大 562px(Figma)に抑え、フォーム側を残り幅で中央寄せにする
 * (`profile/register` と同じ方針)。562px 未満ではサイドパネルも可変になり、
 * フォーム幅(max 460px)は常に保たれる。
 */
export default function AuthLayout({ children }: { children: ReactNode }) {
	return (
		<div className="flex min-h-dvh bg-white">
			<aside className="hidden min-w-0 flex-1 p-6 md:flex md:max-w-[562px]">
				<AuthSidePanel />
			</aside>
			<main className="flex w-full flex-1 flex-col md:items-center md:justify-center md:p-6">
				<div className="flex w-full flex-1 flex-col md:max-w-[460px] md:flex-none">
					{children}
				</div>
			</main>
		</div>
	);
}

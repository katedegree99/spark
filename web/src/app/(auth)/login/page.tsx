import { AuthTabs } from "@/features/auth/components/auth-tabs";
import { LoginForm } from "@/features/auth/components/login-form";

/**
 * ログイン画面 (`/login`)。Server Component。
 * - SP: 上部のグラデ装飾(白い波線 SVG)+ 白カードを被せる + フォーム
 * - PC: サイドパネル(layout)+ タブ + フォーム
 */
export default function LoginPage() {
	return (
		<div className="flex flex-1 flex-col">
			{/* SP: グラデ装飾ヘッダー + 白い装飾 SVG */}
			<div className="relative flex h-[240px] w-full shrink-0 items-end overflow-hidden bg-brand-gradient-top pb-[26px] pl-4 md:hidden">
				{/* biome-ignore lint/performance/noImgElement: 装飾 SVG は最適化不要のため img を使う */}
				<img
					src="/images/auth-header.svg"
					alt=""
					aria-hidden="true"
					className="h-[144px] w-full object-cover"
				/>
			</div>

			{/* 白カード(SP: 装飾グラデに被せる / PC: 通常) */}
			<div className="relative -mt-7 flex flex-1 flex-col rounded-t-[28px] bg-white px-5 pt-8 pb-8 md:mt-0 md:rounded-none md:px-0 md:pt-0 md:pb-0">
				{/* SP: 見出し / PC: タブ */}
				<h1 className="font-medium text-[32px] text-ink tracking-[2px] md:hidden">
					ログイン
				</h1>
				<div className="hidden md:block">
					<AuthTabs />
				</div>
				<LoginForm />
			</div>
		</div>
	);
}

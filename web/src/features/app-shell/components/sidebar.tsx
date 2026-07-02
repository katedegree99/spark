import { Settings } from "lucide-react";
import { cn } from "@/utils/cn";
import type { MiniProfileVM } from "../data";
import { MiniProfile } from "./mini-profile";
import { SidebarNav } from "./sidebar-nav";

/**
 * PC 左サイドバーの器(Server、Figma 準拠)。
 * 上にロゴ、中央に「マッチ一覧」グラデ見出し + ナビ(下部に設定)、
 * 最下部に区切り線付きのミニプロフィール。`sticky top-0 h-dvh` で追従。
 */
export function Sidebar({
	me,
	className,
}: {
	me: MiniProfileVM;
	className?: string;
}) {
	return (
		<aside
			className={cn(
				"sticky top-0 flex h-dvh w-80 shrink-0 flex-col justify-between border-border border-r bg-white px-4 py-6",
				className,
			)}
		>
			{/* ロゴ */}
			<div className="flex items-center gap-3">
				<div className="size-15 rounded-lg bg-brand-gradient" />
				<span className="bg-brand-gradient bg-clip-text font-bold text-2xl text-transparent tracking-wide">
					SPARK
				</span>
			</div>

			{/* ナビ(flex-1 で上下に伸ばし、設定を下端へ) */}
			<nav
				className="flex flex-1 flex-col justify-between pt-12 pb-3"
				aria-label="メインメニュー"
			>
				<div className="flex flex-col gap-2">
					<div className="flex items-center rounded-2xl bg-brand-gradient p-4">
						<span className="font-bold text-white text-xl">マッチ一覧</span>
					</div>
					<SidebarNav />
				</div>
				<button
					type="button"
					aria-disabled="true"
					className="pointer-events-none flex items-center gap-5 rounded-xl p-3 font-bold text-ink text-xl tracking-wide opacity-60"
				>
					<Settings className="size-10 shrink-0 p-1.5" aria-hidden="true" />
					設定
				</button>
			</nav>

			{/* ミニプロフィール */}
			<div className="border-border border-t pt-6">
				<MiniProfile me={me} />
			</div>
		</aside>
	);
}

import { HeartHandshake, Settings } from "lucide-react";
import Image from "next/image";
import { cn } from "@/utils/cn";
import type { MiniProfileVM } from "../data";
import { MiniProfile } from "./mini-profile";
import { SidebarNav } from "./sidebar-nav";

/**
 * PC 左サイドバーの器(Server、Figma 準拠)。
 * 上にロゴ、中央に「マッチ一覧」グラデ見出し + ナビ(下部に設定)、
 * 最下部に区切り線付きのミニプロフィール。`sticky top-0 h-dvh` で追従。
 *
 * レスポンシブ: フル幅(w-80)だと md〜lg の中間幅でメインコンテンツが潰れるため、
 * lg 未満はテキストを隠したアイコンのみの細い表示(w-22)にし、lg 以上でフル表示にする。
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
				// コンテンツ最小高がビューポートを超える低い画面でも最下部のミニプロフィール
				// へ到達できるよう、固定高 + overflow-y-auto でサイドバー内スクロールを許可する。
				"sticky top-0 flex h-dvh w-22 shrink-0 flex-col justify-between overflow-y-auto overscroll-contain border-border border-r bg-white px-3 py-6 lg:w-80 lg:px-4",
				className,
			)}
		>
			{/* ロゴ(lg 未満はグラデタイルのみ)。タイルの上に白のロゴマークを重ねる */}
			<div className="flex items-center justify-center gap-3 lg:justify-start">
				<div className="flex size-15 shrink-0 items-center justify-center rounded-lg bg-brand-gradient p-2.5">
					<Image
						src="/images/spark-logo.png"
						alt="SPARK"
						width={816}
						height={831}
						className="size-full object-contain"
					/>
				</div>
				<span className="hidden bg-brand-gradient bg-clip-text font-bold text-2xl text-transparent tracking-wide lg:inline">
					SPARK
				</span>
			</div>

			{/* ナビ(flex-1 で上下に伸ばし、設定を下端へ) */}
			<nav
				className="flex flex-1 flex-col justify-between pt-12 pb-3"
				aria-label="メインメニュー"
			>
				<div className="flex flex-col gap-2">
					<div className="flex items-center justify-center rounded-2xl bg-brand-gradient p-4 lg:justify-start">
						<HeartHandshake
							className="size-8 shrink-0 text-white lg:hidden"
							strokeWidth={1.5}
							aria-hidden="true"
						/>
						{/* lg 未満は視覚上アイコンのみ(テキストは sr-only で読み上げには残す) */}
						<span className="sr-only font-bold text-white text-xl lg:not-sr-only">
							マッチ一覧
						</span>
					</div>
					<SidebarNav />
				</div>
				<button
					type="button"
					aria-disabled="true"
					aria-label="設定"
					className="pointer-events-none flex items-center gap-5 rounded-xl p-3 font-bold text-ink text-xl tracking-wide opacity-60"
				>
					<Settings
						className="size-10 shrink-0"
						strokeWidth={1.5}
						aria-hidden="true"
					/>
					<span className="hidden lg:inline">設定</span>
				</button>
			</nav>

			{/* ミニプロフィール */}
			<div className="border-border border-t pt-6">
				<MiniProfile me={me} />
			</div>
		</aside>
	);
}

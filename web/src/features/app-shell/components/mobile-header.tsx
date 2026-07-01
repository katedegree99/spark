import { Bell } from "lucide-react";
import { cn } from "@/utils/cn";

/**
 * SP 上部のグラデヘッダ(Server、Figma 準拠)。
 * サービス名 + キャッチに、右側へ白丸の通知ベル。`sticky top-0` で上部に残す。
 * `notificationCount` > 0 でベル右上に赤い未読バッジ(Figma: #ff3434)を表示。
 */
export function MobileHeader({
	className,
	notificationCount = 0,
}: {
	className?: string;
	notificationCount?: number;
}) {
	const hasNotification = notificationCount > 0;
	const badgeLabel = notificationCount > 99 ? "99+" : String(notificationCount);
	return (
		<header
			className={cn(
				"sticky top-0 z-10 overflow-hidden rounded-b-[20px] bg-brand-gradient-top text-white",
				className,
			)}
		>
			{/* Figma: brand グラデの上に白 20% を重ねて淡いトーンにする */}
			<div className="absolute inset-0 bg-white/20" aria-hidden="true" />
			<div className="relative z-10 flex items-center justify-between px-5 py-5">
				<div className="flex flex-col gap-2">
					<span className="font-semibold text-2xl leading-tight">スパーク</span>
					<span className="text-sm text-white">いい出会いがあるかも</span>
				</div>
				<button
					type="button"
					aria-label={
						hasNotification ? `通知(未読 ${notificationCount} 件)` : "通知"
					}
					className="relative flex items-center justify-center rounded-full bg-white p-2 drop-shadow-[0px_2px_2px_rgba(255,255,255,0.25)]"
				>
					<Bell
						className="size-8"
						aria-hidden="true"
						stroke="url(#icon-gradient)"
						strokeWidth={1.5}
					/>
					{hasNotification && (
						<span
							aria-hidden="true"
							className="-right-1 -top-1 absolute flex size-6 items-center justify-center rounded-full bg-error px-1 font-bold text-white text-xs tabular-nums leading-none"
						>
							{badgeLabel}
						</span>
					)}
				</button>
			</div>
		</header>
	);
}

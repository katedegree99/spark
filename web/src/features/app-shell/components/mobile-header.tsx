"use client";

import { Bell } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import { cn } from "@/utils/cn";

/**
 * SP 上部のグラデヘッダ(Client、Figma 準拠)。
 * サービス名 + キャッチに、右側へ白丸の通知ベル。`notificationCount` > 0 で
 * ベル右上に赤い未読バッジ(Figma: #ff3434)を表示。
 *
 * 下スクロールで隠し、上スワイプで自然に出す(オートハイド)。常時 sticky top-0 で
 * 残すと下のグラデカードと色が被るため、隠す + 影で背景から分離する。
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

	// スクロール方向でヘッダの表示/非表示を切り替える。
	const [hidden, setHidden] = useState(false);
	const lastY = useRef(0);

	useEffect(() => {
		function onScroll() {
			const y = window.scrollY;
			const delta = y - lastY.current;
			// 上端付近は常に表示。少し下れば隠し、上スワイプで出す(微小移動は無視)。
			if (y < 80) {
				setHidden(false);
			} else if (delta > 4) {
				setHidden(true);
			} else if (delta < -4) {
				setHidden(false);
			}
			lastY.current = y;
		}

		// ヘッダは PC(md 以上)では `md:hidden` で非表示だがマウントはされるため、
		// scroll 購読は SP 幅のときだけ張る(ブレークポイント跨ぎで張り替える)。
		const mq = window.matchMedia("(min-width: 768px)");
		let subscribed = false;
		function sync() {
			if (mq.matches) {
				if (subscribed) {
					window.removeEventListener("scroll", onScroll);
					subscribed = false;
				}
				setHidden(false);
				return;
			}
			if (!subscribed) {
				lastY.current = window.scrollY;
				window.addEventListener("scroll", onScroll, { passive: true });
				subscribed = true;
			}
		}
		sync();
		mq.addEventListener("change", sync);
		return () => {
			mq.removeEventListener("change", sync);
			if (subscribed) window.removeEventListener("scroll", onScroll);
		};
	}, []);

	return (
		<header
			className={cn(
				"sticky top-0 z-20 overflow-hidden rounded-b-[20px] bg-brand-gradient-top text-white transition-transform duration-300 ease-out",
				// 背景(オレンジのグラデカード)と被らないよう影で分離する。
				"shadow-[0_4px_12px_rgba(77,77,77,0.25)]",
				hidden ? "-translate-y-full" : "translate-y-0",
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

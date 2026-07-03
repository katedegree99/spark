"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/utils/cn";
import { NAV_ITEMS } from "../nav-items";

/**
 * PC 左サイドバーのナビ本体(Client、Figma 準拠)。
 * アクティブ項目は文字・アイコンをブランドグラデにする(背景塗りはしない)。
 * `enabled:false` の項目は遷移させず非活性表示にする。
 */
export function SidebarNav() {
	const pathname = usePathname();

	return (
		<ul className="flex flex-col gap-2">
			{NAV_ITEMS.map(({ href, label, icon: Icon, enabled }) => {
				const active = enabled && pathname === href;
				return (
					<li key={href}>
						<Link
							href={enabled ? href : "#"}
							aria-disabled={!enabled || undefined}
							aria-current={active ? "page" : undefined}
							aria-label={label}
							tabIndex={enabled ? undefined : -1}
							className={cn(
								"flex items-center gap-5 rounded-xl p-3 font-bold text-xl tracking-wide",
								active ? "text-brand-orange" : "text-ink",
								!enabled && "pointer-events-none opacity-60",
							)}
						>
							{/* 40px 描画では viewBox 拡大で線が太く見えるため細めに指定する */}
							<Icon
								className="size-10 shrink-0"
								strokeWidth={1.5}
								aria-hidden="true"
								{...(active ? { stroke: "url(#icon-gradient)" } : {})}
							/>
							{/* lg 未満はアイコンのみ表示(ラベルは aria-label が担う) */}
							{active ? (
								<span className="hidden bg-brand-gradient bg-clip-text text-transparent lg:inline">
									{label}
								</span>
							) : (
								<span className="hidden lg:inline">{label}</span>
							)}
						</Link>
					</li>
				);
			})}
		</ul>
	);
}

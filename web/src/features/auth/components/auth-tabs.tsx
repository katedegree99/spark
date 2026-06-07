"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/utils/cn";

const TABS = [
	{ href: "/register", label: "新規登録" },
	{ href: "/login", label: "ログイン" },
] as const;

/**
 * 新規登録 / ログインの切替タブ(PC レイアウトで使用)。
 * 現在の path をアクティブ表示し、クリックで対応ページへ遷移する。
 */
export function AuthTabs() {
	const pathname = usePathname();

	return (
		<div className="flex w-full gap-2.5 rounded-xl bg-[#f3f3f3] p-1">
			{TABS.map((tab) => {
				const active = pathname === tab.href;
				return (
					<Link
						key={tab.href}
						href={tab.href}
						className={cn(
							"flex flex-1 items-center justify-center rounded-xl py-5 font-bold text-lg transition-colors",
							active ? "bg-white text-ink" : "text-secondary",
						)}
					>
						{tab.label}
					</Link>
				);
			})}
		</div>
	);
}

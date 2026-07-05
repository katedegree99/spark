"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/utils/cn";
import { NAV_ITEMS, type NavItem } from "../nav-items";

// 中央の MATCH ボタンを挟んで左右に 2 項目ずつ配置する。
// 通知(index 4)は SP では上部ヘッダのベルに集約するため下部タブには出さない。
const LEFT_ITEMS = [NAV_ITEMS[0], NAV_ITEMS[1]]; // ホーム / 検索
const RIGHT_ITEMS = [NAV_ITEMS[2], NAV_ITEMS[3]]; // 気になる / マイページ

function TabItem({ item, pathname }: { item: NavItem; pathname: string }) {
	const { href, label, tabLabel, icon: Icon, enabled } = item;
	const active = enabled && pathname === href;
	return (
		<Link
			href={enabled ? href : "#"}
			aria-disabled={!enabled || undefined}
			aria-current={active ? "page" : undefined}
			tabIndex={enabled ? undefined : -1}
			className={cn(
				"flex w-[66px] flex-col items-center gap-2 px-2 text-ink",
				!enabled && "pointer-events-none",
			)}
		>
			{/* アクティブ時のみ上端にグラデバー */}
			<span
				className={cn(
					"h-1 w-full rounded-b bg-brand-gradient",
					active ? "" : "opacity-0",
				)}
			/>
			<span className="flex flex-col items-center gap-1">
				<Icon
					className="size-8"
					strokeWidth={1}
					aria-hidden="true"
					{...(active ? { stroke: "url(#icon-gradient)" } : {})}
				/>
				{active ? (
					<span className="whitespace-nowrap bg-brand-gradient bg-clip-text font-medium text-transparent text-xs">
						{tabLabel ?? label}
					</span>
				) : (
					<span className="whitespace-nowrap font-medium text-ink text-xs">
						{tabLabel ?? label}
					</span>
				)}
			</span>
		</Link>
	);
}

/**
 * SP 下部固定タブバー(Client、Figma 準拠)。
 * 角丸の上端 + 上向きの影。中央に MATCH の円ボタンを白リングで浮かせる。
 *
 * ホームバー回避: `env(safe-area-inset-bottom)` は viewport-fit=cover を
 * 設定したときのみ非 0 になる。現状 cover は未設定のためブラウザ既定の
 * inset に任せており、env() は将来 cover を導入したときの保険として残す。
 */
export function BottomTabBar({ className }: { className?: string }) {
	const pathname = usePathname();

	return (
		<nav
			aria-label="メインナビゲーション"
			className={cn("fixed inset-x-0 bottom-0 z-20", className)}
		>
			<div className="relative mx-auto max-w-[760px]">
				{/* 1. 影だけの円(バー背面)。上端の弧のみバーからはみ出し、
				    下半分の影は不透明なバー(z-10)が覆い隠す。 */}
				<span
					aria-hidden="true"
					className="pointer-events-none absolute top-0 left-1/2 z-0 size-[83px] -translate-x-1/2 -translate-y-5 rounded-full bg-white shadow-[0px_-2px_2px_rgba(77,77,77,0.25)]"
				/>

				{/* 2. バー本体(不透明白)。 */}
				<div className="relative z-10 flex items-start justify-between rounded-t-[28px] bg-white px-5 pb-[calc(env(safe-area-inset-bottom)+20px)] shadow-[0px_-2px_2px_rgba(77,77,77,0.25)]">
					{LEFT_ITEMS.map((item) => (
						<TabItem key={item.href} item={item} pathname={pathname} />
					))}
					<div className="w-[66px] shrink-0" aria-hidden="true" />
					{RIGHT_ITEMS.map((item) => (
						<TabItem key={item.href} item={item} pathname={pathname} />
					))}
				</div>

				{/* 3. MATCH 本体(前面・影なし。影は 1 が肩代わり)。 */}
				<button
					type="button"
					aria-label="マッチ"
					className="absolute top-0 left-1/2 z-20 flex size-[83px] -translate-x-1/2 -translate-y-5 items-center justify-center rounded-full bg-white p-1"
				>
					<span
						className="flex size-full flex-col items-center justify-center rounded-full"
						// MATCH ボタン固有の角度(111.24°)のため @utility 化せずインラインで指定。
						// ストップ列は共通トークン(--brand-stops)を参照する。
						style={{
							backgroundImage: "linear-gradient(111.24deg, var(--brand-stops))",
						}}
					>
						<span className="font-[900] text-base text-white italic">
							MATCH
						</span>
					</span>
				</button>
			</div>
		</nav>
	);
}

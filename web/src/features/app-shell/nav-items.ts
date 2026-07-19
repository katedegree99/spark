import { Bell, CircleUserRound, Heart, Home, Search } from "lucide-react";
import type { ComponentType, SVGProps } from "react";

export type NavItem = {
	href: string;
	/** サイドバー用ラベル。 */
	label: string;
	/** 下部タブ用ラベル(サイドバーと文言が違う場合に指定)。 */
	tabLabel?: string;
	icon: ComponentType<SVGProps<SVGSVGElement>>;
	/** false の項目は遷移先未実装。`aria-disabled` + `pointer-events-none` にする。 */
	enabled: boolean;
};

/**
 * ナビゲーション項目の定義(PC サイドバー / SP 下部タブ共用)。
 * 並び順はサイドバー基準(ホーム → 探す → 気になる → マイページ → 通知)。
 * 実装済みはホームと探すのみ。それ以外は遷移先を作らないため `enabled: false`。
 */
export const NAV_ITEMS: NavItem[] = [
	{ href: "/home", label: "ホーム", icon: Home, enabled: true },
	{
		href: "/search",
		label: "探す",
		tabLabel: "検索",
		icon: Search,
		enabled: true,
	},
	{ href: "/likes", label: "気になる", icon: Heart, enabled: false },
	{
		href: "/mypage",
		label: "マイページ",
		icon: CircleUserRound,
		enabled: false,
	},
	{ href: "/notifications", label: "通知", icon: Bell, enabled: false },
];

import { ChevronRight } from "lucide-react";
import Link from "next/link";

type SectionHeaderProps = {
	title: string;
	/**
	 * 「もっとみる」の遷移先。未指定なら遷移先未実装としてラベルを
	 * `aria-disabled` の非活性表示にする(今回は一覧ページを作らない)。
	 */
	moreHref?: string;
};

/** セクション見出し + 右上「もっとみる >」。 */
export function SectionHeader({ title, moreHref }: SectionHeaderProps) {
	return (
		<div className="flex items-center justify-between">
			<h2 className="font-bold text-ink text-xl">{title}</h2>
			{moreHref ? (
				<Link
					href={moreHref}
					className="inline-flex items-center gap-0.5 text-secondary text-sm"
				>
					もっとみる
					<ChevronRight className="size-4" aria-hidden="true" />
				</Link>
			) : (
				// 遷移先未実装の間は視覚的な装飾として残しつつ、機能しないテキストを
				// 支援技術に読ませない(span への aria-disabled は無意味なため aria-hidden)。
				<span
					aria-hidden="true"
					className="pointer-events-none inline-flex items-center gap-0.5 text-secondary text-sm opacity-60"
				>
					もっとみる
					<ChevronRight className="size-4" />
				</span>
			)}
		</div>
	);
}

"use client";

import { useEffect, useRef, useState } from "react";
import { TagChip } from "@/components/ui/tag-chip";
import { cn } from "@/utils/cn";
import type { TagVM } from "../types";

const GAP = 8; // gap-2 (px)

/** 末尾の「...」の見た目(表示行・計測用で共有)。 */
const ELLIPSIS_CLASS =
	"shrink-0 self-center bg-brand-gradient bg-clip-text font-bold text-2xl text-transparent leading-none tracking-[3px]";

/**
 * カード共通のタグ行(Client)。ピックアップ / おすすめの両カードで使う。
 *
 * タグを 1 行に並べ、コンテナ幅に「完全に収まる」タグだけを表示する。実データは
 * ユーザーごとにタグが最大 40 個ほど付くため、入りきらないぶんは表示せず末尾に
 * 「...」を出す。タグを途中で切らないのがポイント(Figma の overflow-clip 相当。
 * カードは固定高さで縦に伸びない)。
 *
 * 各タグ・「...」の実幅を非表示の計測レイヤーで測り、収まる件数を算出する。
 */
export function TagRow({ tags }: { tags: TagVM[] }) {
	const containerRef = useRef<HTMLDivElement>(null);
	const measureRef = useRef<HTMLDivElement>(null);
	const ellipsisRef = useRef<HTMLSpanElement>(null);
	const [visibleCount, setVisibleCount] = useState(tags.length);
	// 計測前の初回ペイントでは「全タグが並び末尾チップが途中で切れた状態」が
	// 露出してちらつくため、件数が確定するまで表示行を invisible にする
	// (visibility なので行の高さは確保され、レイアウトシフトは起きない)。
	const [ready, setReady] = useState(false);

	useEffect(() => {
		const container = containerRef.current;
		const measure = measureRef.current;
		if (!container || !measure) return;
		let alive = true;

		const compute = () => {
			if (!alive) return;
			const available = container.clientWidth;
			const chips = Array.from(measure.children) as HTMLElement[];
			let used = 0;
			let count = 0;
			for (const chip of chips) {
				const add = (count === 0 ? 0 : GAP) + chip.offsetWidth;
				if (used + add > available) break;
				used += add;
				count += 1;
			}
			// 全部は入らない → 末尾の「...」ぶんを確保できるまで後ろから外す。
			if (count < tags.length) {
				const ellipsisW = ellipsisRef.current?.offsetWidth ?? 24;
				while (count > 0 && used + GAP + ellipsisW > available) {
					used -= chips[count - 1].offsetWidth + GAP;
					count -= 1;
				}
			}
			setVisibleCount(count);
			setReady(true);
		};

		compute();
		const observer = new ResizeObserver(compute);
		observer.observe(container);
		// 日本語 Web フォント適用でタグ幅が変わるため、確定後に再計測する。
		document.fonts?.ready.then(compute);
		return () => {
			alive = false;
			observer.disconnect();
		};
	}, [tags]);

	const hiddenCount = tags.length - visibleCount;

	return (
		<div className="relative overflow-hidden pt-2">
			{/* 計測用(非表示・全タグ + 「...」)。各要素の実幅を測って表示件数を決める。 */}
			<div
				aria-hidden="true"
				className="pointer-events-none invisible absolute top-0 left-0 flex gap-2"
			>
				<div ref={measureRef} className="flex gap-2">
					{tags.map((tag) => (
						<TagChip
							key={`measure-${tag.matched}-${tag.id}`}
							label={tag.name}
							variant={tag.matched ? "filled" : "outline"}
							className="shrink-0"
						/>
					))}
				</div>
				<span ref={ellipsisRef} className={ELLIPSIS_CLASS}>
					...
				</span>
			</div>
			{/* 表示行: 入りきるタグだけ + 余れば末尾に「...」。 */}
			<div
				ref={containerRef}
				className={cn("flex items-center gap-2", !ready && "invisible")}
			>
				{tags.slice(0, visibleCount).map((tag) => (
					<TagChip
						key={`${tag.matched}-${tag.id}`}
						label={tag.name}
						variant={tag.matched ? "filled" : "outline"}
						className="shrink-0"
					/>
				))}
				{hiddenCount > 0 ? (
					<span aria-hidden="true" className={ELLIPSIS_CLASS}>
						...
					</span>
				) : null}
			</div>
		</div>
	);
}

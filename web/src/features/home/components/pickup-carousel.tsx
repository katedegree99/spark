"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { cn } from "@/utils/cn";
import type { PickupCardVM } from "../types";
import { PickupCard } from "./pickup-card";

/**
 * ピックアップカードの横スクロールカルーセル(Client)。
 *
 * 中央スナップ(snap-center)の peek レイアウトで両隣のカードが覗く。スクロール位置に
 * 連動して各カードを連続的に拡縮し(中央=1.0 / 隣=0.8)、中央側の辺を transform-origin
 * で固定するので、縮小しても両隣が覗いたままになる(指の動きに滑らかに追従)。
 * snap-always で 1 スワイプ 1 カードに制限。ドットは中央に最も近いカードへ追従。
 * カードが 1 件以下のときはドットを出さない。
 */
export function PickupCarousel({ cards }: { cards: PickupCardVM[] }) {
	const [activeIndex, setActiveIndex] = useState(0);
	const scrollerRef = useRef<HTMLUListElement | null>(null);
	const slideRefs = useRef<(HTMLLIElement | null)[]>([]);
	const rafRef = useRef<number | null>(null);

	// スクロール位置から各カードのスケール(中央1.0→隣0.8)を計算して反映する。
	// 中央側の辺を origin に固定することで、縮小しても覗き量が保たれる。
	// あわせて中央に最も近いカードを activeIndex にしてドットを追従させる。
	const update = useCallback(() => {
		const scroller = scrollerRef.current;
		if (!scroller) return;
		// PC(md 以上)は Figma 準拠のプレーン横スクロール。中央拡大演出をやめ、
		// SP 幅で付いた inline transform をリセットしてから抜ける(ドットも別途 md:hidden)。
		if (window.matchMedia("(min-width: 768px)").matches) {
			slideRefs.current.forEach((el) => {
				if (!el) return;
				el.style.transform = "";
				el.style.transformOrigin = "";
			});
			return;
		}
		const center = scroller.scrollLeft + scroller.clientWidth / 2;
		let nearest = 0;
		let min = Number.POSITIVE_INFINITY;
		slideRefs.current.forEach((el, i) => {
			if (!el) return;
			const offset = el.offsetLeft + el.offsetWidth / 2 - center;
			const t = Math.min(Math.abs(offset) / el.offsetWidth, 1);
			// 下揃え: 垂直は bottom を固定して縮小(下端が中央カードと揃う)。
			el.style.transformOrigin =
				offset > 0
					? "left bottom"
					: offset < 0
						? "right bottom"
						: "center bottom";
			el.style.transform = `scale(${1 - 0.2 * t})`;
			if (Math.abs(offset) < min) {
				min = Math.abs(offset);
				nearest = i;
			}
		});
		setActiveIndex(nearest);
	}, []);

	// スクロール中は rAF で 1 フレームに集約して追従させる。
	function handleScroll() {
		if (rafRef.current != null) return;
		rafRef.current = requestAnimationFrame(() => {
			rafRef.current = null;
			update();
		});
	}

	useEffect(() => {
		update();
		window.addEventListener("resize", update);
		return () => {
			window.removeEventListener("resize", update);
			if (rafRef.current != null) cancelAnimationFrame(rafRef.current);
		};
	}, [update]);

	function scrollToSlide(index: number) {
		slideRefs.current[index]?.scrollIntoView({
			behavior: "smooth",
			inline: "center",
			block: "nearest",
		});
	}

	return (
		<section className="flex flex-col gap-3" aria-label="今日のピックアップ">
			{/* full-bleed(-mx)でスクロール領域を広げ、px でカードを中央寄せする。
			    px と gap の差が両隣カードの覗き量になる(w-[84vw] + px-[8vw])。
			    左右に白グラデを重ねて端で自然にフェードさせる(新着セクションと同様)。 */}
			<div className="relative -mx-4 md:-mx-10">
				<ul
					ref={scrollerRef}
					onScroll={handleScroll}
					className="flex items-end snap-x snap-mandatory gap-2 overflow-x-auto scroll-smooth scroll-px-[8vw] px-[8vw] pb-1 md:gap-6 md:scroll-px-10 md:px-10 [-ms-overflow-style:none] [scrollbar-width:none] [&::-webkit-scrollbar]:hidden"
				>
					{cards.map((card, index) => (
						<li
							key={card.userId}
							ref={(el) => {
								slideRefs.current[index] = el;
							}}
							className="w-[84vw] shrink-0 snap-center snap-always will-change-transform md:w-[362px] md:snap-start"
						>
							<PickupCard card={card} />
						</li>
					))}
				</ul>
				<div className="pointer-events-none absolute inset-y-0 left-0 w-6 bg-gradient-to-r from-white to-transparent" />
				<div className="pointer-events-none absolute inset-y-0 right-0 w-6 bg-gradient-to-l from-white to-transparent" />
			</div>
			{cards.length > 1 ? (
				<div className="flex items-center justify-center gap-1.5 md:hidden">
					{cards.map((card, index) => (
						<button
							key={card.userId}
							type="button"
							aria-label={`${index + 1}枚目を表示`}
							aria-current={index === activeIndex ? "true" : undefined}
							className={cn(
								"h-2 rounded-full bg-brand-gradient transition-[width,opacity] duration-300 ease-out",
								index === activeIndex ? "w-14" : "w-4 opacity-25",
							)}
							onClick={() => scrollToSlide(index)}
						/>
					))}
				</div>
			) : null}
		</section>
	);
}

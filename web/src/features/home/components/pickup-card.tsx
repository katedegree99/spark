import { Avatar } from "@/components/ui/avatar";
import type { PickupCardVM } from "../types";
import { TagRow } from "./tag-row";

/**
 * カード上段グラデに重なる装飾模様(Figma の imgCard 相当)。
 * 白のリング(右・不透明度 0.1)と同心ストライプ(左・0.07)。装飾なので
 * `aria-hidden` + `pointer-events-none`。上段の overflow-hidden でクリップされる。
 */
function PickupCardDecoration() {
	return (
		<svg
			aria-hidden="true"
			className="pointer-events-none absolute top-0 left-0 h-[352px] w-[562px]"
			viewBox="0 0 580 370"
			fill="none"
			xmlns="http://www.w3.org/2000/svg"
		>
			<path
				d="M388 2C485.202 2 564 80.7979 564 178C564 275.202 485.202 354 388 354C290.798 354 212 275.202 212 178C212 80.7979 290.798 2 388 2ZM388.438 128.09C360.632 128.09 338.09 150.631 338.09 178.438C338.09 206.244 360.632 228.786 388.438 228.786C416.245 228.786 438.786 206.244 438.786 178.438C438.786 150.631 416.245 128.09 388.438 128.09Z"
				fill="white"
				fillOpacity="0.1"
			/>
			<path
				d="M203.298 109C200.172 116.761 197.539 124.774 195.439 133H2V109H203.298ZM220.503 73C215.518 80.6399 211.064 88.6568 207.193 97H2V73H220.503ZM247.715 37C240.084 44.3983 233.068 52.4263 226.749 61H2V55.417C2 49.0053 3.21556 42.7739 5.47852 37H247.715ZM291.799 2C279.58 8.44346 268.144 16.1707 257.667 25H12.1914C16.5726 19.194 22.2285 14.3042 28.8896 10.7822C39.7978 5.01474 51.95 2.00004 64.2891 2H291.799Z"
				fill="white"
				fillOpacity="0.07"
			/>
		</svg>
	);
}

/**
 * 今日のピックアップ 1 件分のカード(presentational・Server 可、Figma 準拠)。
 *
 * カード全体がブランドグラデで、左上だけ大きな角丸(56px)。上段(グラデ地に
 * 白文字)にアバター+名前+共通タグ数、下段は白パネルで自己紹介+タグ群。
 * 上段には装飾模様(白のリング/ストライプ)を薄く重ねる。
 */
export function PickupCard({ card }: { card: PickupCardVM }) {
	return (
		<article className="flex flex-col rounded-tl-[56px] rounded-tr-[20px] rounded-br-[20px] rounded-bl-[20px] bg-brand-gradient drop-shadow-[2px_2px_3px_rgba(77,77,77,0.25)]">
			<div className="relative flex items-start gap-2.5 overflow-hidden rounded-tl-[56px] rounded-tr-[20px] p-4 text-white">
				<PickupCardDecoration />
				<Avatar
					src={card.iconUrl}
					name={card.name}
					size="md"
					className="relative z-10"
				/>
				<div className="relative z-10 flex min-w-0 flex-col justify-center gap-1">
					<p className="truncate font-bold text-lg">{card.name}</p>
					<p className="flex items-end">
						<span className="text-xs">共通のタグが</span>
						<span className="font-bold text-lg leading-none">
							{card.matchedCount}
						</span>
						<span className="text-xs">個あります</span>
					</p>
				</div>
			</div>
			<div className="flex flex-col rounded-[20px] bg-white p-4">
				<p className="line-clamp-2 h-10 text-ink text-xs leading-5">
					{card.bio}
				</p>
				{card.tags.length > 0 ? <TagRow tags={card.tags} /> : null}
			</div>
		</article>
	);
}

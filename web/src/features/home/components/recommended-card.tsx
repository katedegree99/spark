import { Avatar } from "@/components/ui/avatar";
import type { RecommendedUserVM } from "../types";
import { TagRow } from "./tag-row";

/**
 * おすすめユーザー 1 件分のカード(Server、Figma 準拠)。
 * 左にアバター+名前+自己紹介+タグ、右に共通タグ数の大きなグラデバッジ。
 */
export function RecommendedCard({ user }: { user: RecommendedUserVM }) {
	return (
		<article className="flex items-center gap-3 rounded-[20px] border border-border bg-white px-3 py-4 shadow-[2px_2px_6px_0px_rgba(77,77,77,0.25)]">
			<div className="flex min-w-0 flex-1 flex-col">
				<div className="flex items-center gap-2">
					<Avatar src={user.iconUrl} name={user.name} size="sm" />
					<div className="flex min-w-0 flex-col gap-2 py-2">
						<p className="truncate font-semibold text-base text-ink">
							{user.name}
						</p>
						{user.bio ? (
							<p className="truncate text-secondary text-xs">{user.bio}</p>
						) : null}
					</div>
				</div>
				{user.tags.length > 0 ? <TagRow tags={user.tags} /> : null}
			</div>
			<div className="flex w-[76px] shrink-0 flex-col items-center justify-center self-stretch rounded-xl bg-brand-gradient px-2 py-3 text-center text-white">
				<span className="font-din font-bold text-[40px] leading-none">
					{user.matchedCount}
				</span>
				<span className="font-medium text-base">共通</span>
			</div>
		</article>
	);
}

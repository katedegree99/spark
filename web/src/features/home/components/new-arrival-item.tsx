import { Avatar } from "@/components/ui/avatar";
import type { NewArrivalVM } from "../types";

/**
 * 新着ユーザー 1 件(黄枠アバター + 名前)。横スクロールリストの要素。
 */
export function NewArrivalItem({ user }: { user: NewArrivalVM }) {
	return (
		<li className="flex w-22 shrink-0 flex-col items-center gap-2">
			<Avatar
				src={user.iconUrl}
				name={user.name}
				size="lg"
				className="border-2 border-brand-gradient"
			/>
			<span className="w-full truncate text-center font-medium text-ink-light text-sm">
				{user.name}
			</span>
		</li>
	);
}

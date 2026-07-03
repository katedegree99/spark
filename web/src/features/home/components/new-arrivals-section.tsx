import { ScrollEdgeFade } from "@/components/ui/scroll-edge-fade";
import { SectionHeader } from "@/components/ui/section-header";
import type { NewArrivalVM } from "../types";
import { NewArrivalItem } from "./new-arrival-item";

/**
 * 新着ユーザーセクション(Server)。
 * 丸アバター + 名前を横スクロールで並べる。0 件のときは非表示にする。
 */
export function NewArrivalsSection({ users }: { users: NewArrivalVM[] }) {
	if (users.length === 0) return null;

	return (
		<section className="flex flex-col gap-3">
			<SectionHeader title="新着" />
			{/* スクロール領域はコンテナ余白を打ち消して画面端まで全幅化し、
			    左右に白グラデを重ねて端で自然にフェードさせる。 */}
			<div className="relative -mx-4 md:-mx-10">
				<ul className="scrollbar-none flex gap-1 overflow-x-auto px-4 pb-1 md:px-10">
					{users.map((user) => (
						<NewArrivalItem key={user.userId} user={user} />
					))}
				</ul>
				<ScrollEdgeFade />
			</div>
		</section>
	);
}

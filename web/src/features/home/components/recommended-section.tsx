import { SectionHeader } from "@/components/ui/section-header";
import type { RecommendedUserVM } from "../types";
import { RecommendedCard } from "./recommended-card";

/**
 * おすすめユーザーセクション(Server)。
 * SP は 1 カラム、PC は 2 カラムのグリッドでカードを並べる。0 件なら非表示。
 */
export function RecommendedSection({ users }: { users: RecommendedUserVM[] }) {
	if (users.length === 0) return null;

	return (
		<section className="flex flex-col gap-3">
			<SectionHeader title="おすすめのユーザー" />
			<ul className="grid grid-cols-1 gap-4 md:grid-cols-2">
				{users.map((user) => (
					<li key={user.userId}>
						<RecommendedCard user={user} />
					</li>
				))}
			</ul>
		</section>
	);
}

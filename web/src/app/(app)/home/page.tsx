import { NewArrivalsSection } from "@/features/home/components/new-arrivals-section";
import { PickupSection } from "@/features/home/components/pickup-section";
import { RecommendedSection } from "@/features/home/components/recommended-section";
import {
	getNewArrivals,
	getPickupUsers,
	getRecommendedUsers,
} from "@/features/home/data";

/**
 * 認証後ホーム (`/home`)。Server Component。
 *
 * 認証ガードとアプリシェルは `(app)/layout.tsx` が担う。ここでは 3 セクション
 * (今日のピックアップ / 新着 / おすすめ)のデータを並行取得して縦に並べる。
 * 各セクションはセクション単位で degrade するため、1 つの失敗で画面全体は落ちない。
 */
export default async function HomePage() {
	const [pickup, newArrivals, recommended] = await Promise.all([
		getPickupUsers(),
		getNewArrivals(),
		getRecommendedUsers(),
	]);

	return (
		<main className="flex-1">
			{/* PC のみのページ見出しバー(SP は MobileHeader が担う) */}
			<div className="hidden border-border border-b px-5 py-5 md:block">
				<h1 className="font-bold text-2xl text-ink">ホーム</h1>
			</div>
			<div className="flex w-full flex-col gap-8 px-4 py-6 pb-24 md:px-10 md:pb-10">
				<PickupSection data={pickup} />
				<NewArrivalsSection users={newArrivals} />
				<RecommendedSection users={recommended} />
			</div>
		</main>
	);
}

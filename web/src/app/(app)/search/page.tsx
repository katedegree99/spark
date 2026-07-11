import { RecommendedCard } from "@/features/home/components/recommended-card";
import { SelectedTagList } from "@/features/search/components/selected-tag-list";
import { TagSearchBox } from "@/features/search/components/tag-search-box";
import { searchUsers } from "@/features/search/data";
import { parseTagsParam } from "@/features/search/search-params";

/**
 * 探すページ (`/search`)。Server Component。
 *
 * URL の `?tags=<id>:<name>`(繰り返し)が唯一の真実。ここでパースして
 * `GET /users?tagIds=` で AND 絞り込みし、ホームと同じ `RecommendedCard` で
 * 一覧表示する。タグの選択・削除は `TagSearchBox` が URL を更新して再フェッチ。
 *
 * ヘッダは 1 インスタンスをレスポンシブ切替(SP: sticky グラデ / PC: 白地 + border-b)。
 * SP は独自ヘッダを持つため共通 `MobileHeader` は `/search` で非表示になる。
 */
export default async function SearchPage({
	searchParams,
}: {
	searchParams: Promise<{ [key: string]: string | string[] | undefined }>;
}) {
	const { tags } = await searchParams;
	const selectedTags = parseTagsParam(tags);
	const result = await searchUsers(selectedTags.map((t) => t.id));

	return (
		// group: TagSearchBox が遷移中に立てる data-pending を結果領域の減光に使う
		<main className="group flex flex-1 flex-col">
			<div className="sticky top-0 z-20 bg-white md:static md:border-border md:border-b">
				<div className="flex flex-col md:flex-row md:items-center md:justify-between md:p-5">
					{/* SP: グラデ帯は検索入力の縦中央(≈25px)で終わり、入力の下半分は
					    白地にはみ出す(Figma: bg h156 に対し入力が 69〜119px で跨る)。
					    帯の pb-[45px] = 見出しとの間隔 20px + 入力の上半分 25px。 */}
					<div className="relative rounded-b-[20px] bg-brand-gradient-top px-5 pt-5 pb-[45px] text-white md:rounded-none md:bg-none md:p-0 md:pb-0 md:text-ink">
						{/* Figma: SP はグラデの上に白 20% を重ねて淡いトーンにする(MobileHeader と同様) */}
						<div
							className="absolute inset-0 rounded-b-[20px] bg-white/20 md:hidden"
							aria-hidden="true"
						/>
						<h1 className="relative z-10 font-bold text-2xl">探す</h1>
					</div>
					<div className="-mt-[25px] relative z-10 px-5 md:mt-0 md:px-0">
						<TagSearchBox selectedTags={selectedTags} />
						{/* SP: チップは検索入力の直下 */}
						<SelectedTagList
							selectedTags={selectedTags}
							className="mt-3 md:hidden"
						/>
					</div>
				</div>
				{/* PC: チップは見出し行の下に全幅で並べる(入力の右カラム内で
				    折り返すとヘッダが縦に膨らみ「探す」が浮くため) */}
				<SelectedTagList
					selectedTags={selectedTags}
					className="hidden px-5 pb-6 md:flex"
				/>
			</div>
			<div className="flex-1 px-4 py-6 pb-24 transition-opacity md:bg-[#fafafa] md:px-10 md:pb-10 group-has-[[data-pending]]:opacity-60">
				{!Array.isArray(result) ? (
					<p className="py-10 text-center text-secondary text-sm">
						ユーザーを取得できませんでした。時間をおいて再度お試しください。
					</p>
				) : result.length === 0 ? (
					<p className="py-10 text-center text-secondary text-sm">
						{selectedTags.length > 0
							? "選択したタグに一致するユーザーが見つかりませんでした"
							: "ユーザーが見つかりませんでした"}
					</p>
				) : (
					<div className="grid grid-cols-1 gap-5 md:grid-cols-2 md:gap-4">
						{result.map((user) => (
							<RecommendedCard key={user.userId} user={user} />
						))}
					</div>
				)}
			</div>
		</main>
	);
}

import type { PickupError } from "../data";
import type { PickupCardVM } from "../types";
import { PickupCarousel } from "./pickup-carousel";

/**
 * 今日のピックアップセクション(Server)。
 * 取得エラー / 0 件 / 正常 を出し分け、セクション単位で degrade する
 * (ホーム全体やほかのセクションには影響させない)。
 */
export function PickupSection({
	data,
}: {
	data: PickupCardVM[] | PickupError;
}) {
	return (
		<section className="flex flex-col gap-3">
			<h2 className="font-bold text-ink text-xl">今日のピックアップ</h2>
			{!Array.isArray(data) ? (
				<p className="text-secondary text-sm">
					ピックアップを取得できませんでした。時間をおいて再度お試しください。
				</p>
			) : data.length === 0 ? (
				<p className="text-secondary text-sm">
					今日のピックアップはまだありません。
				</p>
			) : (
				<PickupCarousel cards={data} />
			)}
		</section>
	);
}

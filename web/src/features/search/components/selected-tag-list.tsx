"use client";

import { cn } from "@/utils/cn";
import type { SelectedTag } from "../types";
import { usePushTags } from "../use-push-tags";
import { SelectedTagChip } from "./selected-tag-chip";

/**
 * 選択中タグのチップ一覧(Client)。× で削除して URL を更新する。
 * SP は検索入力の直下、PC は見出し行の下の全幅の段、と置き場所が違うため
 * `TagSearchBox` から独立させ、page が表示位置ごとに配置する。0 個なら描画しない。
 */
export function SelectedTagList({
	selectedTags,
	className,
}: {
	selectedTags: SelectedTag[];
	className?: string;
}) {
	const { isPending, pushTags } = usePushTags();

	if (selectedTags.length === 0) return null;

	return (
		<div
			data-pending={isPending || undefined}
			className={cn("flex flex-wrap gap-2", className)}
		>
			{selectedTags.map((tag) => (
				<SelectedTagChip
					key={tag.id}
					name={tag.name}
					// 遷移中の連打を無視する。props の selectedTags は RSC 応答まで古いままの
					// ため、続けて押すと 1 つ前の削除を巻き戻した URL を push してしまう。
					onRemove={() => {
						if (isPending) return;
						pushTags(selectedTags.filter((t) => t.id !== tag.id));
					}}
				/>
			))}
		</div>
	);
}

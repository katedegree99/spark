"use client";

import { useRouter } from "next/navigation";
import { useTransition } from "react";
import { serializeTagsParam } from "./search-params";
import type { SelectedTag } from "./types";

/**
 * 選択タグ配列で `/search` の URL を更新する hook。
 * URL が唯一の真実なので、タグの選択・削除はすべてこの push で行う
 * (0 個なら `/search` に戻す)。`isPending` は遷移中の減光表示に使う。
 */
export function usePushTags(): {
	isPending: boolean;
	pushTags: (next: SelectedTag[]) => void;
} {
	const router = useRouter();
	const [isPending, startTransition] = useTransition();

	function pushTags(next: SelectedTag[]) {
		const queryString = serializeTagsParam(next);
		startTransition(() => {
			router.push(queryString ? `/search?${queryString}` : "/search", {
				scroll: false,
			});
		});
	}

	return { isPending, pushTags };
}

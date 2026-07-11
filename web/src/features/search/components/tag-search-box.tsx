"use client";

import { Search } from "lucide-react";
import { useRouter } from "next/navigation";
import { useEffect, useRef, useState, useTransition } from "react";
import { Input } from "@/components/ui/input";
import { serializeTagsParam } from "../search-params";
import type { SelectedTag } from "../types";
import { useTagSuggestions } from "../use-tag-suggestions";
import { SelectedTagChip } from "./selected-tag-chip";

/**
 * タグ検索ボックス(Client)。入力 + サジェストドロップダウン + 選択済みチップ行。
 *
 * URL(`?tags=`)が唯一の真実。選択・削除はローカル state を持たず
 * `router.push` で URL を更新し、SC の再フェッチに任せる。ローカル state は
 * 「入力テキスト」「ドロップダウン開閉」のみ。遷移中は root の `data-pending`
 * 属性を立て、親(page)が CSS で結果領域を減光する。
 */
export function TagSearchBox({
	selectedTags,
}: {
	selectedTags: SelectedTag[];
}) {
	const router = useRouter();
	const [isPending, startTransition] = useTransition();
	const [query, setQuery] = useState("");
	const [open, setOpen] = useState(false);
	const rootRef = useRef<HTMLDivElement>(null);

	const { suggestions } = useTagSuggestions(query);
	// 選択済みタグはサジェスト候補から除外して表示する。
	const selectedIds = new Set(selectedTags.map((t) => t.id));
	const visibleSuggestions = suggestions.filter((s) => !selectedIds.has(s.id));

	// ドロップダウンの外側クリックで閉じる。
	useEffect(() => {
		function onPointerDown(e: PointerEvent) {
			if (rootRef.current && !rootRef.current.contains(e.target as Node)) {
				setOpen(false);
			}
		}
		document.addEventListener("pointerdown", onPointerDown);
		return () => document.removeEventListener("pointerdown", onPointerDown);
	}, []);

	/** 選択タグ配列で URL を更新する(0 個なら `/search` に戻す)。 */
	function pushTags(next: SelectedTag[]) {
		const queryString = serializeTagsParam(next);
		startTransition(() => {
			router.push(queryString ? `/search?${queryString}` : "/search", {
				scroll: false,
			});
		});
	}

	function selectTag(tag: SelectedTag) {
		pushTags([...selectedTags, tag]);
		setQuery("");
		setOpen(false);
	}

	function removeTag(id: number) {
		pushTags(selectedTags.filter((t) => t.id !== id));
	}

	function onKeyDown(e: React.KeyboardEvent<HTMLInputElement>) {
		if (e.key === "Escape") {
			setOpen(false);
			return;
		}
		if (e.key === "Enter") {
			// IME 変換確定の Enter で誤選択しないようガードする。
			if (e.nativeEvent.isComposing) return;
			e.preventDefault();
			const first = visibleSuggestions[0];
			if (open && first) selectTag(first);
		}
	}

	return (
		<div
			ref={rootRef}
			data-pending={isPending || undefined}
			className="flex w-full flex-col gap-3 md:w-[500px]"
		>
			<div className="relative">
				<Input
					icon={Search}
					placeholder="タグで検索"
					value={query}
					onChange={(e) => {
						setQuery(e.target.value);
						setOpen(e.target.value.trim() !== "");
					}}
					onFocus={() => setOpen(query.trim() !== "")}
					onKeyDown={onKeyDown}
					aria-label="タグで検索"
					// SP: 白地 + brand-yellow 枠(Input 既定)/ PC: グレー地 #f3f3f3・枠なし(Figma)
					wrapperClassName="rounded-lg p-3 md:border-transparent md:bg-[#f3f3f3] md:focus-within:border-transparent"
				/>
				{open && query.trim() !== "" && (
					<ul className="absolute top-full right-0 left-0 z-30 mt-2 max-h-64 overflow-y-auto rounded-xl border border-border bg-white py-1 shadow-[2px_2px_6px_0px_rgba(77,77,77,0.25)]">
						{visibleSuggestions.length === 0 ? (
							<li className="px-4 py-3 text-secondary text-sm">
								該当するタグがありません
							</li>
						) : (
							visibleSuggestions.map((tag) => (
								<li key={tag.id}>
									<button
										type="button"
										onClick={() => selectTag(tag)}
										className="w-full px-4 py-3 text-left text-base text-ink transition-colors hover:bg-[#f3f3f3]"
									>
										{tag.name}
									</button>
								</li>
							))
						)}
					</ul>
				)}
			</div>
			{selectedTags.length > 0 && (
				<div className="flex flex-wrap gap-2">
					{selectedTags.map((tag) => (
						<SelectedTagChip
							key={tag.id}
							name={tag.name}
							onRemove={() => removeTag(tag.id)}
						/>
					))}
				</div>
			)}
		</div>
	);
}

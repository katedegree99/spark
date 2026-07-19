"use client";

import { Search } from "lucide-react";
import { useEffect, useId, useRef, useState } from "react";
import { Input } from "@/components/ui/input";
import { cn } from "@/utils/cn";
import type { SelectedTag } from "../types";
import { usePushTags } from "../use-push-tags";
import { useTagSuggestions } from "../use-tag-suggestions";

/**
 * タグ検索ボックス(Client)。入力 + サジェストドロップダウン(combobox)。
 *
 * URL(`?tags=`)が唯一の真実。選択はローカル state を持たず `usePushTags` で
 * URL を更新し、SC の再フェッチに任せる。ローカル state は「入力テキスト」
 * 「ドロップダウン開閉」「アクティブ候補」のみ。遷移中は root の `data-pending`
 * 属性を立て、親(page)が CSS で結果領域を減光する。選択済みチップの表示・削除は
 * `SelectedTagList` が担う(表示位置が SP/PC で異なるため分離)。
 *
 * キーボード操作: ↑↓ で候補移動、Enter でアクティブ候補を確定(IME 変換中は無視)、
 * Escape で閉じる。候補が現入力に対応するまで(debounce + fetch 完了前)は
 * 「検索中…」を表示し、Enter でのステイル候補の誤選択を防ぐ。
 */
export function TagSearchBox({
	selectedTags,
}: {
	selectedTags: SelectedTag[];
}) {
	const { isPending, pushTags } = usePushTags();
	const [query, setQuery] = useState("");
	const [open, setOpen] = useState(false);
	const [activeIndex, setActiveIndex] = useState(0);
	const rootRef = useRef<HTMLDivElement>(null);
	const listboxId = useId();

	const { suggestions, isLoading, debouncedQuery } = useTagSuggestions(query);
	const trimmed = query.trim();
	// 候補が現在の入力に対応しているか(debounce 待ち・fetch 中は false)。
	const isSettled = !isLoading && debouncedQuery === trimmed;
	// 選択済みタグはサジェスト候補から除外して表示する。
	const selectedIds = new Set(selectedTags.map((t) => t.id));
	const visibleSuggestions = suggestions.filter((s) => !selectedIds.has(s.id));

	const dropdownVisible = open && trimmed !== "";
	const activeSuggestion = isSettled ? visibleSuggestions[activeIndex] : null;
	const activeOptionId = activeSuggestion
		? `${listboxId}-opt-${activeSuggestion.id}`
		: undefined;

	// 候補セットが変わったらアクティブ候補を先頭に戻す。
	// biome-ignore lint/correctness/useExhaustiveDependencies: 候補の切り替わり(debouncedQuery)にだけ反応させる
	useEffect(() => {
		setActiveIndex(0);
	}, [debouncedQuery]);

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

	function selectTag(tag: SelectedTag) {
		// 遷移中の連続選択を無視する。props の selectedTags は RSC 応答まで古いままの
		// ため、続けて選ぶと 1 つ前の選択を巻き戻した URL を push してしまう。
		if (isPending) return;
		pushTags([...selectedTags, tag]);
		setQuery("");
		setOpen(false);
		// debounce 内の同一 query 再入力ではリセット effect が発火しないため、
		// 選択時にもアクティブ候補を先頭へ戻す(範囲外 index の残留防止)。
		setActiveIndex(0);
	}

	function onKeyDown(e: React.KeyboardEvent<HTMLInputElement>) {
		if (e.key === "Escape") {
			setOpen(false);
			return;
		}
		if (!dropdownVisible) return;
		if (e.key === "ArrowDown" || e.key === "ArrowUp") {
			e.preventDefault();
			// ロード中は「検索中…」しか見えておらず、不可視のステイル候補を
			// 操作して範囲外 index を残さないよう無視する。
			if (!isSettled || visibleSuggestions.length === 0) return;
			const delta = e.key === "ArrowDown" ? 1 : -1;
			const next =
				(activeIndex + delta + visibleSuggestions.length) %
				visibleSuggestions.length;
			setActiveIndex(next);
			// アクティブ候補が overflow スクロールの可視領域外に出たら追従させる。
			document
				.getElementById(`${listboxId}-opt-${visibleSuggestions[next].id}`)
				?.scrollIntoView({ block: "nearest" });
			return;
		}
		if (e.key === "Enter") {
			// IME 変換確定の Enter で誤選択しないようガードする。
			if (e.nativeEvent.isComposing) return;
			e.preventDefault();
			// 候補が現入力に追いつくまでは確定しない(前 query のステイル候補の誤選択防止)。
			if (activeSuggestion) selectTag(activeSuggestion);
		}
	}

	return (
		<div
			ref={rootRef}
			data-pending={isPending || undefined}
			className="relative w-full md:w-[500px]"
		>
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
				role="combobox"
				aria-label="タグで検索"
				aria-expanded={dropdownVisible}
				aria-controls={dropdownVisible ? listboxId : undefined}
				aria-activedescendant={activeOptionId}
				aria-autocomplete="list"
				// SP: 白地 + brand-yellow 枠(Input 既定)/ PC: グレー地 #f3f3f3・枠なし(Figma)
				wrapperClassName="rounded-lg p-3 md:border-transparent md:bg-[#f3f3f3] md:focus-within:border-transparent"
			/>
			{dropdownVisible && (
				<ul
					id={listboxId}
					// biome-ignore lint/a11y/noNoninteractiveElementToInteractiveRole: combobox パターンの候補リスト。操作は input の aria-activedescendant 経由で行う
					role="listbox"
					aria-label="タグ候補"
					className="absolute top-full right-0 left-0 z-30 mt-2 max-h-64 overflow-y-auto rounded-xl border border-border bg-white py-1 shadow-[2px_2px_6px_0px_rgba(77,77,77,0.25)]"
				>
					{!isSettled ? (
						<li className="px-4 py-3 text-secondary text-sm" aria-hidden="true">
							検索中…
						</li>
					) : visibleSuggestions.length === 0 ? (
						<li className="px-4 py-3 text-secondary text-sm" aria-hidden="true">
							該当するタグがありません
						</li>
					) : (
						visibleSuggestions.map((tag, index) => (
							<li key={tag.id}>
								<button
									type="button"
									id={`${listboxId}-opt-${tag.id}`}
									role="option"
									aria-selected={index === activeIndex}
									onClick={() => selectTag(tag)}
									onPointerEnter={() => setActiveIndex(index)}
									className={cn(
										"w-full px-4 py-3 text-left text-base text-ink transition-colors",
										// Enter で確定される候補(アクティブ)を視覚的に予告する。
										index === activeIndex && "bg-[#f3f3f3]",
									)}
								>
									{tag.name}
								</button>
							</li>
						))
					)}
				</ul>
			)}
		</div>
	);
}

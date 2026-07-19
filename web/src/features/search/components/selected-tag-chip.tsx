"use client";

import { X } from "lucide-react";

/**
 * 選択済みタグの削除可能チップ(Client)。
 * 表示は `TagChip` の filled 風(グラデ塗り + 白文字)だが、削除ボタンを持つため
 * display-only の既存 `TagChip` は改変せず別コンポーネントにする。
 */
export function SelectedTagChip({
	name,
	onRemove,
}: {
	name: string;
	onRemove: () => void;
}) {
	return (
		<span className="inline-flex items-center gap-1 rounded-full border bg-brand-gradient-tag-fill py-2 pr-2 pl-3 font-semibold text-white text-xs leading-none">
			{name}
			<button
				type="button"
				aria-label={`${name} を外す`}
				onClick={onRemove}
				className="flex items-center justify-center rounded-full transition-opacity hover:opacity-70"
			>
				<X className="size-3.5" strokeWidth={2.5} aria-hidden="true" />
			</button>
		</span>
	);
}

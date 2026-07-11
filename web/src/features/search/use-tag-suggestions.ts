"use client";

import { useEffect, useState } from "react";
import useSWR from "swr";
import type { ThingsResponse } from "@/lib/api/generated/model";
import type { SelectedTag } from "./types";

/**
 * `/api/things` プロキシ経由でタグ候補を取得し `SelectedTag[]` にする。
 * id/name 欠落要素は除外。失敗は throw して SWR のエラー扱いにする
 * (握りつぶして空配列を返すと「成功 = 0 件」としてキャッシュされ、
 * リトライも効かなくなるため)。UI 側ではエラーは候補なしとして黙って扱う。
 */
async function fetchSuggestions([path, q]: [string, string]): Promise<
	SelectedTag[]
> {
	const res = await fetch(`${path}?q=${encodeURIComponent(q)}`);
	if (!res.ok) throw new Error(`tag suggestions failed: ${res.status}`);
	const body = (await res.json()) as ThingsResponse;
	return (body.things ?? [])
		.filter((t): t is { id: number; name: string } => {
			return t.id != null && t.name != null;
		})
		.map((t) => ({ id: t.id, name: t.name }));
}

/**
 * 入力テキストからタグ候補を返すフック。
 * 300ms の debounce 後に `/api/things?q=` を叩く。空 query は fetch しない
 * (SWR の null key)。`keepPreviousData` で入力中のちらつきを抑える。
 *
 * `debouncedQuery` は現在の候補がどの query に対するものかの判定用。
 * `suggestions` は keepPreviousData により前 query の結果でありうるため、
 * Enter 確定などは「`debouncedQuery` が現入力と一致 && !isLoading」を
 * 確認してから行うこと(ステイル候補の誤選択防止)。
 */
export function useTagSuggestions(query: string): {
	suggestions: SelectedTag[];
	isLoading: boolean;
	debouncedQuery: string;
} {
	const [debouncedQuery, setDebouncedQuery] = useState("");

	useEffect(() => {
		const timer = setTimeout(() => setDebouncedQuery(query.trim()), 300);
		return () => clearTimeout(timer);
	}, [query]);

	const { data, isLoading } = useSWR(
		debouncedQuery !== "" ? (["/api/things", debouncedQuery] as const) : null,
		fetchSuggestions,
		{ keepPreviousData: true },
	);

	return { suggestions: data ?? [], isLoading, debouncedQuery };
}

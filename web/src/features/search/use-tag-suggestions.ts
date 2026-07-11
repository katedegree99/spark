"use client";

import { useEffect, useState } from "react";
import useSWR from "swr";
import type { ThingsResponse } from "@/lib/api/generated/model";
import type { SelectedTag } from "./types";

/**
 * `/api/things` プロキシ経由でタグ候補を取得し `SelectedTag[]` にする。
 * id/name 欠落要素は除外。fetch エラーはサイレント(候補なし扱い)。
 */
async function fetchSuggestions([path, q]: [string, string]): Promise<
	SelectedTag[]
> {
	try {
		const res = await fetch(`${path}?q=${encodeURIComponent(q)}`);
		if (!res.ok) return [];
		const body = (await res.json()) as ThingsResponse;
		return (body.things ?? [])
			.filter((t): t is { id: number; name: string } => {
				return t.id != null && t.name != null;
			})
			.map((t) => ({ id: t.id, name: t.name }));
	} catch {
		return [];
	}
}

/**
 * 入力テキストからタグ候補を返すフック。
 * 300ms の debounce 後に `/api/things?q=` を叩く。空 query は fetch しない
 * (SWR の null key)。`keepPreviousData` で入力中のちらつきを抑える。
 */
export function useTagSuggestions(query: string): {
	suggestions: SelectedTag[];
	isLoading: boolean;
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

	return { suggestions: data ?? [], isLoading };
}

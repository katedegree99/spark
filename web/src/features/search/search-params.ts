/**
 * `/search?tags=<id>:<name>` パラメータの parse / serialize(純関数)。
 *
 * `/things` に ID → 名前の逆引き API が無いため、リロード後もチップ名を
 * 復元できるよう name を URL に持たせる。SC・クライアント両用なので
 * `server-only` は付けない。
 */
import type { SelectedTag } from "./types";

/**
 * `?tags=` の値(単数 / 複数 / 未指定)を選択タグ配列にパースする。
 * 最初の `:` で分割し、`id が正の整数 && name 非空` を満たさない
 * エントリは黙って捨てる。id で重複排除する。
 */
export function parseTagsParam(
	value: string | string[] | undefined,
): SelectedTag[] {
	const entries = value == null ? [] : Array.isArray(value) ? value : [value];
	const seen = new Set<number>();
	const tags: SelectedTag[] = [];

	for (const entry of entries) {
		const sep = entry.indexOf(":");
		if (sep === -1) continue;
		const id = Number(entry.slice(0, sep));
		const name = entry.slice(sep + 1);
		if (!Number.isInteger(id) || id <= 0 || name === "" || seen.has(id)) {
			continue;
		}
		seen.add(id);
		tags.push({ id, name });
	}
	return tags;
}

/** 選択タグ配列を `tags=1:xxx&tags=2:yyy` 形式のクエリ文字列にする。0 個なら空文字。 */
export function serializeTagsParam(tags: SelectedTag[]): string {
	const params = new URLSearchParams();
	for (const tag of tags) {
		params.append("tags", `${tag.id}:${tag.name}`);
	}
	return params.toString();
}

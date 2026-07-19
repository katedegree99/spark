/**
 * 探すページで共有する型。SC・クライアント両方から import する。
 */

/** URL(`?tags=`)とサジェスト・チップ間で受け渡す選択中タグ。 */
export type SelectedTag = {
	id: number;
	name: string;
};

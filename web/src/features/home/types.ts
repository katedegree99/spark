/**
 * ホーム画面 UI 用のローカル View Model。
 *
 * 生成型(`PickupUserResponse` 等)は全フィールドが optional で UI に持ち込むと
 * 取り回しが悪いため、`data.ts` のマッパーでこの非 optional な VM に変換し、
 * コンポーネントは VM のみを受け取る。
 */

export type TagVM = {
	id: number;
	name: string;
	/** ログインユーザーと共通のタグなら true。 */
	matched: boolean;
};

export type PickupCardVM = {
	userId: number;
	name: string;
	bio: string;
	iconUrl: string | null;
	/** 共通タグ数(「共通のタグが N 個」)。 */
	matchedCount: number;
	tags: TagVM[];
};

export type NewArrivalVM = {
	userId: number;
	name: string;
	iconUrl: string | null;
};

export type RecommendedUserVM = {
	userId: number;
	name: string;
	bio: string;
	iconUrl: string | null;
	matchedCount: number;
	tags: TagVM[];
};

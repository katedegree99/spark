import { z } from "zod";

/**
 * プロフィール設定フォームのバリデーションスキーマ (zod v4)。
 *
 * OpenAPI (`schema/openapi/openapi.yaml`) の `ProfileCreateRequest` に準拠する。
 * `doingThingIds` / `wantThingIds` は API 上は ID 配列だが、フォーム内部では
 * 表示用に `{ id, name }` のタグとして保持し、送信時に ID へ写像する。
 * エラーメッセージは日本語。
 */

/** やっていること / やってみたいこと の 1 タグ(thing)。 */
export const thingTagSchema = z.object({
	id: z.number(),
	name: z.string(),
});

export const profileSchema = z.object({
	name: z
		.string()
		.min(1, { error: "入力してください" })
		.max(100, { error: "100文字以内で入力してください" }),
	bio: z
		.string()
		.min(1, { error: "入力してください" })
		.max(1000, { error: "1000文字以内で入力してください" }),
	doings: z
		.array(thingTagSchema)
		.max(20, { error: "最大20個まで選択できます" }),
	wants: z.array(thingTagSchema).max(20, { error: "最大20個まで選択できます" }),
});

/** やっていること / やってみたいこと に設定できるタグの上限。 */
export const MAX_THINGS = 20;

export type ThingTag = z.infer<typeof thingTagSchema>;
export type ProfileInput = z.infer<typeof profileSchema>;

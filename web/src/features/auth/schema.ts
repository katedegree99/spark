import { z } from "zod";

/**
 * 認証フォームの入力バリデーションスキーマ (zod v4)。
 *
 * OpenAPI (`schema/openapi/openapi.yaml`) の auth リクエストに準拠する。
 * `@hookform/resolvers/zod` の `zodResolver` でクライアント検証に使う前提。
 * エラーメッセージは日本語。
 */

/** ログインフォーム: email(メール形式) + password(必須)。 */
export const loginSchema = z.object({
	email: z.email({ error: "形式が正しくありません" }),
	password: z.string().min(1, { error: "入力してください" }),
});

/**
 * 新規登録フォーム: email + password(8文字以上) + confirmPassword(確認用)。
 * confirmPassword はクライアント専用で API には送らない。
 * password と confirmPassword の一致を検証し、エラーは confirmPassword 側に出す。
 */
export const registerSchema = z
	.object({
		email: z.email({ error: "形式が正しくありません" }),
		password: z.string().min(8, { error: "8文字以上で入力してください" }),
		confirmPassword: z.string().min(1, { error: "入力してください" }),
	})
	.refine((data) => data.password === data.confirmPassword, {
		error: "パスワードが一致しません",
		path: ["confirmPassword"],
	});

/** OTP 検証フォーム: code(6桁の数字)。email はフロー上別途保持する。 */
export const otpSchema = z.object({
	code: z.string().regex(/^\d{6}$/, { error: "6桁の数字を入力してください" }),
});

export type LoginInput = z.infer<typeof loginSchema>;
export type RegisterInput = z.infer<typeof registerSchema>;
export type OtpInput = z.infer<typeof otpSchema>;

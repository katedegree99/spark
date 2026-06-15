"use server";

import { authErrorMessage } from "@/features/auth/errors";
import {
	type LoginInput,
	loginSchema,
	otpSchema,
	type RegisterInput,
	registerSchema,
} from "@/features/auth/schema";
import { apiFetch } from "@/lib/api/fetcher";
import type {
	AuthTokensResponse,
	GoogleLoginRequest,
	LoginRequest,
	OtpSentResponse,
	OtpVerifyRequest,
	RegisterRequest,
} from "@/lib/api/generated/model";
import { setAuthCookies } from "@/lib/auth/session";

/** OTP 送信系(register/login)の戻り値。成功時は次ステップ用の宛先メールを返す。 */
export type AuthActionResult =
	| { ok: true; email: string }
	| { ok: false; message: string };

/** OTP 検証の戻り値。成功時はトークンを Cookie に保存済みなので追加データは不要。 */
export type VerifyOtpActionResult =
	| { ok: true }
	| { ok: false; message: string };

/**
 * メールアドレスで新規登録する Server Action。
 *
 * - クライアント検証はバイパスされうるため、サーバ側で再度 Zod 検証する。
 * - `confirmPassword` はクライアント専用なので API へは送らない。
 * - 成功時はサーバが確認用 OTP をメール送信する(トークンはまだ返らない)。
 */
export async function registerAction(
	input: RegisterInput,
): Promise<AuthActionResult> {
	const parsed = registerSchema.safeParse(input);
	if (!parsed.success) {
		return { ok: false, message: "入力内容を確認してください" };
	}

	const { email, password } = parsed.data;
	const payload: RegisterRequest = { email, password };

	const res = await apiFetch<OtpSentResponse>("/auth/register", {
		method: "POST",
		body: JSON.stringify(payload),
	});

	if (res.ok) {
		return { ok: true, email };
	}
	return {
		ok: false,
		message: authErrorMessage(
			res.error,
			"登録に失敗しました。時間をおいて再度お試しください",
		),
	};
}

/**
 * メールアドレス + パスワードでログインする Server Action。
 * 成功時もサーバは OTP を送るため、宛先メールを返して `/otp` へ進ませる。
 */
export async function loginAction(
	input: LoginInput,
): Promise<AuthActionResult> {
	const parsed = loginSchema.safeParse(input);
	if (!parsed.success) {
		return { ok: false, message: "入力内容を確認してください" };
	}

	const payload: LoginRequest = parsed.data;

	const res = await apiFetch<OtpSentResponse>("/auth/login", {
		method: "POST",
		body: JSON.stringify(payload),
	});

	if (res.ok) {
		return { ok: true, email: parsed.data.email };
	}
	return {
		ok: false,
		message: authErrorMessage(
			res.error,
			"ログインに失敗しました。時間をおいて再度お試しください",
		),
	};
}

/**
 * OTP を検証する Server Action。
 * 成功時に `AuthTokensResponse` のトークンを **httpOnly Cookie に保存**する
 * (ブラウザ JS にトークンを露出させない)。宛先メールはフロー上から渡す。
 */
export async function verifyOtpAction(input: {
	email: string;
	code: string;
}): Promise<VerifyOtpActionResult> {
	const parsed = otpSchema.safeParse({ code: input.code });
	if (!parsed.success || !input.email) {
		return { ok: false, message: "認証コードを確認してください" };
	}

	const payload: OtpVerifyRequest = {
		email: input.email,
		code: parsed.data.code,
	};

	const res = await apiFetch<AuthTokensResponse>("/auth/otp/verify", {
		method: "POST",
		body: JSON.stringify(payload),
	});

	if (res.ok) {
		await setAuthCookies(res.data);
		return { ok: true };
	}
	return {
		ok: false,
		message: authErrorMessage(
			res.error,
			"認証に失敗しました。時間をおいて再度お試しください",
		),
	};
}

/**
 * Google の ID トークンでログイン/連携する Server Action。
 *
 * フロントが GIS で取得した `id_token` を検証のため backend へ送る。
 * Email フローと違い OTP は無く、成功時に直接 `AuthTokensResponse` が返るので
 * トークンを httpOnly Cookie に保存する。
 */
export async function googleLoginAction(
	idToken: string,
): Promise<VerifyOtpActionResult> {
	if (!idToken) {
		return { ok: false, message: "Google 認証情報を取得できませんでした" };
	}

	const payload: GoogleLoginRequest = { id_token: idToken };

	const res = await apiFetch<AuthTokensResponse>("/auth/google", {
		method: "POST",
		body: JSON.stringify(payload),
	});

	if (res.ok) {
		await setAuthCookies(res.data);
		return { ok: true };
	}
	return {
		ok: false,
		message: authErrorMessage(
			res.error,
			"Google ログインに失敗しました。時間をおいて再度お試しください",
		),
	};
}

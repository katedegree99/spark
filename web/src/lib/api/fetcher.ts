/**
 * サーバ専用 API フェッチャ。Server Action / Route Handler からのみ使う。
 *
 * backend へ server-to-server で直接アクセスする(ブラウザの rewrite プロキシは経由しない)。
 * baseURL は `next.config.ts` の rewrite と同じ `API_BASE_URL` を流用する。
 * トークンを httpOnly Cookie に隔離するため、auth の mutation はこの経路に通す。
 * `server-only` でクライアントからの誤 import をビルド時に弾く。
 */
import "server-only";
import type {
	ErrorResponse,
	ValidationErrorResponse,
} from "@/lib/api/generated/model";
import { getAccessToken } from "@/lib/auth/session";

const API_BASE_URL = process.env.API_BASE_URL ?? "http://localhost:3001";

export type ApiResult<T> =
	| { ok: true; status: number; data: T }
	| {
			ok: false;
			status: number;
			error: ErrorResponse | ValidationErrorResponse | null;
	  };

/**
 * JSON API をサーバ側で叩く。例外は投げず `ApiResult` で返す。
 * `Content-Type: application/json` を既定で付与し、レスポンスはキャッシュしない。
 *
 * `opts.auth = true` のとき、httpOnly Cookie のアクセストークンを
 * `Authorization: Bearer` で付与する(認証後の取得=B 経路で使う)。
 * 認証前の register/login/verify は `auth` を付けない。
 */
export async function apiFetch<T>(
	path: string,
	init?: RequestInit,
	opts?: { auth?: boolean },
): Promise<ApiResult<T>> {
	const headers: Record<string, string> = {
		"Content-Type": "application/json",
		...(init?.headers as Record<string, string>),
	};

	if (opts?.auth) {
		const token = await getAccessToken();
		if (token) {
			headers.Authorization = `Bearer ${token}`;
		}
	}

	const res = await fetch(`${API_BASE_URL}${path}`, {
		...init,
		headers,
		cache: "no-store",
	});

	const raw = [204, 205, 304].includes(res.status) ? "" : await res.text();
	// backend が不正な JSON(500 の HTML 等)を返しても例外を投げないよう、
	// パース失敗は body=null に倒す(「例外は投げず ApiResult で返す」契約を守る)。
	let body: unknown = null;
	if (raw) {
		try {
			body = JSON.parse(raw);
		} catch {
			body = null;
		}
	}

	if (res.ok) {
		return { ok: true, status: res.status, data: body as T };
	}
	return {
		ok: false,
		status: res.status,
		error: body as ErrorResponse | ValidationErrorResponse | null,
	};
}

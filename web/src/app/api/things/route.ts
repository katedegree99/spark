import type { NextRequest } from "next/server";
import { apiFetch } from "@/lib/api/fetcher";
import type { ThingsResponse } from "@/lib/api/generated/model";

/**
 * タグ候補検索のプロキシ Route Handler。
 *
 * 認証トークンが httpOnly Cookie にありクライアント JS から backend を
 * 直接叩けないため、サジェストのクライアント fetch はここを経由して
 * `GET /things?q=` に転送する。
 */
export async function GET(request: NextRequest) {
	const q = request.nextUrl.searchParams.get("q")?.trim() ?? "";
	// 空 query は backend を叩かず 0 件を返す(入力クリア時の無駄打ち防止)。
	if (q === "") {
		return Response.json({ things: [] } satisfies ThingsResponse);
	}

	const res = await apiFetch<ThingsResponse>(
		`/things?q=${encodeURIComponent(q)}`,
		undefined,
		{ auth: true },
	);
	if (!res.ok) {
		// backend のステータスを透過する。ネットワークエラー(status 0)は 502 に倒す。
		return Response.json(res.error ?? {}, {
			status: res.status === 0 ? 502 : res.status,
		});
	}
	return Response.json(res.data ?? ({ things: [] } satisfies ThingsResponse));
}

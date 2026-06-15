"use server";

import { redirect } from "next/navigation";
import { profileErrorMessage } from "@/features/profile/errors";
import { apiFetch } from "@/lib/api/fetcher";
import type {
	ProfileCreateRequest,
	ProfileResponse,
	ThingResponse,
	ThingsResponse,
} from "@/lib/api/generated/model";

/** 失敗時のみ返す(成功時はサーバ側で /home へ redirect する)。 */
export type ProfileActionResult = { ok: false; message: string };

/** 事柄作成の戻り値。 */
export type CreateThingResult =
	| { ok: true; thing: ThingResponse }
	| { ok: false; message: string };

/**
 * プロフィールを初回作成する Server Action(`POST /profiles/me`)。
 *
 * 認証後フローのため `auth: true`(httpOnly Cookie のトークンを Bearer 付与)。
 * 成功時は Cookie 不要だが、認証ガード配下の /home へ確実に遷移させるため
 * サーバ側 `redirect` で完結させる(auth フローと同じ方針)。
 */
export async function createProfileAction(input: {
	name: string;
	bio?: string | null;
	iconImageId?: number | null;
	doingThingIds: number[];
	wantThingIds: number[];
}): Promise<ProfileActionResult> {
	const payload: ProfileCreateRequest = {
		name: input.name,
		bio: input.bio ?? null,
		icon_image_id: input.iconImageId ?? null,
		doing_thing_ids: input.doingThingIds,
		want_thing_ids: input.wantThingIds,
	};

	const res = await apiFetch<ProfileResponse>(
		"/profiles/me",
		{ method: "POST", body: JSON.stringify(payload) },
		{ auth: true },
	);

	if (res.ok) {
		redirect("/home");
	}
	return {
		ok: false,
		message: profileErrorMessage(
			res.error,
			"プロフィールの作成に失敗しました。時間をおいて再度お試しください",
		),
	};
}

/**
 * 事柄(thing)をキーワード前方一致で検索する(`GET /things?q=`)。
 * 入力補助(サジェスト)用。失敗時は握り潰して空配列を返す。
 */
export async function searchThingsAction(q: string): Promise<ThingResponse[]> {
	const trimmed = q.trim();
	const path = trimmed ? `/things?q=${encodeURIComponent(trimmed)}` : "/things";

	const res = await apiFetch<ThingsResponse>(
		path,
		{ method: "GET" },
		{ auth: true },
	);
	return res.ok ? (res.data.things ?? []) : [];
}

/**
 * 候補に無い事柄をユーザー入力から新規作成する(`POST /things`)。
 * 成功時は作成された thing を返し、タグとして即追加できるようにする。
 */
export async function createThingAction(
	name: string,
): Promise<CreateThingResult> {
	const res = await apiFetch<ThingResponse>(
		"/things",
		{ method: "POST", body: JSON.stringify({ name }) },
		{ auth: true },
	);

	if (res.ok) {
		return { ok: true, thing: res.data };
	}
	return {
		ok: false,
		message: profileErrorMessage(res.error, "事柄の追加に失敗しました"),
	};
}

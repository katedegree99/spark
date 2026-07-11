"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { GoogleButton } from "@/components/ui/google-button";
import { googleLoginAction } from "@/features/auth/actions";

const CLIENT_ID = process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID;
const GSI_SRC = "https://accounts.google.com/gsi/client";

/** Google Identity Services の最小型(必要分のみ)。 */
type CredentialResponse = { credential?: string };
type GoogleIdApi = {
	initialize: (config: {
		client_id: string;
		callback: (res: CredentialResponse) => void;
		ux_mode?: "popup" | "redirect";
	}) => void;
	renderButton: (
		parent: HTMLElement,
		options: {
			type?: "standard" | "icon";
			theme?: "outline" | "filled_blue" | "filled_black";
			size?: "small" | "medium" | "large";
			text?: "signin_with" | "signup_with" | "continue_with" | "signin";
			width?: number;
		},
	) => void;
};
declare global {
	interface Window {
		google?: { accounts: { id: GoogleIdApi } };
	}
}

/** GIS スクリプトを一度だけ読み込む。 */
function loadGsi(): Promise<void> {
	return new Promise((resolve, reject) => {
		if (window.google?.accounts?.id) {
			resolve();
			return;
		}
		const existing = document.querySelector<HTMLScriptElement>(
			`script[src="${GSI_SRC}"]`,
		);
		if (existing) {
			existing.addEventListener("load", () => resolve());
			existing.addEventListener("error", () => reject(new Error("GIS load")));
			return;
		}
		const script = document.createElement("script");
		script.src = GSI_SRC;
		script.async = true;
		script.defer = true;
		script.onload = () => resolve();
		script.onerror = () => reject(new Error("GIS load failed"));
		document.head.appendChild(script);
	});
}

/**
 * 「Google でログイン/登録」ボタン。
 *
 * 見た目は自前の {@link GoogleButton} を維持しつつ、その上に **透明な公式 GIS ボタン**
 * を重ねてクリックを拾い、ID トークンを取得する(オーバーレイ方式)。
 * 取得した `credential`(= id_token)を {@link googleLoginAction} に渡し、
 * 成功でプロフィール設定済みなら `/home`、未設定なら `/profile/register` へ遷移する。
 *
 * `NEXT_PUBLIC_GOOGLE_CLIENT_ID` 未設定時はボタンを disabled 表示にして無効化する。
 */
export function GoogleLoginButton({ label }: { label: string }) {
	const wrapRef = useRef<HTMLDivElement>(null);
	const overlayRef = useRef<HTMLDivElement>(null);
	const [error, setError] = useState<string | null>(null);

	const onCredential = useCallback(async (res: CredentialResponse) => {
		if (!res.credential) {
			setError("Google 認証情報を取得できませんでした");
			return;
		}
		setError(null);
		// 成功時の Cookie 保存と遷移(プロフィール設定済みなら /home、未設定なら
		// /profile/register)は Server Action 側で完結する(redirect)。ここに戻るのは失敗時のみ。
		const result = await googleLoginAction(res.credential);
		if (result?.ok === false) {
			setError(result.message);
		}
	}, []);

	useEffect(() => {
		if (!CLIENT_ID) {
			return;
		}
		let cancelled = false;
		loadGsi()
			.then(() => {
				const overlay = overlayRef.current;
				if (cancelled || !overlay || !window.google) {
					return;
				}
				window.google.accounts.id.initialize({
					client_id: CLIENT_ID,
					callback: onCredential,
					ux_mode: "popup",
				});
				// 公式ボタンは px 幅指定(最大 400)。自前ボタンの実幅に合わせる。
				const width = Math.min(wrapRef.current?.offsetWidth ?? 360, 400);
				overlay.replaceChildren();
				window.google.accounts.id.renderButton(overlay, {
					type: "standard",
					theme: "outline",
					size: "large",
					text: "continue_with",
					width,
				});
			})
			.catch(() => setError("Google の読み込みに失敗しました"));
		return () => {
			cancelled = true;
		};
	}, [onCredential]);

	return (
		<div ref={wrapRef} className="relative w-full">
			<GoogleButton disabled={!CLIENT_ID}>{label}</GoogleButton>
			{/* 透明な公式 GIS ボタンを重ねてクリックを拾う(見た目は上の自前ボタン)。 */}
			<div
				ref={overlayRef}
				aria-hidden="true"
				className="absolute inset-0 z-10 flex items-stretch justify-center overflow-hidden opacity-0 [&_iframe]:!h-full [&>div]:!h-full"
			/>
			{error ? (
				<p className="mt-2 text-center text-error text-sm">{error}</p>
			) : null}
		</div>
	);
}

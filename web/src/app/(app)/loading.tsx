import { Loader2 } from "lucide-react";

/**
 * `(app)`(認証後)セグメントのローディング。
 *
 * App Router の Suspense フォールバック。Server Component の描画が終わるまで
 * 自動表示される(OTP / Google 認証後の `/home` 遷移中などに出る)。
 */
export default function AppLoading() {
	return (
		<div className="flex min-h-dvh items-center justify-center bg-white">
			<Loader2
				className="size-8 animate-spin text-brand-orange"
				aria-hidden="true"
			/>
			<span className="sr-only">読み込み中</span>
		</div>
	);
}

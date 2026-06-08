/**
 * 認証後ホーム (`/home`)。Server Component。
 *
 * 認証ガードは `(app)/layout.tsx` が担う。実データ取得は OpenAPI に認証後
 * エンドポイントが追加され次第、ここで `apiFetch(path, undefined, { auth: true })`
 * を用いて Server Component 取得する(推奨の B 経路)。
 */
export default function HomePage() {
	return (
		<main className="flex flex-1 flex-col items-center justify-center gap-2 p-6">
			<h1 className="font-bold text-2xl text-ink">ログインしました</h1>
			<p className="text-secondary text-sm">認証後ホーム(仮)</p>
		</main>
	);
}

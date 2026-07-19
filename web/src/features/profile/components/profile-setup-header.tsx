/**
 * プロフィール設定画面のヘッダー(中央タイトルのみ)。
 * プロフィール初回設定は認証直後の必須ステップで「戻る」で行ける正当な画面が
 * 存在しない(履歴上の前画面は使用済み OTP / ログイン)ため、戻るボタンは置かない。
 */
export function ProfileSetupHeader() {
	return (
		<header className="relative flex h-12 items-center justify-center">
			<h1 className="font-semibold text-ink text-xl">プロフィール設定</h1>
		</header>
	);
}

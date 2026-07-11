import { AuthSidePanel } from "@/features/auth/components/auth-side-panel";
import { ProfileRegisterForm } from "@/features/profile/components/profile-register-form";
import { ProfileSetupHeader } from "@/features/profile/components/profile-setup-header";

/**
 * プロフィール初回設定画面(`/profile/register`)。
 *
 * OTP / Google 認証後、`profileExists` が false のユーザーをここへ遷移させる。
 * 認証必須(`(setup)/layout.tsx` の requireSession ガード配下。アプリシェル無し)。
 *
 * レスポンシブ(auth 画面と同じ 2 カラム方針):
 * - SP(~md): 1 カラム。戻るボタン + 中央タイトルのヘッダー + フォーム
 * - PC(md~): 左にブランドサイドパネル(auth と共用)+ 右にフォーム。
 *   フォーム上部は左寄せタイトル「プロフィールを設定」(戻るボタンなし)。
 */
export default function ProfileRegisterPage() {
	return (
		<div className="flex w-full flex-1 bg-white">
			{/* PC: 左サイドパネル(最大幅 562px = Figma。それ未満は可変) */}
			<aside className="hidden min-w-0 flex-1 p-6 md:flex md:max-w-[562px]">
				<AuthSidePanel />
			</aside>

			{/* フォーム列(残り幅を取り中央寄せ、SP は全幅) */}
			<main className="flex w-full flex-col md:flex-1 md:items-center md:justify-center md:p-6">
				<div className="mx-auto flex w-full max-w-[430px] flex-1 flex-col px-5 pt-6 pb-10 md:max-w-[580px] md:flex-none md:px-0 md:py-0">
					{/* SP: 戻る + 中央タイトル */}
					<div className="md:hidden">
						<ProfileSetupHeader />
					</div>
					{/* PC: 左寄せタイトル */}
					<h1 className="hidden font-bold text-2xl text-ink tracking-[0.48px] md:block">
						プロフィールを設定
					</h1>

					<ProfileRegisterForm />
				</div>
			</main>
		</div>
	);
}

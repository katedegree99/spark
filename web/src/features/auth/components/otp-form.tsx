"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { ChevronLeft, Loader2 } from "lucide-react";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { Controller, useForm } from "react-hook-form";
import { Button } from "@/components/ui/button";
import { OtpInput } from "@/components/ui/otp-input";
import { verifyOtpAction } from "@/features/auth/actions";
// コンポーネント名 `OtpInput` と衝突するため、フォーム入力型は別名で import する。
import {
	type OtpInput as OtpFormInput,
	otpSchema,
} from "@/features/auth/schema";
import { useAuthFlowStore } from "@/features/auth/store";

/**
 * 認証コード(OTP)入力フォーム。`"use client"`。
 * `pendingEmail`(register / login で送信した宛先)宛のコードを検証する。
 * フロー外からの直接アクセス時は `/login` へリダイレクトする。
 *
 * - SP: 戻る + タイトルのヘッダー、本文中央、下部固定ボタン
 * - PC: サイドパネル(layout)+ 中央集約フォーム(ヘッダーなし)
 */
export function OtpForm() {
	const router = useRouter();
	const pendingEmail = useAuthFlowStore((s) => s.pendingEmail);
	const setPendingEmail = useAuthFlowStore((s) => s.setPendingEmail);
	const hasHydrated = useAuthFlowStore((s) => s.hasHydrated);
	const [formError, setFormError] = useState<string | null>(null);

	const {
		control,
		handleSubmit,
		formState: { isValid, isSubmitting },
	} = useForm<OtpFormInput>({
		resolver: zodResolver(otpSchema),
		mode: "onChange",
		defaultValues: { code: "" },
	});

	// フロー外からの直接アクセス(宛先メール未保持)はフローへ戻す。
	// sessionStorage からの復元完了(hasHydrated)を待ってから判定する
	// (復元前は pendingEmail が null のため、待たないと誤って /login へ飛ぶ)。
	useEffect(() => {
		if (hasHydrated && !pendingEmail) {
			router.replace("/login");
		}
	}, [hasHydrated, pendingEmail, router]);

	const onSubmit = handleSubmit(async ({ code }) => {
		// onComplete(自動送信)とボタン送信の二重起動を防ぐ。
		if (!pendingEmail || isSubmitting) {
			return;
		}
		setFormError(null);
		// 検証成功時、トークンは Server Action 側で httpOnly Cookie に保存される。
		const result = await verifyOtpAction({ email: pendingEmail, code });
		if (!result.ok) {
			setFormError(result.message);
			return;
		}
		// フロー完了。一時状態(宛先メール)を破棄してから遷移する。
		setPendingEmail(null);
		router.push("/home");
	});

	// sessionStorage 復元前、または宛先メール未保持(リダイレクト確定)の間は
	// 空メール表示・誤送信を避けるためローダーを出す。
	if (!hasHydrated || !pendingEmail) {
		return (
			<div className="flex flex-1 items-center justify-center">
				<Loader2
					className="size-8 animate-spin text-brand-orange"
					aria-hidden="true"
				/>
				<span className="sr-only">読み込み中</span>
			</div>
		);
	}

	return (
		<div className="flex flex-1 flex-col px-5 pt-4 pb-8 md:px-0 md:pt-0 md:pb-0">
			{/* SP のみ: 戻る + タイトルのヘッダー(PC は中央集約のため非表示) */}
			<header className="relative flex h-12 items-center justify-center md:hidden">
				<button
					type="button"
					onClick={() => router.back()}
					aria-label="戻る"
					className="absolute left-0 flex items-center justify-center rounded-full bg-white p-2 shadow-[0px_2px_2px_rgba(77,77,77,0.25)]"
				>
					<ChevronLeft
						className="size-8 text-ink"
						strokeWidth={1}
						aria-hidden="true"
					/>
				</button>
				<h1 className="font-bold text-ink text-lg">認証コードを入力</h1>
			</header>

			<form
				onSubmit={onSubmit}
				className="flex flex-1 flex-col md:items-center md:justify-center"
			>
				<div className="mt-20 flex flex-col items-center gap-2 text-center md:mt-0">
					<h2 className="font-bold text-ink text-xl">
						認証コードを入力してください
					</h2>
					<p className="text-secondary text-sm">
						{pendingEmail}に
						<br className="md:hidden" />
						6桁の認証コードを送信しました
					</p>
				</div>

				<div className="mt-8 w-full max-w-[360px] self-center">
					<Controller
						control={control}
						name="code"
						render={({ field }) => (
							<OtpInput {...field} onComplete={() => onSubmit()} />
						)}
					/>
				</div>

				{formError ? (
					<p className="mt-4 text-center text-error text-sm">{formError}</p>
				) : null}

				{/* TODO: 再送方式(register/login 再呼び出し)未確定のため一旦無効化 */}
				<button
					type="button"
					disabled
					className="mt-6 self-center text-secondary text-sm underline underline-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
				>
					コードを再送信
				</button>

				<p className="mt-6 text-center text-secondary text-sm">
					メールが届かない場合は迷惑メールフォルダをご確認ください
				</p>

				<Button
					type="submit"
					disabled={!isValid}
					loading={isSubmitting}
					className="mt-auto md:mt-8 md:w-full"
				>
					認証して次へ
				</Button>
			</form>
		</div>
	);
}

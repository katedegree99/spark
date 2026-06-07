"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { ChevronLeft } from "lucide-react";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { Controller, useForm } from "react-hook-form";
import { Button } from "@/components/ui/button";
import { OtpInput } from "@/components/ui/otp-input";
// コンポーネント名 `OtpInput` と衝突するため、フォーム入力型は別名で import する。
import {
	type OtpInput as OtpFormInput,
	otpSchema,
} from "@/features/auth/schema";
import { useAuthFlowStore } from "@/features/auth/store";
import { useVerifyOtp } from "@/lib/api/generated";

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
	const { trigger, isMutating } = useVerifyOtp();

	const {
		control,
		handleSubmit,
		formState: { isValid },
	} = useForm<OtpFormInput>({
		resolver: zodResolver(otpSchema),
		mode: "onChange",
		defaultValues: { code: "" },
	});

	// フロー外からの直接アクセス(宛先メール未保持)はフローへ戻す。
	useEffect(() => {
		if (!pendingEmail) {
			router.replace("/login");
		}
	}, [pendingEmail, router]);

	const onSubmit = handleSubmit(async ({ code }) => {
		if (!pendingEmail) {
			return;
		}
		await trigger({ email: pendingEmail, code });
		// TODO: トークン保管方式が未確定(nextjs-best-practices.md)。確定後にここで保存する
		router.push("/"); // 仮の遷移先
	});

	const handleResend = () => {
		// TODO: 再送方式(register/login 再呼び出し)未確定
	};

	return (
		<div className="flex flex-1 flex-col px-5 pt-4 pb-8 md:px-0 md:pt-0 md:pb-0">
			{/* SP のみ: 戻る + タイトルのヘッダー(PC は中央集約のため非表示) */}
			<header className="relative flex h-12 items-center justify-center md:hidden">
				<button
					type="button"
					onClick={() => router.back()}
					aria-label="戻る"
					className="absolute left-0 flex size-10 items-center justify-center rounded-full border border-border"
				>
					<ChevronLeft
						className="size-5"
						strokeWidth={1}
						stroke="url(#icon-gradient)"
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

				<button
					type="button"
					onClick={handleResend}
					className="mt-6 self-center text-secondary text-sm underline underline-offset-2"
				>
					コードを再送信
				</button>

				<p className="mt-6 text-center text-secondary text-sm">
					メールが届かない場合は迷惑メールフォルダをご確認ください
				</p>

				<Button
					type="submit"
					disabled={!isValid || isMutating}
					className="mt-auto md:mt-8 md:w-full"
				>
					認証して次へ
				</Button>
			</form>
		</div>
	);
}

"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { Mail } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { PasswordInput } from "@/components/ui/password-input";
import { loginAction } from "@/features/auth/actions";
import { GoogleLoginButton } from "@/features/auth/components/google-login-button";
import { type LoginInput, loginSchema } from "@/features/auth/schema";
import { useAuthFlowStore } from "@/features/auth/store";

/**
 * ログインフォーム。react-hook-form + zod でクライアント検証し、
 * 送信は `loginAction`(Server Action)に委譲する。
 * 成功時は宛先メールを `useAuthFlowStore` に保存して `/otp` へ遷移する。
 * (login 成功時もサーバは OTP を送る)
 *
 * 要素の並びは SP / PC で異なるため `order` で出し分ける:
 * - SP: 入力 → 登録導線 → Google → 利用規約 → ログインボタン
 * - PC: 入力 → 利用規約 → ログインボタン → Google(タブは page 側)
 */
export function LoginForm() {
	const router = useRouter();
	const setPendingEmail = useAuthFlowStore((s) => s.setPendingEmail);
	const [formError, setFormError] = useState<string | null>(null);

	const {
		register,
		handleSubmit,
		formState: { errors, isValid, isSubmitting },
	} = useForm<LoginInput>({
		resolver: zodResolver(loginSchema),
		mode: "onChange",
	});

	const onSubmit = handleSubmit(async (data) => {
		setFormError(null);
		const result = await loginAction(data);
		if (!result.ok) {
			setFormError(result.message);
			return;
		}
		setPendingEmail(result.email);
		router.push("/otp");
	});

	return (
		<form onSubmit={onSubmit} className="mt-7 flex flex-1 flex-col md:mt-5">
			<div className="order-1 flex flex-col gap-4">
				<div className="flex flex-col gap-1">
					<div className="flex items-baseline justify-between gap-2">
						<label htmlFor="email" className="font-semibold text-ink text-sm">
							メールアドレス
						</label>
						{errors.email ? (
							<span className="text-error text-sm text-right">
								{errors.email.message}
							</span>
						) : null}
					</div>
					<Input
						id="email"
						type="email"
						icon={Mail}
						placeholder="dummy@example.com"
						autoComplete="email"
						error={!!errors.email}
						{...register("email")}
					/>
				</div>

				<div className="flex flex-col gap-1">
					<div className="flex items-baseline justify-between gap-2">
						<label
							htmlFor="password"
							className="font-semibold text-ink text-sm"
						>
							パスワード
						</label>
						{errors.password ? (
							<span className="text-error text-sm text-right">
								{errors.password.message}
							</span>
						) : null}
					</div>
					<PasswordInput
						id="password"
						placeholder="password"
						autoComplete="current-password"
						error={!!errors.password}
						{...register("password")}
					/>
				</div>
			</div>

			{/* SP のみ: 新規登録への導線(PC はタブで切替) */}
			<p className="order-2 pt-10 text-center text-secondary text-xs md:hidden">
				アカウントをお持ちでない方は
				<Link
					href="/register"
					className="font-semibold underline underline-offset-2"
				>
					こちら
				</Link>
			</p>

			{/* 利用規約: SP は Google の後 / PC は入力の後 */}
			<p className="order-4 pt-3 text-center text-secondary text-xs underline underline-offset-2 md:order-2 md:pt-5 md:text-left">
				利用規約・プライバシーポリシー
			</p>

			{/* 送信失敗時のフォーム全体エラー(ボタン直上) */}
			{formError ? (
				<p className="order-5 mt-8 text-center text-error text-sm md:order-3 md:mt-5">
					{formError}
				</p>
			) : null}

			{/* ログインボタン: SP は最後 / PC は利用規約の後 */}
			<Button
				type="submit"
				disabled={!isValid}
				loading={isSubmitting}
				className="order-5 mt-8 md:order-3 md:mt-5"
			>
				ログイン
			</Button>

			{/* Google: SP は 登録導線 の後 / PC は最後 */}
			<div className="order-3 mt-5 border-border border-t pt-5 md:order-4">
				<GoogleLoginButton label="Googleでログイン" />
			</div>
		</form>
	);
}

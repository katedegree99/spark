"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { Mail } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { Button } from "@/components/ui/button";
import { GoogleButton } from "@/components/ui/google-button";
import { Input } from "@/components/ui/input";
import { PasswordInput } from "@/components/ui/password-input";
import { type LoginInput, loginSchema } from "@/features/auth/schema";
import { useAuthFlowStore } from "@/features/auth/store";
import { useLoginWithEmail } from "@/lib/api/generated";

/**
 * ログインフォーム。react-hook-form + zod で検証し、
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
	const { trigger, isMutating } = useLoginWithEmail();

	const {
		register,
		handleSubmit,
		formState: { errors, isValid },
	} = useForm<LoginInput>({
		resolver: zodResolver(loginSchema),
		mode: "onChange",
	});

	const onSubmit = handleSubmit(async ({ email, password }) => {
		await trigger({ email, password });
		setPendingEmail(email);
		router.push("/otp");
	});

	return (
		<form onSubmit={onSubmit} className="mt-7 flex flex-1 flex-col md:mt-5">
			<div className="order-1 flex flex-col gap-4">
				<div className="flex flex-col gap-1">
					<label htmlFor="email" className="font-semibold text-ink text-sm">
						メールアドレス
					</label>
					<Input
						id="email"
						type="email"
						icon={Mail}
						placeholder="dummy@example.com"
						autoComplete="email"
						{...register("email")}
					/>
					{errors.email ? (
						<p className="text-brand-red text-sm">{errors.email.message}</p>
					) : null}
				</div>

				<div className="flex flex-col gap-1">
					<label htmlFor="password" className="font-semibold text-ink text-sm">
						パスワード
					</label>
					<PasswordInput
						id="password"
						placeholder="password"
						autoComplete="current-password"
						{...register("password")}
					/>
					{errors.password ? (
						<p className="text-brand-red text-sm">{errors.password.message}</p>
					) : null}
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

			{/* ログインボタン: SP は最後 / PC は利用規約の後 */}
			<Button
				type="submit"
				disabled={!isValid || isMutating}
				className="order-5 mt-8 md:order-3 md:mt-5"
			>
				ログイン
			</Button>

			{/* Google: SP は 登録導線 の後 / PC は最後 */}
			<div className="order-3 mt-5 border-border border-t pt-5 md:order-4">
				<GoogleButton
					onClick={() => {
						// TODO: Google OAuth
					}}
				>
					Googleでログイン
				</GoogleButton>
			</div>
		</form>
	);
}

"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";
import { Controller, useForm } from "react-hook-form";
import { Button } from "@/components/ui/button";
import { createProfileAction } from "@/features/profile/actions";
import { AvatarPicker } from "@/features/profile/components/avatar-picker";
import { ThingTagInput } from "@/features/profile/components/thing-tag-input";
import { type ProfileInput, profileSchema } from "@/features/profile/schema";

/** ラベル + フィールド + エラー文言の縦並びラッパ。 */
function Field({
	label,
	error,
	children,
}: {
	label: string;
	error?: string;
	children: React.ReactNode;
}) {
	return (
		<div className="flex flex-col gap-2">
			<div className="flex items-center justify-between">
				<span className="font-semibold text-ink text-sm">{label}</span>
				{error ? <span className="text-error text-sm">{error}</span> : null}
			</div>
			{children}
		</div>
	);
}

/**
 * プロフィール設定フォーム。`"use client"`。
 * react-hook-form + zod でクライアント検証し、送信は `createProfileAction` に委譲する。
 * 成功時はサーバ側 redirect("/home") で完結するため、戻り値は失敗時のみ参照する。
 */
export function ProfileRegisterForm() {
	const [formError, setFormError] = useState<string | null>(null);

	const {
		register,
		handleSubmit,
		control,
		watch,
		formState: { errors, isSubmitting },
	} = useForm<ProfileInput>({
		resolver: zodResolver(profileSchema),
		mode: "onBlur",
		defaultValues: { name: "", bio: "", doings: [], wants: [] },
	});

	const name = watch("name");
	const bio = watch("bio");

	const onSubmit = handleSubmit(async (data) => {
		setFormError(null);
		const result = await createProfileAction({
			name: data.name,
			bio: data.bio,
			doingThingIds: data.doings.map((d) => d.id),
			wantThingIds: data.wants.map((w) => w.id),
		});
		// 成功時は createProfileAction 内で redirect 済み。戻るのは失敗時のみ。
		if (result?.ok === false) {
			setFormError(result.message);
		}
	});

	const inputClass =
		"w-full rounded-lg border border-border bg-white p-4 text-base text-ink outline-none transition-colors placeholder:text-border focus:border-gradient";

	return (
		<form onSubmit={onSubmit} className="flex flex-col gap-5 pt-5 md:gap-8">
			<AvatarPicker />

			<div className="flex flex-col gap-3">
				<Field label="ニックネーム" error={errors.name?.message}>
					<input
						{...register("name")}
						placeholder="山田　たかし"
						autoComplete="nickname"
						className={inputClass}
					/>
				</Field>

				<Field label="やっていること">
					<Controller
						control={control}
						name="doings"
						render={({ field }) => (
							<ThingTagInput
								value={field.value}
								onChange={field.onChange}
								placeholder="やっていること...."
							/>
						)}
					/>
				</Field>

				<Field label="これからやってみたいこと">
					<Controller
						control={control}
						name="wants"
						render={({ field }) => (
							<ThingTagInput
								value={field.value}
								onChange={field.onChange}
								placeholder="これからやってみたいこと...."
							/>
						)}
					/>
				</Field>

				<Field label="自己紹介" error={errors.bio?.message}>
					<textarea
						{...register("bio")}
						placeholder="自己紹介...."
						className={`${inputClass} h-[120px] resize-none`}
					/>
				</Field>
			</div>

			{formError ? (
				<p className="text-center text-error text-sm">{formError}</p>
			) : null}

			<Button
				type="submit"
				disabled={!name?.trim() || !bio?.trim()}
				loading={isSubmitting}
				className="mt-2 w-full text-xl tracking-[0.8px]"
			>
				はじめる
			</Button>
		</form>
	);
}

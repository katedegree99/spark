"use client";

import { Eye, EyeOff, Lock } from "lucide-react";
import { forwardRef, useState } from "react";
import { cn } from "@/utils/cn";

export type PasswordInputProps = Omit<
	React.InputHTMLAttributes<HTMLInputElement>,
	"type"
> & {
	/** 入力欄を内包する wrapper の className (レイアウト調整用, 任意) */
	wrapperClassName?: string;
	/** バリデーションエラー状態。枠線・アイコン・入力文字を赤にする。 */
	error?: boolean;
};

/**
 * パスワード入力。左に鍵 (`Lock`)、右に表示/非表示トグル (`Eye`/`EyeOff`)。
 * 角丸 `rounded-xl` + `border-border`、フォーカスで枠が `border-brand-orange`。
 * `error` 時は枠線・アイコン・入力文字を `brand-red` にする。
 * `forwardRef` でネイティブ `<input>` props を透過し、RHF の `register()` と結合可能。
 */
export const PasswordInput = forwardRef<HTMLInputElement, PasswordInputProps>(
	({ className, wrapperClassName, error, ...props }, ref) => {
		const [visible, setVisible] = useState(false);

		return (
			<div
				className={cn(
					"flex w-full items-center gap-2 rounded-xl border-[1.5px] bg-white py-3 pr-5 pl-4 transition-colors",
					error
						? "border-error drop-shadow-[2px_2px_2px_rgba(255,110,110,0.25)] focus-within:border-error"
						: "border-brand-yellow focus-within:border-brand-orange",
					wrapperClassName,
				)}
			>
				<Lock
					className={cn("size-6 shrink-0", error && "text-error")}
					strokeWidth={1}
					stroke={error ? "currentColor" : "url(#icon-gradient)"}
					aria-hidden="true"
				/>
				<input
					ref={ref}
					type={visible ? "text" : "password"}
					aria-invalid={error || undefined}
					className={cn(
						"min-w-0 flex-1 bg-transparent text-base outline-none placeholder:text-secondary disabled:cursor-not-allowed disabled:opacity-50",
						error ? "text-error" : "text-ink",
						className,
					)}
					{...props}
				/>
				<button
					type="button"
					onClick={() => setVisible((v) => !v)}
					className={cn(
						"shrink-0 transition-opacity hover:opacity-70 focus-visible:outline-none",
						error ? "text-error" : "text-brand-orange",
					)}
					aria-label={visible ? "パスワードを非表示" : "パスワードを表示"}
					aria-pressed={visible}
					tabIndex={-1}
				>
					{visible ? (
						<EyeOff
							className="size-8"
							strokeWidth={1}
							stroke={error ? "currentColor" : "url(#icon-gradient)"}
							aria-hidden="true"
						/>
					) : (
						<Eye
							className="size-8"
							strokeWidth={1}
							stroke={error ? "currentColor" : "url(#icon-gradient)"}
							aria-hidden="true"
						/>
					)}
				</button>
			</div>
		);
	},
);
PasswordInput.displayName = "PasswordInput";

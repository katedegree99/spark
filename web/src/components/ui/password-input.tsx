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
};

/**
 * パスワード入力。左に鍵 (`Lock`)、右に表示/非表示トグル (`Eye`/`EyeOff`)。
 * 角丸 `rounded-xl` + `border-border`、フォーカスで枠が `border-brand-orange`。
 * `forwardRef` でネイティブ `<input>` props を透過し、RHF の `register()` と結合可能。
 */
export const PasswordInput = forwardRef<HTMLInputElement, PasswordInputProps>(
	({ className, wrapperClassName, ...props }, ref) => {
		const [visible, setVisible] = useState(false);

		return (
			<div
				className={cn(
					"flex w-full items-center gap-2 rounded-xl border-[1.5px] border-brand-yellow bg-white py-3 pr-5 pl-4 transition-colors focus-within:border-brand-orange",
					wrapperClassName,
				)}
			>
				<Lock
					className="size-6 shrink-0"
					strokeWidth={1}
					stroke="url(#icon-gradient)"
					aria-hidden="true"
				/>
				<input
					ref={ref}
					type={visible ? "text" : "password"}
					className={cn(
						"min-w-0 flex-1 bg-transparent text-base text-ink outline-none placeholder:text-secondary disabled:cursor-not-allowed disabled:opacity-50",
						className,
					)}
					{...props}
				/>
				<button
					type="button"
					onClick={() => setVisible((v) => !v)}
					className="shrink-0 text-brand-orange transition-opacity hover:opacity-70 focus-visible:outline-none"
					aria-label={visible ? "パスワードを非表示" : "パスワードを表示"}
					aria-pressed={visible}
					tabIndex={-1}
				>
					{visible ? (
						<EyeOff
							className="size-8"
							strokeWidth={1}
							stroke="url(#icon-gradient)"
							aria-hidden="true"
						/>
					) : (
						<Eye
							className="size-8"
							strokeWidth={1}
							stroke="url(#icon-gradient)"
							aria-hidden="true"
						/>
					)}
				</button>
			</div>
		);
	},
);
PasswordInput.displayName = "PasswordInput";

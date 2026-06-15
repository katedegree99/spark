import { forwardRef } from "react";
import { GoogleIcon } from "@/components/icons/google";
import { cn } from "@/utils/cn";

export type GoogleButtonProps =
	React.ButtonHTMLAttributes<HTMLButtonElement> & {
		/** ボタンラベル。`children` が無い場合に使う (例: "Googleでログイン") */
		label?: React.ReactNode;
	};

/**
 * Google ロゴ (4色) + ラベル付きボタン。白背景 + `border-border`。
 * ラベルは `children` 優先、無ければ `label` prop。
 * `forwardRef` でネイティブ `<button>` props を透過する。
 */
export const GoogleButton = forwardRef<HTMLButtonElement, GoogleButtonProps>(
	({ className, label, children, type = "button", ...props }, ref) => {
		return (
			<button
				ref={ref}
				type={type}
				className={cn(
					"inline-flex w-full items-center justify-center gap-2 rounded-lg border-2 border-border bg-white py-4 font-semibold text-secondary tracking-wide transition-colors hover:bg-foreground/[0.03] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-orange focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50",
					className,
				)}
				{...props}
			>
				<GoogleIcon size={24} className="shrink-0" />
				{children ?? label}
			</button>
		);
	},
);
GoogleButton.displayName = "GoogleButton";

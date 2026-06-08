import type { LucideIcon } from "lucide-react";
import { forwardRef } from "react";
import { cn } from "@/utils/cn";

export type InputProps = React.InputHTMLAttributes<HTMLInputElement> & {
	/** 左側に表示する lucide アイコンコンポーネント (任意)。例: `icon={Mail}` */
	icon?: LucideIcon;
	/** 入力欄を内包する wrapper の className (レイアウト調整用, 任意) */
	wrapperClassName?: string;
	/** バリデーションエラー状態。枠線・アイコン・入力文字を赤にする。 */
	error?: boolean;
};

/**
 * 左アイコン (任意) 付きテキスト入力。
 * 角丸 `rounded-xl` + `border-border`、フォーカスで枠が `border-brand-orange`。
 * `error` 時は枠線・アイコン・入力文字を `brand-red` にする。
 * `forwardRef` でネイティブ `<input>` props (ref/name/onChange/onBlur 等) を透過し、
 * react-hook-form の `register()` 戻り値を spread できる。
 */
export const Input = forwardRef<HTMLInputElement, InputProps>(
	(
		{ className, wrapperClassName, icon: Icon, error, type = "text", ...props },
		ref,
	) => {
		return (
			<div
				className={cn(
					"flex w-full items-center gap-3 rounded-xl border-[1.5px] bg-white p-4 transition-colors",
					error
						? "border-error drop-shadow-[2px_2px_2px_rgba(255,110,110,0.25)] focus-within:border-error"
						: "border-brand-yellow focus-within:border-brand-orange",
					wrapperClassName,
				)}
			>
				{Icon ? (
					<Icon
						className={cn("size-6 shrink-0", error && "text-error")}
						strokeWidth={1}
						stroke={error ? "currentColor" : "url(#icon-gradient)"}
						aria-hidden="true"
					/>
				) : null}
				<input
					ref={ref}
					type={type}
					aria-invalid={error || undefined}
					className={cn(
						"min-w-0 flex-1 bg-transparent text-base outline-none placeholder:text-secondary disabled:cursor-not-allowed disabled:opacity-50",
						error ? "text-error" : "text-ink",
						className,
					)}
					{...props}
				/>
			</div>
		);
	},
);
Input.displayName = "Input";

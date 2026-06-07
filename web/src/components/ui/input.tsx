import type { LucideIcon } from "lucide-react";
import { forwardRef } from "react";
import { cn } from "@/utils/cn";

export type InputProps = React.InputHTMLAttributes<HTMLInputElement> & {
	/** 左側に表示する lucide アイコンコンポーネント (任意)。例: `icon={Mail}` */
	icon?: LucideIcon;
	/** 入力欄を内包する wrapper の className (レイアウト調整用, 任意) */
	wrapperClassName?: string;
};

/**
 * 左アイコン (任意) 付きテキスト入力。
 * 角丸 `rounded-xl` + `border-border`、フォーカスで枠が `border-brand-orange`。
 * `forwardRef` でネイティブ `<input>` props (ref/name/onChange/onBlur 等) を透過し、
 * react-hook-form の `register()` 戻り値を spread できる。
 */
export const Input = forwardRef<HTMLInputElement, InputProps>(
	(
		{ className, wrapperClassName, icon: Icon, type = "text", ...props },
		ref,
	) => {
		return (
			<div
				className={cn(
					"flex w-full items-center gap-3 rounded-xl border-[1.5px] border-brand-yellow bg-white p-4 transition-colors focus-within:border-brand-orange",
					wrapperClassName,
				)}
			>
				{Icon ? (
					<Icon
						className="size-6 shrink-0"
						strokeWidth={1}
						stroke="url(#icon-gradient)"
						aria-hidden="true"
					/>
				) : null}
				<input
					ref={ref}
					type={type}
					className={cn(
						"min-w-0 flex-1 bg-transparent text-base text-ink outline-none placeholder:text-secondary disabled:cursor-not-allowed disabled:opacity-50",
						className,
					)}
					{...props}
				/>
			</div>
		);
	},
);
Input.displayName = "Input";

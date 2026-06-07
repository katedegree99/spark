import { cva, type VariantProps } from "class-variance-authority";
import { forwardRef } from "react";
import { cn } from "@/utils/cn";

const buttonVariants = cva(
	"inline-flex items-center justify-center gap-2 rounded-xl font-bold text-white transition-opacity focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-brand-orange focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50",
	{
		variants: {
			variant: {
				/** メインアクション: orange→yellow グラデーション + 白文字 */
				gradient: "bg-brand-gradient",
			},
			size: {
				sm: "h-10 px-4 text-sm",
				md: "h-12 px-6 text-base",
				lg: "py-5 text-xl tracking-wide",
			},
			fullWidth: {
				true: "w-full",
			},
		},
		defaultVariants: {
			variant: "gradient",
			size: "lg",
			fullWidth: true,
		},
	},
);

export type ButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement> &
	VariantProps<typeof buttonVariants>;

/**
 * メインボタン。`cva` バリアント (`gradient` / `size` / `fullWidth`) を持ち、
 * ネイティブ `<button>` props を透過する。非活性時は `disabled:opacity-50`。
 */
export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
	({ className, variant, size, fullWidth, type = "button", ...props }, ref) => {
		return (
			<button
				ref={ref}
				type={type}
				className={cn(buttonVariants({ variant, size, fullWidth }), className)}
				{...props}
			/>
		);
	},
);
Button.displayName = "Button";

export { buttonVariants };

import { cn } from "@/utils/cn";

type TagChipProps = {
	label: string;
	/** filled=共通タグ(グラデ塗り)/ outline=非共通(白地・グラデ文字)。 */
	variant?: "filled" | "outline";
	className?: string;
};

// 両バリアント共通: 1px 枠線(色は variant で指定) + 角丸 full + 12/8 パディング(Figma)。
// 枠線の「幅」だけ共通化し、色は filled/outline で上書きする。
const BASE =
	"inline-flex items-center rounded-full border px-3 py-2 font-semibold text-xs leading-none";

/**
 * タグチップ(Figma 準拠)。
 * filled(共通タグ): グラデ塗り + 塗りより明るいグラデ枠 + 白文字。Figma の
 *   tag/main は塗りの外側にブランドグラデの枠リングが乗る二重構造。
 * outline(非共通): 白背景 + ブランドグラデ枠 + ブランドグラデのクリップ文字。
 */
export function TagChip({
	label,
	variant = "outline",
	className,
}: TagChipProps) {
	if (variant === "filled") {
		return (
			<span
				className={cn(BASE, "bg-brand-gradient-tag-fill text-white", className)}
			>
				{label}
			</span>
		);
	}
	return (
		<span className={cn(BASE, "border-brand-gradient-tag", className)}>
			<span className="bg-brand-gradient-tag bg-clip-text text-transparent">
				{label}
			</span>
		</span>
	);
}

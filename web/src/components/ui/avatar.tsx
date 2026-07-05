import { UserRound } from "lucide-react";
import { cn } from "@/utils/cn";

// Figma の avatar トークンに合わせる(sm=44 / md=52 / lg=72)。
const SIZE_CLASSES = {
	sm: "size-11",
	md: "size-13",
	lg: "size-18",
} as const;

type AvatarProps = {
	/** 画像 URL。null/欠落時はデフォルトの人型アイコンを表示する。 */
	src: string | null;
	/** alt / aria-label に使う表示名。 */
	name: string;
	size?: keyof typeof SIZE_CLASSES;
	className?: string;
};

/**
 * 円形アバター。`next/image` は使わない(`next.config.ts` に
 * `remotePatterns` 未設定のため、設定を触らず素の `<img>` で実装する)。
 * 画像が無い(アイコン未設定の)ときは人型のデフォルトアイコンを表示する。
 */
export function Avatar({ src, name, size = "md", className }: AvatarProps) {
	// name が空(API 由来で `?? ""` になりうる)のとき、アクセシブル名が
	// 空文字の img ロールにならないよう代替ラベルにフォールバックする。
	const label = name.trim() || "ユーザー";
	return (
		<div
			className={cn(
				"inline-flex shrink-0 items-center justify-center overflow-hidden rounded-full bg-border/40 text-secondary",
				SIZE_CLASSES[size],
				className,
			)}
			{...(src ? {} : { role: "img", "aria-label": label })}
		>
			{src ? (
				// biome-ignore lint/performance/noImgElement: remotePatterns 未設定のため next/image は使わない
				<img src={src} alt={label} className="size-full object-cover" />
			) : (
				<UserRound className="size-3/5" strokeWidth={1.5} aria-hidden="true" />
			)}
		</div>
	);
}

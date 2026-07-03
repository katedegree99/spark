import { cn } from "@/utils/cn";

// Figma の avatar トークンに合わせる(sm=44 / md=52 / lg=72)。
const SIZE_CLASSES = {
	sm: "size-11 text-base",
	md: "size-13 text-lg",
	lg: "size-18 text-2xl",
} as const;

type AvatarProps = {
	/** 画像 URL。null/欠落時はイニシャル fallback を表示する。 */
	src: string | null;
	/** alt / イニシャル算出に使う表示名。 */
	name: string;
	size?: keyof typeof SIZE_CLASSES;
	className?: string;
};

/** 表示名の先頭 1 文字を取り出す(空なら "?")。 */
function initial(name: string): string {
	return [...name.trim()][0] ?? "?";
}

/**
 * 円形アバター。`next/image` は使わない(`next.config.ts` に
 * `remotePatterns` 未設定のため、設定を触らず素の `<img>` で実装する)。
 * 画像が無いときは表示名のイニシャルにフォールバックする。
 */
export function Avatar({ src, name, size = "md", className }: AvatarProps) {
	// name が空(API 由来で `?? ""` になりうる)のとき、アクセシブル名が
	// 空文字の img ロールにならないよう代替ラベルにフォールバックする。
	const label = name.trim() || "ユーザー";
	return (
		<div
			className={cn(
				"inline-flex shrink-0 items-center justify-center overflow-hidden rounded-full bg-border/40 font-bold text-secondary",
				SIZE_CLASSES[size],
				className,
			)}
			{...(src ? {} : { role: "img", "aria-label": label })}
		>
			{src ? (
				// biome-ignore lint/performance/noImgElement: remotePatterns 未設定のため next/image は使わない
				<img src={src} alt={label} className="size-full object-cover" />
			) : (
				<span aria-hidden="true">{initial(name)}</span>
			)}
		</div>
	);
}

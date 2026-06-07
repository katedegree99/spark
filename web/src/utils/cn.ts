import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

/**
 * Tailwind クラスを条件付きで結合し、競合するユーティリティを解決する。
 * clsx で条件付きクラスをまとめ、tailwind-merge で後勝ちにマージする。
 *
 * @example
 * cn("px-2 py-1", isActive && "bg-blue-500", "px-4") // => "py-1 bg-blue-500 px-4"
 */
export function cn(...inputs: ClassValue[]): string {
	return twMerge(clsx(inputs));
}

import type { SVGProps } from "react";

type GoogleIconProps = SVGProps<SVGSVGElement> & {
	/** アイコンの一辺の長さ (px)。width/height に適用する。default: 20 */
	size?: number | string;
};

/**
 * Google 公式 4色 "G" ロゴ。lucide には含まれないため自作。
 * Google ブランドガイドラインに沿った 4 パス構成。
 */
export function GoogleIcon({ size = 20, ...props }: GoogleIconProps) {
	return (
		<svg
			width={size}
			height={size}
			viewBox="0 0 48 48"
			xmlns="http://www.w3.org/2000/svg"
			role="img"
			aria-label="Google"
			{...props}
		>
			<path
				fill="#4285F4"
				d="M47.532 24.552c0-1.638-.147-3.213-.42-4.725H24.48v8.939h12.94c-.557 3.003-2.252 5.547-4.799 7.25v6.02h7.76c4.542-4.183 7.151-10.342 7.151-17.484Z"
			/>
			<path
				fill="#34A853"
				d="M24.48 48c6.48 0 11.913-2.149 15.884-5.815l-7.76-6.02c-2.149 1.44-4.9 2.291-8.124 2.291-6.249 0-11.54-4.218-13.43-9.888H2.998v6.215C6.948 42.62 14.114 48 24.48 48Z"
			/>
			<path
				fill="#FBBC05"
				d="M11.05 28.568c-.48-1.44-.756-2.977-.756-4.568 0-1.59.276-3.128.756-4.568v-6.215H2.998A23.94 23.94 0 0 0 .48 24c0 3.873.927 7.537 2.518 10.783l8.052-6.215Z"
			/>
			<path
				fill="#EA4335"
				d="M24.48 9.544c3.522 0 6.683 1.21 9.169 3.587l6.882-6.882C36.387 2.376 30.954 0 24.48 0 14.114 0 6.948 5.38 2.998 13.217l8.052 6.215c1.89-5.67 7.181-9.888 13.43-9.888Z"
			/>
		</svg>
	);
}

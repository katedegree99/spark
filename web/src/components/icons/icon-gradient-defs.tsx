/**
 * lucide アイコンを stroke にグラデーションで塗るための共有 <defs>。
 * 各アイコンに `stroke="url(#icon-gradient)"` を指定すると、この勾配で線が塗られる。
 * 画面に出ない hidden SVG。root layout の body 直下に 1 度だけ置く。
 */
export function IconGradientDefs() {
	return (
		<svg
			aria-hidden="true"
			focusable="false"
			width="0"
			height="0"
			className="absolute"
		>
			<title>icon gradient defs</title>
			<defs>
				<linearGradient id="icon-gradient" x1="0" y1="0" x2="1" y2="1">
					<stop offset="0%" stopColor="var(--brand-yellow)" />
					<stop offset="50%" stopColor="var(--brand-orange)" />
					<stop offset="100%" stopColor="var(--brand-red)" />
				</linearGradient>
			</defs>
		</svg>
	);
}

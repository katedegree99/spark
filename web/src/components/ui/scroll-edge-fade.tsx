/**
 * 横スクロール領域の左右端に重ねる白フェード(装飾)。
 * `relative` な親の中に置き、スクロール中の要素が画面端で自然に
 * 消え込むように見せる。クリックは透過する。
 */
export function ScrollEdgeFade() {
	return (
		<>
			<div className="pointer-events-none absolute inset-y-0 left-0 w-6 bg-gradient-to-r from-white to-transparent" />
			<div className="pointer-events-none absolute inset-y-0 right-0 w-6 bg-gradient-to-l from-white to-transparent" />
		</>
	);
}

import Image from "next/image";

/**
 * PC レイアウトの左サイドパネル(md 以上で表示)。
 * ブランドグラデ背景にロゴ・キャッチコピー・装飾を配置する。Server Component。
 */
export function AuthSidePanel() {
	return (
		<div className="relative flex h-full w-full flex-col overflow-hidden rounded-[28px] bg-brand-panel p-12 text-white">
			<div className="flex items-center gap-2.5">
				<Image
					src="/images/spark-logo.png"
					alt=""
					width={816}
					height={831}
					className="size-14 object-contain"
				/>
				<span className="font-bold text-[40px] tracking-[1.6px]">SPARK</span>
			</div>

			<div className="absolute inset-x-12 top-[45%] flex flex-col gap-6">
				<div className="flex items-center gap-2.5">
					<span className="h-px w-6 bg-white" />
					<span className="font-semibold text-lg">WELCOME TO SPARK</span>
				</div>
				<h2 className="font-bold text-[28px] leading-snug">
					「やりたい」で、
					<br />
					１対１でつながる。
				</h2>
				<p className="font-semibold text-base leading-relaxed tracking-wide">
					いま夢中なこと、これから挑戦したいこと。
					<br />
					やっていること・やりたいことが重なる人と、1on1でじっっくり話せます。
				</p>
			</div>

			<Image
				src="/images/spark-logo.png"
				alt=""
				aria-hidden="true"
				width={816}
				height={831}
				className="pointer-events-none absolute right-[-21%] bottom-[-13%] h-[44%] w-auto opacity-25 blur-[2px]"
			/>
		</div>
	);
}

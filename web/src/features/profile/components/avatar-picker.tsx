"use client";

import { useEffect, useRef, useState } from "react";

const ACCEPT = "image/jpeg,image/png";
const MAX_BYTES = 5 * 1024 * 1024; // 5MB

/**
 * プロフィールアイコン選択。
 *
 * アバター(または「アイコンを設定」/「写真をアップロード」)をタップすると
 * ファイル選択ダイアログを開き、選んだ画像をプレビュー表示する。
 * JPG/PNG・5MB までを検証する(Figma のキャプション準拠)。
 *
 * NOTE: R2 への実アップロード(POST /images)は未配線のため、現状は
 * ローカルプレビューのみ(icon_image_id は付与しない)。アップロード実装時に
 * 選択ファイルを送信し、返却 id をフォームへ渡す配線を追加する。
 */
export function AvatarPicker() {
	const inputRef = useRef<HTMLInputElement>(null);
	const [previewUrl, setPreviewUrl] = useState<string | null>(null);
	const [error, setError] = useState<string | null>(null);

	// objectURL のメモリリークを防ぐ。
	useEffect(() => {
		return () => {
			if (previewUrl) URL.revokeObjectURL(previewUrl);
		};
	}, [previewUrl]);

	const openPicker = () => inputRef.current?.click();

	const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
		const file = e.target.files?.[0];
		if (!file) return;
		if (file.type !== "image/jpeg" && file.type !== "image/png") {
			setError("JPG / PNG 画像を選択してください");
			return;
		}
		if (file.size > MAX_BYTES) {
			setError("5MB 以下の画像を選択してください");
			return;
		}
		setError(null);
		if (previewUrl) URL.revokeObjectURL(previewUrl);
		setPreviewUrl(URL.createObjectURL(file));
	};

	return (
		<div className="flex flex-col items-center gap-2 md:items-start">
			<div className="flex w-full flex-col items-center gap-3 md:flex-row md:gap-5">
				<button
					type="button"
					onClick={openPicker}
					aria-label="アイコン画像を選択"
					className="rounded-full bg-brand-gradient p-[2px]"
				>
					<div className="flex size-[100px] items-center justify-center overflow-hidden rounded-full bg-white">
						{previewUrl ? (
							// biome-ignore lint/performance/noImgElement: ローカル objectURL のプレビューのため next/image は不要
							<img
								src={previewUrl}
								alt="選択したアイコン"
								className="size-full object-cover"
							/>
						) : (
							// Figma のアイコン: 細線の人型シルエット。肩は大きな円の上端で表現し、
							// 親の rounded-full + overflow-hidden で円外をクリップして「はみ出す」見た目にする。
							<svg
								viewBox="0 0 100 100"
								className="size-full text-border"
								fill="none"
								stroke="currentColor"
								strokeWidth={2}
								aria-hidden="true"
							>
								<title>アイコン未設定</title>
								{/* 頭(r15)と、肩を表す大きめの円(上端=59)。頭下端(56)との間に
								    3px ほどのネック隙間。肩は円縁で overflow-hidden にクリップされ、
								    肩先が円からはみ出して切れる Figma の見た目になる。 */}
								<circle cx="50" cy="40" r="17" />
								<circle cx="50" cy="95" r="36" />
							</svg>
						)}
					</div>
				</button>

				{/* SP: ラベル(タップでピッカー) */}
				<button
					type="button"
					onClick={openPicker}
					className="font-semibold text-ink text-sm md:hidden"
				>
					アイコンを設定
				</button>

				{/* PC: アップロードボタン + 制約キャプション(タップでピッカー) */}
				<div className="hidden flex-col items-start gap-2 md:flex">
					<button
						type="button"
						onClick={openPicker}
						className="rounded-full border-[1.5px] border-brand-yellow px-4 py-2.5"
					>
						<span className="bg-brand-gradient bg-clip-text font-medium text-base text-transparent tracking-[0.4px]">
							写真をアップロード
						</span>
					</button>
					<span className="font-medium text-secondary text-sm tracking-[0.4px]">
						JPG/PNG 5MBまで
					</span>
				</div>

				<input
					ref={inputRef}
					type="file"
					accept={ACCEPT}
					onChange={onChange}
					className="hidden"
				/>
			</div>

			{error ? <p className="text-error text-sm">{error}</p> : null}
		</div>
	);
}

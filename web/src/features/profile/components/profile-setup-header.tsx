"use client";

import { ChevronLeft } from "lucide-react";
import { useRouter } from "next/navigation";

/**
 * プロフィール設定画面のヘッダー(戻るボタン + 中央タイトル)。
 * 見た目は OTP 画面の SP ヘッダーに合わせる(白丸 + shadow の戻るボタン)。
 */
export function ProfileSetupHeader() {
	const router = useRouter();

	return (
		<header className="relative flex h-12 items-center justify-center">
			<button
				type="button"
				onClick={() => router.back()}
				aria-label="戻る"
				className="absolute left-0 flex items-center justify-center rounded-full bg-white p-2 shadow-[0px_2px_2px_rgba(77,77,77,0.25)]"
			>
				<ChevronLeft
					className="size-8 text-ink"
					strokeWidth={1}
					aria-hidden="true"
				/>
			</button>
			<h1 className="font-semibold text-ink text-xl">プロフィール設定</h1>
		</header>
	);
}

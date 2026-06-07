"use client";

import { useId, useRef } from "react";
import { cn } from "@/utils/cn";

const OTP_LENGTH = 6;
/** 固定長ボックスの安定キー(index をキーに使わないため)。 */
const SLOTS = ["s0", "s1", "s2", "s3", "s4", "s5"] as const;

export type OtpInputProps = {
	/** 制御値。長さ 0〜6 の数字文字列。 */
	value: string;
	/** 値が変わるたびに呼ばれる(常に長さ 0〜6 の数字文字列)。 */
	onChange: (value: string) => void;
	/** 全 6 桁が揃った瞬間に呼ばれる。 */
	onComplete?: (value: string) => void;
	/** react-hook-form 連携用(任意)。 */
	name?: string;
	/** react-hook-form 連携用(任意)。最後のボックスを抜けたとき発火。 */
	onBlur?: () => void;
	disabled?: boolean;
	/** 入力全体に対する aria-label(各ボックスの連番ラベルの土台)。 */
	"aria-label"?: string;
	className?: string;
};

/** 数字以外を除去し、最大 6 桁に切り詰める。 */
function sanitize(raw: string): string {
	return raw.replace(/\D/g, "").slice(0, OTP_LENGTH);
}

/**
 * 6 桁の OTP(ワンタイムコード)入力。1 桁ずつの分割ボックスで、
 * 入力で次へ自動フォーカス、Backspace で前へ、6 桁の貼り付けで一括入力する。
 * value / onChange による controlled component。
 */
export function OtpInput({
	value,
	onChange,
	onComplete,
	name,
	onBlur,
	disabled,
	"aria-label": ariaLabel = "認証コード",
	className,
}: OtpInputProps) {
	const groupId = useId();
	const inputsRef = useRef<(HTMLInputElement | null)[]>([]);

	const digits = sanitize(value);

	const focusIndex = (index: number) => {
		const target = inputsRef.current[index];
		if (target) {
			target.focus();
			target.select();
		}
	};

	const emit = (next: string) => {
		const sanitized = sanitize(next);
		onChange(sanitized);
		if (sanitized.length === OTP_LENGTH) {
			onComplete?.(sanitized);
		}
	};

	const handleChange = (index: number, raw: string) => {
		const incoming = sanitize(raw);
		if (incoming.length === 0) {
			return;
		}
		// 1 ボックスに複数文字が入ったとき(ペースト含む)は index 位置から埋める。
		const chars = digits.split("");
		let cursor = index;
		for (const ch of incoming) {
			if (cursor >= OTP_LENGTH) break;
			chars[cursor] = ch;
			cursor += 1;
		}
		const next = chars.join("").slice(0, OTP_LENGTH);
		emit(next);
		focusIndex(Math.min(cursor, OTP_LENGTH - 1));
	};

	const handleKeyDown = (
		index: number,
		event: React.KeyboardEvent<HTMLInputElement>,
	) => {
		if (event.key === "Backspace") {
			event.preventDefault();
			const chars = digits.split("");
			if (chars[index]) {
				// 現在のボックスに値があればそれを消す。
				chars[index] = "";
				emit(chars.join(""));
			} else if (index > 0) {
				// 空なら前のボックスへ戻って消す。
				chars[index - 1] = "";
				emit(chars.join(""));
				focusIndex(index - 1);
			}
			return;
		}
		if (event.key === "ArrowLeft" && index > 0) {
			event.preventDefault();
			focusIndex(index - 1);
			return;
		}
		if (event.key === "ArrowRight" && index < OTP_LENGTH - 1) {
			event.preventDefault();
			focusIndex(index + 1);
		}
	};

	const handlePaste = (
		index: number,
		event: React.ClipboardEvent<HTMLInputElement>,
	) => {
		event.preventDefault();
		const pasted = sanitize(event.clipboardData.getData("text"));
		if (pasted.length === 0) {
			return;
		}
		const chars = digits.split("");
		let cursor = index;
		for (const ch of pasted) {
			if (cursor >= OTP_LENGTH) break;
			chars[cursor] = ch;
			cursor += 1;
		}
		const next = chars.join("").slice(0, OTP_LENGTH);
		emit(next);
		focusIndex(Math.min(cursor, OTP_LENGTH - 1));
	};

	const handleBlur = (index: number) => {
		// 最後のボックスからフォーカスが抜けたときだけ form 側に通知。
		if (index === OTP_LENGTH - 1) {
			onBlur?.();
		}
	};

	return (
		<fieldset
			aria-label={ariaLabel}
			className={cn(
				"m-0 flex w-full items-center justify-between border-0 p-0",
				className,
			)}
		>
			{SLOTS.map((slot, index) => {
				const char = digits[index] ?? "";
				const filled = char !== "";
				return (
					<input
						key={`${groupId}-${slot}`}
						ref={(el) => {
							inputsRef.current[index] = el;
						}}
						// react-hook-form 連携時は最初のボックスに name を付与して register と紐付けやすくする。
						name={index === 0 ? name : undefined}
						type="text"
						inputMode="numeric"
						autoComplete={index === 0 ? "one-time-code" : "off"}
						pattern="\d*"
						maxLength={OTP_LENGTH}
						disabled={disabled}
						value={char}
						aria-label={`${ariaLabel} ${index + 1} 桁目`}
						onChange={(e) => handleChange(index, e.target.value)}
						onKeyDown={(e) => handleKeyDown(index, e)}
						onPaste={(e) => handlePaste(index, e)}
						onBlur={() => handleBlur(index)}
						onFocus={(e) => e.currentTarget.select()}
						className={cn(
							"h-[84px] w-[50px] rounded-xl border-2 text-center font-semibold text-3xl text-ink",
							"outline-none transition-colors",
							"focus:border-brand-orange",
							filled ? "border-brand-orange" : "border-border",
							disabled && "cursor-not-allowed opacity-50",
						)}
					/>
				);
			})}
		</fieldset>
	);
}

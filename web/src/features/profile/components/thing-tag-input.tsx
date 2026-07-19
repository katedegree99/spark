"use client";

import { Plus, X } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import {
	createThingAction,
	searchThingsAction,
} from "@/features/profile/actions";
import { MAX_THINGS, type ThingTag } from "@/features/profile/schema";

type Props = {
	/** 選択済みタグ。react-hook-form の Controller 経由で渡す。 */
	value: ThingTag[];
	onChange: (next: ThingTag[]) => void;
	placeholder?: string;
	/** 選択できるタグの上限(既定 20)。 */
	max?: number;
};

/** ThingResponse(全フィールド optional)を表示可能なタグへ詰める。 */
function toTag(t: { id?: number; name?: string }): ThingTag | null {
	return t.id != null && t.name != null ? { id: t.id, name: t.name } : null;
}

/**
 * 事柄(thing)の複数選択タグ入力。Figma の「チップ枠 + セレクトボックス」構成。
 *
 * - チップ枠: 選択済みをブランドグラデのチップ(× で削除)で表示。クリックで開く。
 * - セレクトボックス(開いている間だけ表示): 上部に検索入力、下に候補リスト。
 *   各候補は「+ ◯◯」(+ はオレンジアイコン、名前はグラデ文字)。候補に無ければ
 *   「+『◯◯』を追加」で `POST /things` し、作成したタグを選択に加える。
 * - トークンは httpOnly Cookie のためクライアントから直接 API を叩けず、
 *   検索/作成はいずれも Server Action 経由(Bearer 付与はサーバ側で実施)。
 */
export function ThingTagInput({
	value,
	onChange,
	placeholder,
	max = MAX_THINGS,
}: Props) {
	const [query, setQuery] = useState("");
	const [suggestions, setSuggestions] = useState<ThingTag[]>([]);
	// suggestions がどの query に対する結果か(null = まだどの入力にも対応していない)。
	const [settledQuery, setSettledQuery] = useState<string | null>(null);
	const [open, setOpen] = useState(false);
	const [pending, setPending] = useState(false);
	const [error, setError] = useState<string | null>(null);
	const containerRef = useRef<HTMLDivElement>(null);
	const inputRef = useRef<HTMLInputElement>(null);

	// 開いている間だけ検索(空クエリは全件 = 初期候補)。250ms デバウンス。
	useEffect(() => {
		if (!open) {
			setSuggestions([]);
			setSettledQuery(null);
			return;
		}
		let cancelled = false;
		const requestedQuery = query;
		const timer = setTimeout(async () => {
			const things = await searchThingsAction(requestedQuery);
			if (cancelled) return;
			setSuggestions(
				things.map(toTag).filter((t): t is ThingTag => t !== null),
			);
			setSettledQuery(requestedQuery);
		}, 250);
		return () => {
			cancelled = true;
			clearTimeout(timer);
		};
	}, [query, open]);

	// セレクトボックス外をクリックしたら閉じる。
	useEffect(() => {
		if (!open) return;
		const onPointerDown = (e: PointerEvent) => {
			if (
				containerRef.current &&
				!containerRef.current.contains(e.target as Node)
			) {
				setOpen(false);
			}
		};
		document.addEventListener("pointerdown", onPointerDown);
		return () => document.removeEventListener("pointerdown", onPointerDown);
	}, [open]);

	const selectedIds = new Set(value.map((v) => v.id));
	const filtered = suggestions.filter((s) => !selectedIds.has(s.id));
	const trimmed = query.trim();
	// suggestions が現在の入力に対応しているか(debounce 待ち・フェッチ中は false)。
	// ステイルな候補での Enter 確定や、existsの誤判定による重複 POST /things を防ぐ。
	const settled = settledQuery === query;
	const exists =
		suggestions.some((s) => s.name === trimmed) ||
		value.some((v) => v.name === trimmed);
	const atMax = value.length >= max;
	const canCreate = settled && trimmed.length > 0 && !exists && !atMax;

	const openBox = () => {
		setOpen(true);
		// 描画後に検索入力へフォーカスする。
		requestAnimationFrame(() => inputRef.current?.focus());
	};

	const add = (tag: ThingTag) => {
		if (selectedIds.has(tag.id)) {
			setQuery("");
			return;
		}
		if (atMax) {
			setError(`最大${max}個まで選択できます`);
			return;
		}
		setError(null);
		onChange([...value, tag]);
		// suggestions はクリアしない。選択で候補一覧自体は変わらず(選択済みは
		// filtered で除外される)、空クエリから選んだ場合は query が変化せず
		// 再フェッチも走らないため、クリアすると候補が空のまま残ってしまう。
		// suggestions と settledQuery を対で残すので exists の判定も正しいまま。
		setQuery("");
		inputRef.current?.focus();
	};

	const remove = (id: number) => {
		onChange(value.filter((v) => v.id !== id));
	};

	const create = async () => {
		if (!trimmed || pending) return;
		setPending(true);
		setError(null);
		const res = await createThingAction(trimmed);
		setPending(false);
		if (res.ok) {
			const tag = toTag(res.thing);
			if (tag) add(tag);
		} else {
			// 失敗を無音にしない(例: backend 未実装で 404 のとき原因が分かるように)。
			setError(res.message);
		}
	};

	return (
		<div ref={containerRef} className="relative w-full">
			{/* チップ表示エリア(クリックでセレクトボックスを開く) */}
			{/* biome-ignore lint/a11y/useSemanticElements: 内部に削除用 button を持つため button をネストできず、role="button" の div で代替する */}
			<div
				role="button"
				tabIndex={0}
				onClick={openBox}
				onKeyDown={(e) => {
					if (e.key === "Enter" || e.key === " ") {
						e.preventDefault();
						openBox();
					}
				}}
				className={`flex min-h-[104px] w-full cursor-text flex-wrap content-start items-start gap-2 rounded-lg border bg-white p-4 ${open ? "border-brand-gradient" : "border-border"}`}
			>
				{value.length === 0 ? (
					<span className="text-base text-border">{placeholder}</span>
				) : (
					value.map((tag) => (
						<span
							key={tag.id}
							className="inline-flex items-center gap-1 rounded-full bg-brand-gradient py-2 pr-2 pl-3 font-semibold text-white text-xs"
						>
							{tag.name}
							<button
								type="button"
								onClick={(e) => {
									e.stopPropagation();
									remove(tag.id);
								}}
								aria-label={`${tag.name} を削除`}
								className="flex items-center justify-center"
							>
								<X className="size-4" aria-hidden="true" />
							</button>
						</span>
					))
				)}
			</div>

			{/* セレクトボックス(検索入力 + 候補リスト) */}
			{open ? (
				<div className="absolute z-20 mt-1 w-full rounded-xl border border-border bg-white p-4 shadow-[0px_4px_4px_rgba(0,0,0,0.25)]">
					<input
						ref={inputRef}
						value={query}
						onChange={(e) => setQuery(e.target.value)}
						onKeyDown={(e) => {
							if (e.key === "Enter") {
								// IME 変換確定の Enter で誤選択・誤作成(POST /things)しないようガードする。
								if (e.nativeEvent.isComposing) return;
								e.preventDefault();
								// 候補が現在の入力に追いつくまでは確定しない(前 query の候補の誤選択防止)。
								if (!settled) return;
								if (filtered[0]) add(filtered[0]);
								else if (canCreate) create();
							}
						}}
						placeholder={placeholder}
						disabled={atMax}
						className="w-full bg-transparent text-base text-ink outline-none placeholder:text-border disabled:cursor-not-allowed"
					/>

					{atMax ? (
						<p className="mt-2 text-secondary text-sm">
							最大{max}個まで選択できます
						</p>
					) : null}

					{!atMax && (filtered.length > 0 || canCreate) ? (
						<ul className="mt-2 flex max-h-60 flex-col gap-1 overflow-auto">
							{filtered.map((s) => (
								<li key={s.id}>
									<button
										type="button"
										onClick={() => add(s)}
										className="flex w-full items-center gap-3 py-2 text-left"
									>
										<Plus
											className="size-6 shrink-0 text-brand-orange"
											aria-hidden="true"
										/>
										<span className="bg-brand-gradient bg-clip-text font-semibold text-base text-transparent">
											{s.name}
										</span>
									</button>
								</li>
							))}
							{canCreate ? (
								<li>
									<button
										type="button"
										onClick={create}
										disabled={pending}
										className="flex w-full items-center gap-3 py-2 text-left disabled:opacity-50"
									>
										<Plus
											className="size-6 shrink-0 text-brand-orange"
											aria-hidden="true"
										/>
										<span className="bg-brand-gradient bg-clip-text font-semibold text-base text-transparent">
											「{trimmed}」を追加
										</span>
									</button>
								</li>
							) : null}
						</ul>
					) : null}

					{error ? <p className="mt-1 text-error text-sm">{error}</p> : null}
				</div>
			) : null}
		</div>
	);
}

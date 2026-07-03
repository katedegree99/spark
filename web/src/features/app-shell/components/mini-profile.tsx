import { Avatar } from "@/components/ui/avatar";
import type { MiniProfileVM } from "../data";

/**
 * サイドバー最下部の自分のミニプロフィール(Server)。データは props で受ける。
 * サイドバーのアイコンのみ表示(lg 未満)に合わせ、テキストは lg 以上でだけ出す。
 */
export function MiniProfile({ me }: { me: MiniProfileVM }) {
	return (
		<div className="flex w-full items-center justify-center gap-3 lg:justify-start">
			<Avatar src={me.iconUrl} name={me.name} size="md" className="size-15" />
			<div className="hidden min-w-0 flex-col gap-1 lg:flex">
				<p className="truncate font-medium text-2xl text-ink tracking-[0.48px]">
					{me.name}
				</p>
				{me.handle ? (
					<p className="truncate text-base text-secondary tracking-[0.32px]">
						{me.handle}
					</p>
				) : null}
			</div>
		</div>
	);
}

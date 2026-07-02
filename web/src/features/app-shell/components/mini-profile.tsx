import { Avatar } from "@/components/ui/avatar";
import type { MiniProfileVM } from "../data";

/**
 * サイドバー最下部の自分のミニプロフィール(Server)。データは props で受ける。
 */
export function MiniProfile({ me }: { me: MiniProfileVM }) {
	return (
		<div className="flex w-full items-center gap-3">
			<Avatar src={me.iconUrl} name={me.name} size="md" className="size-15" />
			<div className="flex min-w-0 flex-col gap-1">
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

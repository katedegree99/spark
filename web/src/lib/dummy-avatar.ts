/**
 * 開発用のダミーアバター解決。
 *
 * mock サーバー(Redocly)は例データとして無効な `pub-xxx.r2.dev` の URL を返すため、
 * そのままだと画像が壊れる。この「壊れ URL」のときだけ Unsplash のダミー
 * ポートレートに置き換えて見た目を保つ(`seed`=userId 等で 1 ユーザー 1 画像に固定)。
 * 真に未設定(null/空)のときは null を返し、Avatar のデフォルトアイコン表示に委ねる。
 *
 * TODO(mock): backend が本番のプロフィール画像 URL を返すようになったら削除する。
 */
const DUMMY_AVATARS = [
	"1500648767791-00dcc994a43e",
	"1494790108377-be9c29b29330",
	"1438761681033-6461ffad8d80",
	"1507003211169-0a1dd7228f2d",
	"1544005313-94ddf0286df2",
	"1506794778202-cad84cf45f1d",
	"1534528741775-53994a69daeb",
	"1489424731084-a5d8b219a5bb",
	"1531123897727-8f129e1688ce",
	"1522075469751-3a6694fb2f61",
	"1524504388940-b1c1722653e1",
	"1500336624523-d727130c3328",
].map(
	(id) =>
		`https://images.unsplash.com/photo-${id}?w=200&h=200&fit=crop&crop=faces&auto=format&q=80`,
);

/**
 * mock の無効 URL(実体が無い placeholder)。これらは壊れるためダミーに置換する。
 *
 * 'pub-xxx' は openapi.yaml の example に書かれた R2 の placeholder ホスト
 * (`https://pub-xxx.r2.dev/...`)で、Redocly mock がそのまま返してくる。
 * 実際の R2 公開 URL は `pub-<実ハッシュ>.r2.dev` 形式なので誤検知はしないが、
 * mock の例データ文字列に依存した判定である点に注意(このヘルパーごと削除予定)。
 */
function isUsableIconUrl(url: string | null | undefined): url is string {
	return url != null && url !== "" && !url.includes("pub-xxx");
}

/**
 * 使える iconUrl はそのまま、mock の壊れ URL は `seed` から決まるダミー画像に置換、
 * 未設定(null/空)は null を返す(Avatar 側でデフォルトアイコンにフォールバック)。
 */
export function resolveAvatarUrl(
	iconUrl: string | null | undefined,
	seed: number,
): string | null {
	if (iconUrl == null || iconUrl === "") return null;
	if (isUsableIconUrl(iconUrl)) return iconUrl;
	const i =
		((Math.trunc(seed) % DUMMY_AVATARS.length) + DUMMY_AVATARS.length) %
		DUMMY_AVATARS.length;
	return DUMMY_AVATARS[i];
}

/**
 * 新着・おすすめユーザーの一時モックデータ。
 *
 * これらのセクションに対応する API(`GET /users/new-arrivals` /
 * `GET /users/recommended`)は未実装。将来 API 化したら **このファイルを削除**し、
 * `data.ts` の該当関数を `apiFetch` 取得に差し替えるだけで済むよう、モック実体は
 * ここに隔離している(コンポーネントはモックを直接 import しない)。
 */
import type {
	NewArrivalVM,
	PickupCardVM,
	RecommendedUserVM,
	TagVM,
} from "./types";

/** Unsplash のポートレート写真(顔中心・正方形クロップ)をダミー画像に使う。 */
const unsplash = (id: string) =>
	`https://images.unsplash.com/photo-${id}?w=200&h=200&fit=crop&crop=faces&auto=format&q=80`;

/** タグ名プール。実データはユーザーごとに最大 40 個ほどタグが付く想定。 */
const TAG_POOL = [
	"コーディング",
	"IT",
	"水族館",
	"デザイン",
	"カフェ",
	"音楽",
	"読書",
	"料理",
	"写真",
	"散歩",
	"ランニング",
	"キャンプ",
	"ボードゲーム",
	"映画",
	"アート",
	"美術館",
	"旅行",
	"ゲーム",
	"アニメ",
	"ギター",
	"ライブ",
	"コーヒー",
	"お菓子",
	"登山",
	"釣り",
	"ヨガ",
	"筋トレ",
	"カメラ",
	"ドライブ",
	"温泉",
	"ラーメン",
	"ワイン",
	"ダンス",
	"英会話",
	"ボルダリング",
	"サウナ",
	"猫",
	"犬",
	"ガジェット",
	"ポッドキャスト",
];

let tagSeq = 0;
const nextId = (): number => {
	tagSeq += 1;
	return tagSeq;
};
/**
 * 共通タグ(matched)を先頭に、残りをプールから unmatched として補充し計 `total` 個の
 * タグ配列を作る。実データ規模(最大〜40 個)を再現して、カードのタグ行が
 * 「入るだけ表示 + …」で見切れる様子を確認できるようにする。
 */
function makeTags(matched: string[], total: number): TagVM[] {
	const tags: TagVM[] = matched.map((name) => ({
		id: nextId(),
		name,
		matched: true,
	}));
	for (const name of TAG_POOL) {
		if (tags.length >= total) break;
		if (matched.includes(name)) continue;
		tags.push({ id: nextId(), name, matched: false });
	}
	return tags;
}

/**
 * 今日のピックアップのダミーデータ。
 * 実 API `GET /users/pickup` がまだ本番データを返せない開発環境向けに、
 * `data.ts` が取得失敗/空のとき暫定表示するフォールバック。API が本番データを
 * 返せるようになったら不要になる(このファイルごと削除できる)。
 */
export const pickupCardsFixture: PickupCardVM[] = [
	{
		userId: 7001,
		name: "堺 理人",
		bio: "週末はコードを書いたり水族館に行ったりしています。IT 系の話題で気軽に盛り上がれる人と繋がりたいです。",
		iconUrl: unsplash("1500648767791-00dcc994a43e"),
		matchedCount: 3,
		tags: makeTags(["コーディング", "IT", "水族館"], 40),
	},
	{
		userId: 7002,
		name: "藤原 美咲",
		bio: "デザインとカフェ巡りが好きです。同じ趣味の友達を探しています。まずは気軽にお話しできたら嬉しいです。",
		iconUrl: unsplash("1494790108377-be9c29b29330"),
		matchedCount: 2,
		tags: makeTags(["デザイン", "カフェ"], 24),
	},
	{
		userId: 7003,
		name: "田中 蓮",
		bio: "音楽ライブとキャンプによく行きます。アウトドア好きな人、一緒に出かけられる仲間を募集中です。",
		iconUrl: unsplash("1507003211169-0a1dd7228f2d"),
		matchedCount: 1,
		tags: makeTags(["音楽"], 33),
	},
	{
		userId: 7004,
		name: "小林 さくら",
		bio: "読書と料理が趣味のインドア派です。ゆっくり語り合える人と出会えたらいいなと思っています。",
		iconUrl: unsplash("1438761681033-6461ffad8d80"),
		matchedCount: 2,
		tags: makeTags(["読書", "料理"], 18),
	},
];

export const newArrivalsFixture: NewArrivalVM[] = [
	{
		userId: 9001,
		name: "あおい",
		iconUrl: unsplash("1494790108377-be9c29b29330"),
	},
	{
		userId: 9002,
		name: "はると",
		iconUrl: unsplash("1500648767791-00dcc994a43e"),
	},
	{
		userId: 9003,
		name: "ゆい",
		iconUrl: unsplash("1438761681033-6461ffad8d80"),
	},
	{
		userId: 9004,
		name: "そうた",
		iconUrl: unsplash("1507003211169-0a1dd7228f2d"),
	},
	{ userId: 9005, name: "めい", iconUrl: unsplash("1544005313-94ddf0286df2") },
	{
		userId: 9006,
		name: "りく",
		iconUrl: unsplash("1506794778202-cad84cf45f1d"),
	},
	{
		userId: 9007,
		name: "さくら",
		iconUrl: unsplash("1534528741775-53994a69daeb"),
	},
	{
		userId: 9008,
		name: "かいと",
		iconUrl: unsplash("1489424731084-a5d8b219a5bb"),
	},
];

export const recommendedUsersFixture: RecommendedUserVM[] = [
	{
		userId: 8001,
		name: "みなと",
		bio: "週末はカフェ巡りと写真を撮るのが好きです。気軽に話せる人を探しています。",
		iconUrl: unsplash("1531123897727-8f129e1688ce"),
		matchedCount: 3,
		tags: makeTags(["カフェ", "写真", "散歩"], 40),
	},
	{
		userId: 8002,
		name: "ひなた",
		bio: "ランニングとアウトドアが趣味。一緒に体を動かせる友達がほしいです。",
		iconUrl: unsplash("1522075469751-3a6694fb2f61"),
		matchedCount: 2,
		tags: makeTags(["ランニング", "キャンプ"], 22),
	},
	{
		userId: 8003,
		name: "つむぎ",
		bio: "音楽ライブによく行きます。同じアーティストが好きな人と繋がりたい。",
		iconUrl: unsplash("1524504388940-b1c1722653e1"),
		matchedCount: 1,
		tags: makeTags(["音楽"], 30),
	},
	{
		userId: 8004,
		name: "あさひ",
		bio: "読書とボードゲームが好きなインドア派。ゆるく語り合える人を募集中。",
		iconUrl: unsplash("1500336624523-d727130c3328"),
		matchedCount: 2,
		tags: makeTags(["読書", "ボードゲーム"], 16),
	},
	{
		userId: 8005,
		name: "ゆうな",
		bio: "美術館巡りと映画鑑賞が好きです。感想をゆっくり語り合える人と繋がりたい。",
		iconUrl: unsplash("1544005313-94ddf0286df2"),
		matchedCount: 3,
		tags: makeTags(["映画", "アート", "美術館"], 40),
	},
	{
		userId: 8006,
		name: "そら",
		bio: "旅行とカメラが趣味。週末は近場をふらっと撮り歩いています。一緒に出かけられる人募集。",
		iconUrl: unsplash("1506794778202-cad84cf45f1d"),
		matchedCount: 2,
		tags: makeTags(["旅行", "写真"], 27),
	},
	{
		userId: 8007,
		name: "いおり",
		bio: "料理とお菓子作りにはまっています。おすすめのお店やレシピを交換できたら嬉しいです。",
		iconUrl: unsplash("1534528741775-53994a69daeb"),
		matchedCount: 1,
		tags: makeTags(["料理"], 12),
	},
	{
		userId: 8008,
		name: "はるき",
		bio: "ゲームとアニメが好きなインドア派です。同じ作品が好きな人とゆるく話したいです。",
		iconUrl: unsplash("1489424731084-a5d8b219a5bb"),
		matchedCount: 2,
		tags: makeTags(["ゲーム", "アニメ"], 35),
	},
];

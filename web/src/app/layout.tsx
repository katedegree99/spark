import type { Metadata } from "next";
import { Barlow, Roboto } from "next/font/google";
import { IconGradientDefs } from "@/components/icons/icon-gradient-defs";
import "./globals.css";

const roboto = Roboto({
	variable: "--font-roboto",
	subsets: ["latin"],
	weight: ["400", "500", "700", "900"],
	display: "swap",
});

// Figma の数字用フォント「DIN Bold」の代替。DIN は商用のため Google Fonts に
// 無く、DIN 系グロテスクの Barlow(700)で近似する(利用は `font-din` 経由)。
const barlow = Barlow({
	variable: "--font-barlow",
	subsets: ["latin"],
	weight: "700",
	display: "swap",
});

export const metadata: Metadata = {
	title: "Spark",
	description: "Spark",
};

export default function RootLayout({
	children,
}: Readonly<{
	children: React.ReactNode;
}>) {
	return (
		<html
			lang="ja"
			className={`${roboto.variable} ${barlow.variable} h-full antialiased`}
		>
			<body className="flex min-h-full flex-col">
				<IconGradientDefs />
				{children}
			</body>
		</html>
	);
}

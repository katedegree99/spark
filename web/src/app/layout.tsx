import type { Metadata } from "next";
import { Roboto } from "next/font/google";
import localFont from "next/font/local";
import { IconGradientDefs } from "@/components/icons/icon-gradient-defs";
import "./globals.css";

const roboto = Roboto({
	variable: "--font-roboto",
	subsets: ["latin"],
	weight: ["400", "500", "700", "900"],
	display: "swap",
});

// Figma の数字用フォント「DIN Bold」を、OSS の DIN 復刻フォント D-DIN(SIL OFL)で
// 再現する。Google Fonts に DIN は無いため self-host(next/font/local)。
// ライセンス: src/app/fonts/D-DIN-OFL.txt
const dDin = localFont({
	src: "./fonts/D-DIN-Bold.woff2",
	variable: "--font-d-din",
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
			className={`${roboto.variable} ${dDin.variable} h-full antialiased`}
		>
			<body className="flex min-h-full flex-col">
				<IconGradientDefs />
				{children}
			</body>
		</html>
	);
}

import type { Metadata } from "next";
import { Roboto } from "next/font/google";
import { IconGradientDefs } from "@/components/icons/icon-gradient-defs";
import "./globals.css";

const roboto = Roboto({
	variable: "--font-roboto",
	subsets: ["latin"],
	weight: ["400", "500", "700"],
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
		<html lang="ja" className={`${roboto.variable} h-full antialiased`}>
			<body className="flex min-h-full flex-col">
				<IconGradientDefs />
				{children}
			</body>
		</html>
	);
}

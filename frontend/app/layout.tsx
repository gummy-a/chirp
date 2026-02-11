import Initializer from "@/components/initializer";
import "./globals.css";

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="jp">
      <body className="wrap-break-word px-4 max-w-2xl mx-auto">
        <Initializer />
        {children}
      </body>
    </html>
  );
}

import type { Metadata } from "next";
import type { ReactNode } from "react";

import "./globals.css";
import { SessionProvider } from "@/providers/session-provider";

export const metadata: Metadata = {
  title: "Dvorfs Repository Manager",
  description: "A minimal UI for managing repositories and artifacts.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: ReactNode;
}>) {
  return (
    <html lang="en">
      <body>
        <SessionProvider>{children}</SessionProvider>
      </body>
    </html>
  );
}

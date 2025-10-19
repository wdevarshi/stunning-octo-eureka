import type { Metadata } from 'next';
import './globals.css';

export const metadata: Metadata = {
  title: 'LTA Transport Reliability Dashboard',
  description: 'Real-time dashboard for Singapore transport reliability analytics',
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className="antialiased bg-gray-50">{children}</body>
    </html>
  );
}

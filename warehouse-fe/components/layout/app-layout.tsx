"use client";

import { Sidebar } from "./sidebar";

interface AppLayoutProps {
  children: React.ReactNode;
}

export function AppLayout({ children }: AppLayoutProps) {
  return (
    <div className="flex h-screen w-full bg-zinc-50 text-foreground">
      <Sidebar />
      <main className="flex-1 overflow-y-auto p-4 md:p-8 transition-all">
        <div className="mx-auto max-w-7xl animate-fade-in space-y-6">
          {children}
        </div>
      </main>
    </div>
  );
}

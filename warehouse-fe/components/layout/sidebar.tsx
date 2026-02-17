"use client";

import { usePathname, useRouter } from "next/navigation";
import Link from "next/link";
import {
  HomeIcon,
  CubeIcon,
  ShoppingCartIcon,
  ClipboardDocumentListIcon,
  ChartBarIcon,
  UsersIcon,
  ArrowRightOnRectangleIcon,
} from "@heroicons/react/24/outline";
import { useAuth } from "@/components/providers/auth-provider";
import { Avatar, Button, Tooltip } from "@nextui-org/react";
import { toast } from "sonner";

export function Sidebar() {
  const pathname = usePathname();
  const { user, logout } = useAuth();
  const router = useRouter();

  const handleLogout = async () => {
    try {
      await logout();
      toast.success("Logged out successfully");
      router.push("/login"); // Ensure redirect
    } catch (error) {
      toast.error("Failed to logout");
    }
  };

  const menuItems = [
    { key: "dashboard", name: "Dashboard", href: "/dashboard", icon: HomeIcon },
    { key: "barang", name: "Inventory", href: "/barang", icon: CubeIcon },
    {
      key: "penjualan",
      name: "Sales",
      href: "/penjualan",
      icon: ShoppingCartIcon,
    },
    {
      key: "pembelian",
      name: "Purchasing",
      href: "/pembelian",
      icon: ClipboardDocumentListIcon,
    },
    { key: "stok", name: "Stock History", href: "/stok", icon: ChartBarIcon },
  ];

  if (user?.role === "admin") {
    menuItems.push({
      key: "users",
      name: "Users & Roles",
      href: "/users",
      icon: UsersIcon,
    });
  }

  return (
    <aside className="hidden h-screen w-72 flex-col border-r border-slate-200 bg-white/80 backdrop-blur-md px-4 py-8 md:flex shadow-sm z-50">
      <div className="mb-10 flex items-center gap-3 px-3">
        <div className="flex h-11 w-11 items-center justify-center rounded-xl bg-gradient-to-br from-green-600 to-teal-700 text-white shadow-lg shadow-green-600/30 ring-4 ring-green-50">
          <CubeIcon className="h-6 w-6 text-white" />
        </div>
        <div className="flex flex-col">
          <span className="text-lg font-bold tracking-tight text-slate-800">
            ROXY
          </span>
          <span className="text-xs font-medium text-slate-400 tracking-wider uppercase">
            Management System
          </span>
        </div>
      </div>

      <div className="flex-1 overflow-y-auto px-2 space-y-1.5 py-4">
        <p className="px-4 text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">
          Menu
        </p>

        {menuItems.map((item) => {
          const isActive = pathname.startsWith(item.href);
          return (
            <Link
              key={item.href}
              href={item.href}
              className={`group relative flex items-center gap-3 rounded-xl px-4 py-3 text-sm font-medium transition-all duration-300 ease-out border border-transparent ${
                isActive
                  ? "bg-green-50/80 text-green-700 shadow-sm border-green-100"
                  : "text-slate-500 hover:bg-slate-50 hover:text-slate-900 hover:shadow-sm"
              }`}>
              {isActive && (
                <span className="absolute left-0 top-1/2 -translate-y-1/2 w-1 h-6 bg-green-600 rounded-r-full" />
              )}

              <item.icon
                className={`h-5 w-5 transition-transform duration-300 group-hover:scale-110 ${
                  isActive
                    ? "text-green-600"
                    : "text-slate-400 group-hover:text-slate-600"
                }`}
              />
              <span className={isActive ? "font-semibold" : ""}>
                {item.name}
              </span>
            </Link>
          );
        })}
      </div>

      <div className="mt-auto px-2 pt-6 border-t border-slate-100">
        <div className="flex items-center justify-between gap-3 p-3 rounded-2xl bg-slate-50 border border-slate-100 shadow-sm transition-all hover:shadow-md hover:border-slate-200 cursor-pointer group">
          <div className="flex items-center gap-3 overflow-hidden">
            <Avatar
              name={user?.username}
              className="w-9 h-9 text-tiny font-semibold bg-gradient-to-br from-slate-700 to-slate-900 text-white shadow-sm ring-2 ring-white"
            />
            <div className="flex flex-col truncate">
              <span className="text-sm font-semibold text-slate-700 truncate group-hover:text-slate-900 transition-colors">
                {user?.username || "Guest User"}
              </span>
              <span className="text-xs text-slate-400 capitalize truncate group-hover:text-slate-500 transition-colors">
                {user?.role || "Visitor"}
              </span>
            </div>
          </div>

          <Tooltip content="Logout">
            <Button
              isIconOnly
              size="sm"
              variant="light"
              color="danger"
              onClick={handleLogout}
              aria-label="Logout"
              className="text-slate-400 hover:text-red-500 hover:bg-red-50 transition-colors rounded-lg min-w-8 w-8 h-8">
              <ArrowRightOnRectangleIcon className="h-4 w-4" />
            </Button>
          </Tooltip>
        </div>
      </div>
    </aside>
  );
}

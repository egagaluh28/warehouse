"use client";

import { useEffect, useState, useCallback, Key } from "react";
import {
  Card,
  CardBody,
  CardHeader,
  Spinner,
  Divider,
} from "@nextui-org/react";
import {
  UsersIcon,
  CubeIcon,
  CurrencyDollarIcon,
  ChartBarIcon,
} from "@heroicons/react/24/outline";
import { AppLayout } from "@/components/layout/app-layout";
import { dashboardApi } from "@/lib/api";
import { DashboardStats } from "@/lib/types";
import { DataTable, Column } from "@/components/table/data-table";

// Component for Stat Card
interface StatCardProps {
  title: string;
  value?: number | string;
  icon: React.ElementType;
  className?: string; // Tailwind class for background color
  iconClassName?: string; // Tailwind class for icon color
}

const StatCard = ({
  title,
  value,
  icon: Icon,
  className,
  iconClassName,
}: StatCardProps) => (
  <Card className="border-none bg-white shadow-sm ring-1 ring-gray-100 transition-all hover:shadow-md">
    <CardBody className="flex flex-row items-center justify-between gap-4 p-6">
      <div className="flex flex-col gap-1">
        <p className="text-sm font-medium text-gray-500">{title}</p>
        <h3 className="text-2xl font-bold tracking-tight text-gray-900">
          {typeof value === "number"
            ? value.toLocaleString("id-ID")
            : value || 0}
        </h3>
      </div>
      <div className={`rounded-xl p-3 ${className || "bg-gray-100"}`}>
        <Icon className={`h-6 w-6 ${iconClassName || "text-gray-600"}`} />
      </div>
    </CardBody>
  </Card>
);

export default function DashboardPage() {
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const topProductColumns: Column[] = [
    { uid: "nama_barang", name: "NAMA PRODUK" },
    { uid: "total_terjual", name: "TERJUAL", align: "end" },
  ];

  const renderTopProductCell = (
    item: { nama_barang: string; total_terjual: number },
    columnKey: Key,
  ) => {
    switch (columnKey) {
      case "nama_barang":
        return <span className="font-medium">{item.nama_barang}</span>;
      case "total_terjual":
        return (
          <div className="flex justify-end gap-1">
            <span className="font-semibold text-green-600">
              {item.total_terjual}
            </span>
            <span className="text-gray-500">unit</span>
          </div>
        );
      default:
        return null;
    }
  };

  const fetchStats = useCallback(async () => {
    try {
      const data = await dashboardApi.getStats();
      setStats(data);
    } catch (error) {
      console.error("Failed to fetch dashboard stats", error);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchStats();
  }, [fetchStats]);

  const formatCurrency = (value: number = 0) => {
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
    }).format(value);
  };

  if (isLoading) {
    return (
      <AppLayout>
        <div className="flex h-[80vh] w-full items-center justify-center">
          <Spinner size="md" color="success" label="Memuat dashboard..." />
        </div>
      </AppLayout>
    );
  }

  return (
    <AppLayout>
      <div className="flex h-full w-full flex-col gap-6 p-6">
        {/* Header */}
        <div className="flex flex-col gap-1">
          <h1 className="text-2xl font-bold tracking-tight text-gray-900">
            Dashboard Overview
          </h1>
          <p className="text-sm text-gray-500">
            Ringkasan aktivitas hari ini dan status gudang.
          </p>
        </div>

        {/* Stats Grid */}
        <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
          <StatCard
            title="Total Barang"
            value={stats?.total_barang}
            icon={CubeIcon}
            className="bg-blue-50"
            iconClassName="text-blue-600"
          />
          <StatCard
            title="Total Stok"
            value={stats?.total_stok}
            icon={ChartBarIcon}
            className="bg-green-50"
            iconClassName="text-green-600"
          />
          <StatCard
            title="Total Aset"
            value={formatCurrency(stats?.total_nilai_aset)}
            icon={CurrencyDollarIcon}
            className="bg-amber-50"
            iconClassName="text-amber-600"
          />
          <StatCard
            title="Total Pengguna"
            value={stats?.total_user}
            icon={UsersIcon}
            className="bg-purple-50"
            iconClassName="text-purple-600"
          />
        </div>

        {/* Content Section */}
        <div className="grid gap-6 lg:grid-cols-3">
          {/* Top Products Table */}
          <Card className="h-full border-none bg-white shadow-sm ring-1 ring-gray-100 lg:col-span-2">
            <CardHeader className="flex flex-col items-start gap-1 px-6 pt-6">
              <h3 className="text-lg font-bold text-gray-900">
                Produk Terlaris
              </h3>
              <p className="text-sm text-gray-500">
                Barang dengan penjualan tertinggi
              </p>
            </CardHeader>
            <CardBody className="px-6 pb-6">
              <DataTable
                showTopBar={false}
                showRowsPerPage={false}
                showSearch={false}
                columns={topProductColumns}
                data={stats?.top_selling_products || []}
                emptyContent="Belum ada data penjualan."
                totalPages={1}
                rowsPerPage={10}
                renderCell={renderTopProductCell}
                getRowKey={(item) => item.nama_barang}
              />
            </CardBody>
          </Card>

          {/* System Info Widget */}
          <Card className="h-fit border-none bg-white shadow-sm ring-1 ring-gray-100">
            <CardHeader className="flex flex-col items-start gap-1 px-6 pt-6">
              <h3 className="text-lg font-bold text-gray-900">Status Sistem</h3>
              <p className="text-sm text-gray-500">Informasi konektivitas</p>
            </CardHeader>
            <CardBody className="flex flex-col gap-4 px-6 pb-6">
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium text-gray-600">
                  Status Server
                </span>
                <span className="flex items-center gap-2 rounded-full bg-green-50 px-2 py-1 text-xs font-medium text-green-700 ring-1 ring-inset ring-green-600/20">
                  <span className="h-1.5 w-1.5 rounded-full bg-green-600"></span>
                  Online
                </span>
              </div>
              <Divider className="my-1" />
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium text-gray-600">
                  Database
                </span>
                <span className="flex items-center gap-2 rounded-full bg-green-50 px-2 py-1 text-xs font-medium text-green-700 ring-1 ring-inset ring-green-600/20">
                  <span className="h-1.5 w-1.5 rounded-full bg-green-600"></span>
                  Terhubung
                </span>
              </div>
              <Divider className="my-1" />
              <div className="flex items-center justify-between">
                <span className="text-sm font-medium text-gray-600">
                  Terakhir Update
                </span>
                <span className="text-xs text-gray-500">Baru saja</span>
              </div>
            </CardBody>
          </Card>
        </div>
      </div>
    </AppLayout>
  );
}

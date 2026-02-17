"use client";

import { useEffect, useState, useMemo, Key } from "react";
import {
  Chip,
  Tabs,
  Tab,
  Modal,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
  Button,
  useDisclosure,
} from "@nextui-org/react";
import {
  EyeIcon,
  ArrowPathIcon,
  DocumentMagnifyingGlassIcon,
} from "@heroicons/react/24/outline";
import { AppLayout } from "@/components/layout/app-layout";
import { stokApi } from "@/lib/api";
import { HistoryStok, Stok } from "@/lib/types";
import { toast } from "sonner";
import { DataTable, Column } from "@/components/table/data-table";

export default function StokPage() {
  const [historyData, setHistoryData] = useState<HistoryStok[]>([]);
  const [stockData, setStockData] = useState<Stok[]>([]);

  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");
  const [page, setPage] = useState(1);
  const [activeTab, setActiveTab] = useState("history");
  const [filteredBarangId, setFilteredBarangId] = useState<number | null>(null);

  const [selectedHistory, setSelectedHistory] = useState<HistoryStok | null>(
    null,
  );

  const [rowsPerPage, setRowsPerPage] = useState(10);
  const handleSearchChange = (value: string) => {
    setSearch(value);
    setPage(1);
  };

  const handleRowsPerPageChange = (value: number) => {
    setRowsPerPage(value);
    setPage(1);
  };

  const { isOpen, onOpen, onOpenChange } = useDisclosure();

  const [selectedStokDetail, setSelectedStokDetail] = useState<Stok | null>(
    null,
  );
  const {
    isOpen: isDetailOpen,
    onOpen: onDetailOpen,
    onOpenChange: onDetailOpenChange,
  } = useDisclosure();

  const handleViewDetail = (history: HistoryStok) => {
    setSelectedHistory(history);
    onOpen();
  };

  const handleViewStokDetail = async (stok: Stok) => {
    try {
      const detail = await stokApi.getById(stok.barang_id);
      setSelectedStokDetail(detail);
      onDetailOpen();
    } catch (error) {
      console.error("Failed to fetch stok detail", error);
      toast.error("Gagal mengambil detail stok");
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    setLoading(true);
    try {
      const [history, stocks] = await Promise.all([
        stokApi.getHistory(),
        stokApi.getAll(),
      ]);
      setHistoryData(history);
      setStockData(stocks);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const fetchHistoryByBarang = async (barangId: number) => {
    setLoading(true);
    try {
      const history = await stokApi.getHistory(barangId);
      setHistoryData(history);
      setFilteredBarangId(barangId);
      setActiveTab("history");
    } catch (err) {
      console.error(err);
      toast.error("Gagal mengambil riwayat barang");
    } finally {
      setLoading(false);
    }
  };

  const resetHistory = async () => {
    setLoading(true);
    try {
      const history = await stokApi.getHistory();
      setHistoryData(history);
      setFilteredBarangId(null);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const filteredHistory = useMemo(() => {
    const lowerSearch = search.toLowerCase();
    return historyData.filter(
      (item) =>
        item.barang?.nama_barang.toLowerCase().includes(lowerSearch) ||
        item.barang?.kode_barang?.toLowerCase().includes(lowerSearch) ||
        item.jenis_transaksi.toLowerCase().includes(lowerSearch),
    );
  }, [search, historyData]);

  const filteredStock = useMemo(() => {
    const lowerSearch = search.toLowerCase();
    return stockData.filter(
      (item) =>
        item.barang?.nama_barang.toLowerCase().includes(lowerSearch) ||
        item.barang?.kode_barang?.toLowerCase().includes(lowerSearch),
    );
  }, [search, stockData]);

  const itemsHistory = filteredHistory.slice(
    (page - 1) * rowsPerPage,
    page * rowsPerPage,
  );
  const itemsStock = filteredStock.slice(
    (page - 1) * rowsPerPage,
    page * rowsPerPage,
  );

  const stockColumns: Column[] = [
    { uid: "kode_barang", name: "KODE BARANG" },
    { uid: "nama_barang", name: "NAMA BARANG" },
    { uid: "stok_akhir", name: "SISA STOK", align: "center" },
    { uid: "updated_at", name: "TERAKHIR UPDATE" },
    { uid: "aksi", name: "AKSI" },
  ];

  const historyColumns: Column[] = [
    { uid: "created_at", name: "WAKTU" },
    { uid: "barang", name: "BARANG" },
    { uid: "jenis_transaksi", name: "JENIS" },
    { uid: "jumlah", name: "JUMLAH" },
    { uid: "stok_sebelum", name: "SEBELUM" },
    { uid: "stok_sesudah", name: "SESUDAH" },
    { uid: "keterangan", name: "KETERANGAN" },
    { uid: "aksi", name: "AKSI" },
  ];

  const renderStockCell = (item: Stok, columnKey: Key) => {
    switch (columnKey) {
      case "kode_barang":
        return (
          <span className="font-mono font-medium">
            {item.barang?.kode_barang || "-"}
          </span>
        );
      case "nama_barang":
        return item.barang?.nama_barang || "-";
      case "stok_akhir":
        return (
          <div
            className={`font-bold text-center ${
              item.stok_akhir <= 5 ? "text-red-500" : "text-green-600"
            }`}>
            {item.stok_akhir}
          </div>
        );
      case "updated_at":
        return item.updated_at
          ? new Date(item.updated_at).toLocaleString()
          : "-";
      case "aksi":
        return (
          <div className="flex items-center gap-2">
            <Button
              isIconOnly
              size="sm"
              variant="light"
              onPress={() => handleViewStokDetail(item)}>
              <EyeIcon className="h-4 w-4 text-blue-500" />
            </Button>
            <Button
              isIconOnly
              size="sm"
              variant="light"
              onPress={() => fetchHistoryByBarang(item.barang_id)}>
              <DocumentMagnifyingGlassIcon className="h-4 w-4 text-gray-500" />
            </Button>
          </div>
        );
      default:
        return null;
    }
  };

  const renderHistoryCell = (item: HistoryStok, columnKey: Key) => {
    switch (columnKey) {
      case "created_at":
        return (
          <span className="text-xs text-gray-500">
            {new Date(item.created_at).toLocaleString()}
          </span>
        );
      case "barang":
        return (
          <div className="flex flex-col">
            <span className="font-semibold">
              {item.barang?.nama_barang || "-"}
            </span>
            <span className="text-xs text-gray-400">
              {item.barang?.kode_barang || "-"}
            </span>
          </div>
        );
      case "jenis_transaksi":
        return (
          <Chip
            size="sm"
            color={item.jenis_transaksi === "masuk" ? "success" : "danger"}
            variant="flat"
            className="capitalize">
            {item.jenis_transaksi}
          </Chip>
        );
      case "jumlah":
        return (
          <span
            className={`font-bold ${
              item.jenis_transaksi === "masuk"
                ? "text-green-600"
                : "text-red-500"
            }`}>
            {item.jenis_transaksi === "masuk" ? "+" : "-"}
            {item.jumlah}
          </span>
        );
      case "stok_sebelum":
        return item.stok_sebelum;
      case "stok_sesudah":
        return item.stok_sesudah;
      case "keterangan":
        return (
          <div className="text-gray-500 text-sm max-w-xs truncate">
            {item.keterangan}
          </div>
        );
      case "aksi":
        return (
          <div className="flex items-center gap-2">
            <Button
              isIconOnly
              size="sm"
              variant="light"
              onPress={() => handleViewDetail(item)}>
              <EyeIcon className="h-4 w-4 text-gray-500" />
            </Button>
          </div>
        );
      default:
        return null;
    }
  };

  const totalPagesHistory =
    Math.ceil(filteredHistory.length / rowsPerPage) || 1;
  const totalPagesStock = Math.ceil(filteredStock.length / rowsPerPage) || 1;

  return (
    <AppLayout>
      <div className="flex flex-col gap-6">
        <div className="flex flex-col justify-between gap-4 md:flex-row md:items-end">
          <div>
            <h1 className="text-2xl font-bold">Manajemen Stok</h1>
            <p className="text-gray-500">
              Monitor sisa stok dan riwayat pergerakan barang.
            </p>
          </div>
        </div>

        <Tabs
          selectedKey={activeTab}
          aria-label="Stok Menu"
          onSelectionChange={(key) => setActiveTab(key as string)}>
          <Tab key="status" title="Status Stok">
            <DataTable
              columns={stockColumns}
              data={itemsStock}
              isLoading={loading}
              page={page}
              totalPages={totalPagesStock}
              onPageChange={setPage}
              rowsPerPage={rowsPerPage}
              onRowsPerPageChange={handleRowsPerPageChange}
              search={search}
              onSearchChange={handleSearchChange}
              renderCell={renderStockCell}
              emptyContent="Data stok tidak ditemukan."
            />
          </Tab>

          <Tab key="history" title="Riwayat Pergerakan">
            <div className="flex flex-col gap-4">
              {filteredBarangId && (
                <div className="flex items-center justify-between rounded-lg border border-blue-100 bg-blue-50 px-4 py-3 text-blue-700">
                  <div className="flex items-center gap-2">
                    <EyeIcon className="h-5 w-5" />
                    <span className="font-medium">
                      Menampilkan riwayat untuk barang ID: {filteredBarangId}
                    </span>
                  </div>
                  <Button
                    size="sm"
                    variant="flat"
                    color="primary"
                    startContent={<ArrowPathIcon className="h-4 w-4" />}
                    onPress={resetHistory}>
                    Tampilkan Semua
                  </Button>
                </div>
              )}
              <DataTable
                columns={historyColumns}
                data={itemsHistory}
                isLoading={loading}
                page={page}
                totalPages={totalPagesHistory}
                onPageChange={setPage}
                rowsPerPage={rowsPerPage}
                onRowsPerPageChange={handleRowsPerPageChange}
                search={search}
                onSearchChange={handleSearchChange}
                renderCell={renderHistoryCell}
                emptyContent="Belum ada riwayat stok."
              />
            </div>
          </Tab>
        </Tabs>

        {/* Detail Modal History */}
        <Modal isOpen={isOpen} onOpenChange={onOpenChange}>
          <ModalContent>
            {(onClose) => (
              <>
                <ModalHeader className="flex flex-col gap-1">
                  Detail Riwayat Stok
                </ModalHeader>
                <ModalBody>
                  {selectedHistory && (
                    <div className="space-y-4">
                      <div className="p-4 bg-slate-50 rounded-lg border border-slate-100">
                        <div className="flex justify-between items-start mb-2">
                          <h3 className="font-semibold text-lg">
                            {selectedHistory.barang?.nama_barang}
                          </h3>
                          <Chip
                            size="sm"
                            color={
                              selectedHistory.jenis_transaksi === "masuk"
                                ? "success"
                                : "warning"
                            }
                            variant="flat">
                            {selectedHistory.jenis_transaksi.toUpperCase()}
                          </Chip>
                        </div>
                        <p className="text-sm text-slate-500 font-mono">
                          {selectedHistory.barang?.kode_barang}
                        </p>
                      </div>

                      <div className="grid grid-cols-2 gap-4">
                        <div className="space-y-1">
                          <p className="text-xs text-slate-500 uppercase">
                            Jumlah
                          </p>
                          <p
                            className={`font-semibold text-lg ${selectedHistory.jenis_transaksi === "masuk" ? "text-green-600" : "text-amber-600"}`}>
                            {selectedHistory.jenis_transaksi === "masuk"
                              ? "+"
                              : "-"}
                            {selectedHistory.jumlah}
                          </p>
                        </div>
                        <div className="space-y-1">
                          <p className="text-xs text-slate-500 uppercase">
                            Tanggal
                          </p>
                          <p className="text-sm font-medium">
                            {new Date(
                              selectedHistory.created_at,
                            ).toLocaleString()}
                          </p>
                        </div>
                        <div className="space-y-1">
                          <p className="text-xs text-slate-500 uppercase">
                            Stok Sebelum
                          </p>
                          <p className="text-sm font-medium">
                            {selectedHistory.stok_sebelum}
                          </p>
                        </div>
                        <div className="space-y-1">
                          <p className="text-xs text-slate-500 uppercase">
                            Stok Sesudah
                          </p>
                          <p className="text-sm font-medium">
                            {selectedHistory.stok_sesudah}
                          </p>
                        </div>
                      </div>

                      <div className="pt-2 border-t border-slate-100">
                        <p className="text-xs text-slate-500 uppercase mb-1">
                          Keterangan / ID Transaksi
                        </p>
                        <p className="text-sm text-slate-700">
                          {selectedHistory.keterangan || "-"}
                        </p>
                      </div>

                      {selectedHistory.user && (
                        <div className="pt-2 border-t border-slate-100">
                          <p className="text-xs text-slate-500 uppercase mb-1">
                            Petugas
                          </p>
                          <p className="text-sm font-medium">
                            {selectedHistory.user.username} (
                            {selectedHistory.user.role})
                          </p>
                        </div>
                      )}
                    </div>
                  )}
                </ModalBody>
                <ModalFooter>
                  <Button color="primary" onPress={onClose}>
                    Tutup
                  </Button>
                </ModalFooter>
              </>
            )}
          </ModalContent>
        </Modal>

        {/* Detail Modal Stok */}
        <Modal isOpen={isDetailOpen} onOpenChange={onDetailOpenChange}>
          <ModalContent>
            {(onClose) => (
              <>
                <ModalHeader className="flex flex-col gap-1">
                  Detail Stok Barang
                </ModalHeader>
                <ModalBody>
                  {selectedStokDetail && (
                    <div className="space-y-4">
                      <div className="p-4 bg-slate-50 rounded-lg border border-slate-100">
                        <h3 className="font-semibold text-lg text-slate-800">
                          {selectedStokDetail.barang?.nama_barang}
                        </h3>
                        <p className="text-sm text-slate-500 font-mono">
                          {selectedStokDetail.barang?.kode_barang}
                        </p>
                      </div>

                      <div className="grid grid-cols-2 gap-4">
                        <div className="p-3 bg-blue-50 rounded-lg">
                          <p className="text-xs text-blue-500 uppercase font-semibold mb-1">
                            Sisa Stok
                          </p>
                          <p className="text-2xl font-bold text-blue-700">
                            {selectedStokDetail.stok_akhir}{" "}
                            <span className="text-sm font-normal text-blue-600">
                              {selectedStokDetail.barang?.satuan}
                            </span>
                          </p>
                        </div>
                        <div className="p-3 bg-slate-50 rounded-lg">
                          <p className="text-xs text-slate-500 uppercase font-semibold mb-1">
                            Harga Jual
                          </p>
                          <p className="text-lg font-bold text-slate-700">
                            Rp{" "}
                            {selectedStokDetail.barang?.harga_jual.toLocaleString(
                              "id-ID",
                            )}
                          </p>
                        </div>
                      </div>

                      <div className="space-y-3 pt-2 text-sm">
                        <div className="flex justify-between border-b border-slate-100 pb-2">
                          <span className="text-slate-500">ID Stok</span>
                          <span className="font-medium text-slate-700">
                            #{selectedStokDetail.id}
                          </span>
                        </div>
                        <div className="flex justify-between border-b border-slate-100 pb-2">
                          <span className="text-slate-500">
                            Terakhir Update
                          </span>
                          <span className="font-medium text-slate-700">
                            {selectedStokDetail.updated_at
                              ? new Date(
                                  selectedStokDetail.updated_at,
                                ).toLocaleString()
                              : "-"}
                          </span>
                        </div>
                        {selectedStokDetail.barang?.deskripsi && (
                          <div className="pt-1">
                            <span className="block text-slate-500 mb-1">
                              Deskripsi
                            </span>
                            <p className="text-slate-700 p-2 bg-slate-50 rounded border border-slate-100">
                              {selectedStokDetail.barang.deskripsi}
                            </p>
                          </div>
                        )}
                      </div>
                    </div>
                  )}
                </ModalBody>
                <ModalFooter>
                  <Button color="primary" onPress={onClose}>
                    Tutup
                  </Button>
                </ModalFooter>
              </>
            )}
          </ModalContent>
        </Modal>
      </div>
    </AppLayout>
  );
}

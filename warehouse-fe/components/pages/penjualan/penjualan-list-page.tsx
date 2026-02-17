"use client";

import { useEffect, useMemo, useState, Key } from "react";
import {
  Button,
  Chip,
  Spinner,
  Modal,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
  useDisclosure,
  DatePicker,
  Tooltip,
} from "@nextui-org/react";
import {
  PlusIcon,
  EyeIcon,
  ArrowDownTrayIcon,
} from "@heroicons/react/24/outline";
import { AppLayout } from "@/components/layout/app-layout";
import { penjualanApi } from "@/lib/api";
import { JualHeader } from "@/lib/types";
import { useRouter } from "next/navigation";
import { DataTable, Column } from "@/components/table/data-table";
import { downloadPenjualanInvoice } from "@/lib/invoice-pdf";

export default function PenjualanListPage() {
  const router = useRouter();
  const [data, setData] = useState<JualHeader[]>([]);
  const [loading, setLoading] = useState(true);
  const [detailLoading, setDetailLoading] = useState(false);
  const [invoiceLoadingId, setInvoiceLoadingId] = useState<number | null>(null);
  const [selectedTx, setSelectedTx] = useState<JualHeader | null>(null);

  // Filters
  const [search, setSearch] = useState("");
  const [page, setPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [startDate, setStartDate] = useState<any>(null);
  const [endDate, setEndDate] = useState<any>(null);

  const { isOpen, onOpen, onOpenChange } = useDisclosure();

  useEffect(() => {
    fetchData();
  }, [startDate, endDate]);

  const fetchData = async () => {
    setLoading(true);
    try {
      const params: any = {};
      if (startDate) params.start_date = startDate.toString();
      if (endDate) params.end_date = endDate.toString();

      const res = await penjualanApi.getAll(params);
      setData(res);
    } catch (e) {
      console.error(e);
    } finally {
      setLoading(false);
    }
  };

  const filteredData = data.filter(
    (item) =>
      item.no_faktur.toLowerCase().includes(search.toLowerCase()) ||
      item.customer.toLowerCase().includes(search.toLowerCase()),
  );

  const totalPages = Math.ceil(filteredData.length / rowsPerPage) || 1;
  const items = useMemo(
    () => filteredData.slice((page - 1) * rowsPerPage, page * rowsPerPage),
    [filteredData, page, rowsPerPage],
  );

  const columns: Column[] = [
    { uid: "no_faktur", name: "NO FAKTUR" },
    { uid: "created_at", name: "TANGGAL" },
    { uid: "customer", name: "CUSTOMER" },
    { uid: "total", name: "TOTAL" },
    { uid: "status", name: "STATUS" },
    { uid: "actions", name: "ACTIONS", align: "center" },
  ];

  const renderCell = (item: JualHeader, columnKey: Key) => {
    switch (columnKey) {
      case "no_faktur":
        return <span className="font-mono">{item.no_faktur}</span>;
      case "created_at":
        return new Date(item.created_at).toLocaleDateString();
      case "customer":
        return item.customer;
      case "total":
        return `Rp ${item.total.toLocaleString()}`;
      case "status":
        return (
          <Chip
            size="sm"
            color={item.status === "completed" ? "success" : "warning"}
            variant="flat">
            {item.status}
          </Chip>
        );
      case "actions":
        return (
          <div className="flex items-center justify-center gap-1">
            <Tooltip content="Detail">
              <Button
                isIconOnly
                size="sm"
                variant="light"
                onPress={() => handleView(item)}>
                <EyeIcon className="h-4 w-4" />
              </Button>
            </Tooltip>
            <Tooltip content="Download Invoice (PDF)">
              <Button
                isIconOnly
                size="sm"
                variant="light"
                isDisabled={invoiceLoadingId === item.id}
                onPress={() => handleDownloadInvoice(item)}>
                <ArrowDownTrayIcon className="h-4 w-4" />
              </Button>
            </Tooltip>
          </div>
        );
      default:
        return null;
    }
  };

  const handleView = async (item: JualHeader) => {
    setSelectedTx(item);
    onOpen();
    setDetailLoading(true);
    try {
      const detail = await penjualanApi.getById(item.id);
      setSelectedTx(detail);
    } catch (error) {
      console.error("Failed to fetch detail", error);
    } finally {
      setDetailLoading(false);
    }
  };

  const handleDownloadInvoice = async (item: JualHeader) => {
    setInvoiceLoadingId(item.id);
    try {
      const detail = await penjualanApi.getById(item.id);
      await downloadPenjualanInvoice(detail);
    } catch (error) {
      console.error("Failed to download invoice", error);
    } finally {
      setInvoiceLoadingId(null);
    }
  };

  return (
    <AppLayout>
      <div className="flex flex-col gap-6">
        {/* Header & Controls */}
        <div className="flex flex-col justify-between gap-4 md:flex-row md:items-end">
          <div>
            <h1 className="text-2xl font-bold">Riwayat Penjualan</h1>
            <p className="text-gray-500">
              Daftar transaksi penjualan barang keluar.
            </p>
          </div>
        </div>

        <DataTable
          columns={columns}
          data={items}
          isLoading={loading}
          emptyContent="Belum ada transaksi."
          page={page}
          totalPages={totalPages}
          onPageChange={setPage}
          rowsPerPage={rowsPerPage}
          onRowsPerPageChange={(v) => {
            setRowsPerPage(v);
            setPage(1);
          }}
          search={search}
          onSearchChange={(v) => {
            setSearch(v);
            setPage(1);
          }}
          searchPlaceholder="Cari faktur / customer..."
          renderCell={renderCell}
          filtersContent={
            <div className="flex gap-2 rounded">
              <DatePicker
                variant="bordered"
                size="md"
                value={startDate}
                onChange={(v) => {
                  setStartDate(v);
                  setPage(1);
                }}
                className="w-full sm:w-44"
                classNames={{
                  inputWrapper: "!rounded-2xl",
                }}
              />
              <DatePicker
                variant="bordered"
                size="md"
                value={endDate}
                onChange={(v) => {
                  setEndDate(v);
                  setPage(1);
                }}
                className="w-full sm:w-44"
                classNames={{
                  inputWrapper: "!rounded-2xl",
                }}
              />
            </div>
          }
          topContent={
            <Button
              color="primary"
              size="md"
              className="flex-shrink-0 font-medium"
              endContent={<PlusIcon className="h-4 w-4 font-bold" />}
              onPress={() => router.push("/penjualan/create")}>
              Tambah Baru
            </Button>
          }
        />

        {/* Detail Modal */}
        <Modal isOpen={isOpen} onOpenChange={onOpenChange} size="3xl">
          <ModalContent>
            {(onClose) => (
              <>
                <ModalHeader className="flex flex-col gap-1">
                  Detail Transaksi {selectedTx?.no_faktur}
                </ModalHeader>
                <ModalBody>
                  <div className="flex justify-between mb-4 text-sm">
                    <div>
                      <p className="text-gray-500">Customer</p>
                      <p className="font-semibold">{selectedTx?.customer}</p>
                    </div>
                    <div className="text-right">
                      <p className="text-gray-500">Tanggal</p>
                      <p className="font-semibold">
                        {selectedTx?.created_at &&
                          new Date(selectedTx.created_at).toLocaleString()}
                      </p>
                    </div>
                  </div>

                  {detailLoading ? (
                    <div className="flex w-full justify-center p-8">
                      <Spinner label="Memuat detail transaksi..." />
                    </div>
                  ) : (
                    <div className="border rounded-lg overflow-hidden">
                      <table className="w-full text-sm text-left">
                        <thead className="bg-gray-50 text-gray-700">
                          <tr>
                            <th className="px-4 py-2">Barang</th>
                            <th className="px-4 py-2 text-right">Qty</th>
                            <th className="px-4 py-2 text-right">Harga</th>
                            <th className="px-4 py-2 text-right">Subtotal</th>
                          </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-100">
                          {selectedTx?.details?.map((det) => (
                            <tr key={det.id}>
                              <td className="px-4 py-2">
                                {det.barang?.nama_barang} (
                                {det.barang?.kode_barang})
                              </td>
                              <td className="px-4 py-2 text-right">
                                {det.qty}
                              </td>
                              <td className="px-4 py-2 text-right">
                                Rp {det.harga.toLocaleString()}
                              </td>
                              <td className="px-4 py-2 text-right font-medium">
                                Rp {det.subtotal.toLocaleString()}
                              </td>
                            </tr>
                          ))}
                        </tbody>
                        <tfoot className="bg-gray-50 font-semibold">
                          <tr>
                            <td colSpan={3} className="px-4 py-2 text-right">
                              Total
                            </td>
                            <td className="px-4 py-2 text-right text-green-600">
                              Rp {selectedTx?.total.toLocaleString()}
                            </td>
                          </tr>
                        </tfoot>
                      </table>
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

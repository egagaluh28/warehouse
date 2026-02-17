"use client";

import { useEffect, useMemo, useState, Key } from "react";
import { useRouter } from "next/navigation";
import {
  Table,
  TableHeader,
  TableColumn,
  TableBody,
  TableRow,
  TableCell,
  Button,
  Chip,
  Spinner,
  useDisclosure,
  Modal,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
  DatePicker,
  Tooltip,
} from "@nextui-org/react";
import {
  PlusIcon,
  EyeIcon,
  ArrowDownTrayIcon,
} from "@heroicons/react/24/outline";
import { AppLayout } from "@/components/layout/app-layout";
import { pembelianApi } from "@/lib/api";
import { BeliHeader } from "@/lib/types";
import { DataTable, Column } from "@/components/table/data-table";
import { downloadPembelianInvoice } from "@/lib/invoice-pdf";

export default function PembelianListPage() {
  const router = useRouter();
  const [data, setData] = useState<BeliHeader[]>([]);
  const [loading, setLoading] = useState(true);
  const [detailLoading, setDetailLoading] = useState(false);
  const [invoiceLoadingId, setInvoiceLoadingId] = useState<number | null>(null);
  const [selectedTx, setSelectedTx] = useState<BeliHeader | null>(null);

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

      const res = await pembelianApi.getAll(params);
      setData(res);
    } catch (error) {
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  const filteredData = data.filter(
    (item) =>
      item.no_faktur.toLowerCase().includes(search.toLowerCase()) ||
      item.supplier.toLowerCase().includes(search.toLowerCase()),
  );

  const totalPages = Math.ceil(filteredData.length / rowsPerPage) || 1;
  const items = useMemo(
    () => filteredData.slice((page - 1) * rowsPerPage, page * rowsPerPage),
    [filteredData, page, rowsPerPage],
  );

  const columns: Column[] = [
    { uid: "no_faktur", name: "NO FAKTUR" },
    { uid: "supplier", name: "SUPPLIER" },
    { uid: "created_at", name: "TANGGAL" },
    { uid: "total", name: "TOTAL" },
    { uid: "status", name: "STATUS" },
    { uid: "actions", name: "ACTIONS", align: "center" },
  ];

  const renderCell = (item: BeliHeader, columnKey: Key) => {
    switch (columnKey) {
      case "no_faktur":
        return <span className="font-mono">{item.no_faktur}</span>;
      case "supplier":
        return item.supplier;
      case "created_at":
        return new Date(item.created_at).toLocaleDateString();
      case "total":
        return `Rp ${item.total.toLocaleString()}`;
      case "status":
        return (
          <Chip
            size="sm"
            color={item.status === "completed" ? "primary" : "warning"}
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

  const handleView = async (item: BeliHeader) => {
    setSelectedTx(item);
    onOpen();
    setDetailLoading(true);
    try {
      const detail = await pembelianApi.getById(item.id);
      setSelectedTx(detail);
    } catch (err) {
      console.error(err);
    } finally {
      setDetailLoading(false);
    }
  };

  const handleDownloadInvoice = async (item: BeliHeader) => {
    setInvoiceLoadingId(item.id);
    try {
      const detail = await pembelianApi.getById(item.id);
      await downloadPembelianInvoice(detail);
    } catch (err) {
      console.error(err);
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
            <h1 className="text-2xl font-bold">Riwayat Pembelian</h1>
            <p className="text-gray-500">
              Daftar transaksi pembelian (Restock) barang masuk.
            </p>
          </div>
        </div>

        <DataTable
          columns={columns}
          data={items}
          isLoading={loading}
          emptyContent="Belum ada data pembelian."
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
          searchPlaceholder="Cari faktur / supplier..."
          renderCell={renderCell}
          filtersContent={
            <div className="flex gap-2">
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
              startContent={<PlusIcon className="h-4 w-4 font-bold" />}
              onPress={() => router.push("/pembelian/create")}
              className="font-medium">
              Tambah Baru
            </Button>
          }
        />

        {/* Detail Modal */}
        <Modal isOpen={isOpen} onOpenChange={onOpenChange} size="3xl">
          <ModalContent>
            {(onClose) => (
              <>
                <ModalHeader>
                  Detail Pembelian {selectedTx?.no_faktur}
                </ModalHeader>
                <ModalBody>
                  <div className="flex justify-between mb-4 text-sm">
                    <div>
                      <p className="text-gray-500">Supplier</p>
                      <p className="font-semibold">{selectedTx?.supplier}</p>
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
                      <Spinner label="Memuat detail pembelian..." />
                    </div>
                  ) : (
                    <Table removeWrapper aria-label="Items">
                      <TableHeader>
                        <TableColumn>BARANG</TableColumn>
                        <TableColumn align="end">QTY</TableColumn>
                        <TableColumn align="end">HARGA BELI</TableColumn>
                        <TableColumn align="end">SUBTOTAL</TableColumn>
                      </TableHeader>
                      <TableBody>
                        {(selectedTx?.details || []).map((det) => (
                          <TableRow key={det.id}>
                            <TableCell>{det.barang?.nama_barang}</TableCell>
                            <TableCell>{det.qty}</TableCell>
                            <TableCell>
                              Rp {det.harga.toLocaleString()}
                            </TableCell>
                            <TableCell>
                              Rp {det.subtotal.toLocaleString()}
                            </TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>
                  )}
                </ModalBody>
                <ModalFooter>
                  <div className="flex w-full justify-between items-center">
                    <span className="font-bold">
                      Total: Rp {selectedTx?.total.toLocaleString()}
                    </span>
                    <Button onPress={onClose}>Tutup</Button>
                  </div>
                </ModalFooter>
              </>
            )}
          </ModalContent>
        </Modal>
      </div>
    </AppLayout>
  );
}

"use client";

import { useEffect, useState, useCallback, useMemo, Key } from "react";
import {
  Input,
  Button,
  Tooltip,
  Spinner,
  useDisclosure,
  Modal,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
  Chip,
  SortDescriptor,
  Dropdown,
  DropdownTrigger,
  DropdownMenu,
  DropdownItem,
} from "@nextui-org/react";
import { AppLayout } from "@/components/layout/app-layout";
import { barangApi } from "@/lib/api";
import { Barang, CreateBarangRequest } from "@/lib/types";
import { toast } from "sonner";
import { DataTable, Column } from "@/components/table/data-table";

import {
  PlusIcon,
  PencilIcon,
  TrashIcon,
  EyeIcon,
  EllipsisVerticalIcon,
} from "@heroicons/react/24/outline";

export default function BarangPage() {
  // State
  const [data, setData] = useState<Barang[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isDetailLoading, setIsDetailLoading] = useState(false);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [search, setSearch] = useState("");
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [selectedBarang, setSelectedBarang] = useState<Barang | null>(null);
  const [sortDescriptor, setSortDescriptor] = useState<SortDescriptor>({
    column: "kode_barang",
    direction: "ascending",
  });

  // Handlers
  const {
    isOpen: isFormOpen,
    onOpen: onFormOpen,
    onOpenChange: onFormOpenChange,
    onClose: onFormClose,
  } = useDisclosure();

  const {
    isOpen: isDeleteOpen,
    onOpen: onDeleteOpen,
    onOpenChange: onDeleteOpenChange,
    onClose: onDeleteClose,
  } = useDisclosure();

  const {
    isOpen: isDetailOpen,
    onOpen: onDetailOpen,
    onOpenChange: onDetailOpenChange,
  } = useDisclosure();

  const [formData, setFormData] = useState<CreateBarangRequest>({
    nama_barang: "",
    deskripsi: "",
    satuan: "pcs",
    harga_beli: 0,
    harga_jual: 0,
  });

  // Fetch Data
  const fetchData = useCallback(async () => {
    setIsLoading(true);
    try {
      const res = await barangApi.getAll({
        page,
        limit: rowsPerPage,
        search,
        sort_by: sortDescriptor.column as string,
        order: sortDescriptor.direction === "ascending" ? "asc" : "desc",
      });
      setData(res.data || []);
      setTotalPages(res.meta?.total_pages || 1);
    } catch (error) {
      toast.error("Gagal memuat data barang");
      console.error("Failed to fetch barang", error);
    } finally {
      setIsLoading(false);
    }
  }, [page, rowsPerPage, search, sortDescriptor]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  // Handlers
  const handleSortChange = (descriptor: SortDescriptor) => {
    setSortDescriptor(descriptor);
    setPage(1);
  };

  const handleSearch = (value: string) => {
    setSearch(value);
    setPage(1);
  };

  const resetForm = () => {
    setFormData({
      nama_barang: "",
      deskripsi: "",
      satuan: "pcs",
      harga_beli: 0,
      harga_jual: 0,
    });
    setSelectedBarang(null);
  };

  const handleCreateOpen = () => {
    resetForm();
    onFormOpen();
  };

  const handleEditOpen = (item: Barang) => {
    setSelectedBarang(item);
    setFormData({
      nama_barang: item.nama_barang,
      deskripsi: item.deskripsi || "",
      satuan: item.satuan,
      harga_beli: item.harga_beli,
      harga_jual: item.harga_jual,
    });
    onFormOpen();
  };

  const handleDeleteOpen = (item: Barang) => {
    setSelectedBarang(item);
    onDeleteOpen();
  };

  const handleDetailOpen = async (item: Barang) => {
    setSelectedBarang(item);
    onDetailOpen();

    setIsDetailLoading(true);
    try {
      const detail = await barangApi.getById(item.id);
      setSelectedBarang(detail);
    } catch (error) {
      toast.error("Gagal memuat detail barang");
      console.error("Failed to fetch barang detail", error);
    } finally {
      setIsDetailLoading(false);
    }
  };

  const handleSubmit = async () => {
    if (!formData.nama_barang.trim()) {
      toast.error("Nama barang wajib diisi");
      return;
    }
    if (!formData.harga_beli || formData.harga_beli <= 0) {
      toast.error("Harga beli wajib diisi");
      return;
    }
    if (!formData.harga_jual || formData.harga_jual <= 0) {
      toast.error("Harga jual wajib diisi");
      return;
    }

    try {
      if (selectedBarang) {
        await barangApi.update(selectedBarang.id, formData);
        toast.success("Barang berhasil diperbarui");
      } else {
        await barangApi.create(formData);
        toast.success("Barang berhasil ditambahkan");
      }
      fetchData();
      onFormClose();
    } catch (error) {
      toast.error("Gagal menyimpan data barang");
      console.error("Failed to save data", error);
    }
  };

  const handleDelete = async () => {
    if (!selectedBarang) return;
    try {
      await barangApi.delete(selectedBarang.id);
      toast.success("Barang berhasil dihapus");
      fetchData();
      onDeleteClose();
    } catch (error) {
      toast.error("Gagal menghapus barang");
      console.error("Failed to delete barang", error);
    }
  };

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
    }).format(value);
  };

  // Table Config
  const columns: Column[] = useMemo(
    () => [
      { name: "KODE", uid: "kode_barang", sortable: true },
      { name: "NAMA BARANG", uid: "nama_barang", sortable: true },
      { name: "HARGA BELI", uid: "harga_beli", sortable: true, align: "end" },
      { name: "HARGA JUAL", uid: "harga_jual", sortable: true, align: "end" },
      { name: "AKSI", uid: "actions", sortable: false, align: "center" },
    ],
    [],
  );

  const renderCell = useCallback((item: Barang, columnKey: Key) => {
    const cellValue = item[columnKey as keyof Barang];

    switch (columnKey) {
      case "kode_barang":
        return <p className="font-mono text-sm">{item.kode_barang}</p>;
      case "nama_barang":
        return (
          <div className="flex flex-col">
            <p className="text-bold text-sm capitalize">{item.nama_barang}</p>
            <p className="text-bold text-tiny text-default-400">
              {item.deskripsi || "-"}
            </p>
          </div>
        );
      case "harga_beli":
      case "harga_jual":
        return <p className="text-sm">{formatCurrency(Number(cellValue))}</p>;
      case "actions":
        return (
          <div className="relative flex items-center justify-center gap-2">
            <Tooltip content="Detail">
              <span
                className="cursor-pointer text-lg text-default-400 active:opacity-50 hover:text-primary"
                onClick={() => handleDetailOpen(item)}>
                <EyeIcon className="h-5 w-5" />
              </span>
            </Tooltip>
            <Tooltip content="Edit">
              <span
                className="cursor-pointer text-lg text-default-400 active:opacity-50 hover:text-warning"
                onClick={() => handleEditOpen(item)}>
                <PencilIcon className="h-5 w-5" />
              </span>
            </Tooltip>
            <Tooltip color="danger" content="Hapus">
              <span
                className="cursor-pointer text-lg text-danger active:opacity-50 hover:text-danger-400"
                onClick={() => handleDeleteOpen(item)}>
                <TrashIcon className="h-5 w-5" />
              </span>
            </Tooltip>
          </div>
        );
      default:
        return cellValue;
    }
  }, []);

  return (
    <AppLayout>
      <div className="flex flex-col gap-6">
        {/* Header Section */}
        <div className="flex flex-col justify-between gap-4 md:flex-row md:items-end bg-content1 p-6 rounded-large shadow-sm">
          <div className="flex flex-col gap-1">
            <h1 className="text-2xl font-bold tracking-tight">Daftar Barang</h1>
            <p className="text-sm text-default-500">
              Kelola inventaris dan stok barang gudang.
            </p>
          </div>
        </div>

        {/* Table Section */}
        <DataTable
          columns={columns}
          data={data}
          isLoading={isLoading}
          emptyContent="Tidak ada data barang."
          page={page}
          totalPages={totalPages}
          onPageChange={setPage}
          rowsPerPage={rowsPerPage}
          onRowsPerPageChange={(v) => {
            setRowsPerPage(v);
            setPage(1);
          }}
          search={search}
          onSearchChange={handleSearch}
          searchPlaceholder="Cari nama / kode..."
          sortDescriptor={sortDescriptor}
          onSortChange={handleSortChange}
          renderCell={renderCell}
          topContent={
            <Button
              color="primary"
              size="md"
              startContent={<PlusIcon className="h-4 w-4" />}
              onPress={handleCreateOpen}
              className="font-medium">
              Tambah Barang
            </Button>
          }
        />

        {/* Modal Form */}
        <Modal
          isOpen={isFormOpen}
          onOpenChange={onFormOpenChange}
          placement="top-center"
          backdrop="opaque">
          <ModalContent className="bg-white dark:bg-zinc-900 text-black dark:text-white">
            {(onClose) => (
              <>
                <ModalHeader className="flex flex-col gap-1 text-inherit">
                  {selectedBarang ? "Edit Barang" : "Tambah Barang Baru"}
                </ModalHeader>
                <ModalBody>
                  <Input
                    autoFocus
                    label="Nama Barang"
                    placeholder="Contoh: Kertas A4"
                    variant="bordered"
                    value={formData.nama_barang}
                    onValueChange={(v) =>
                      setFormData({ ...formData, nama_barang: v })
                    }
                  />
                  <Input
                    label="Deskripsi"
                    placeholder="Keterangan tambahan"
                    variant="bordered"
                    value={formData.deskripsi}
                    onValueChange={(v) =>
                      setFormData({ ...formData, deskripsi: v })
                    }
                  />
                  <div className="flex gap-2">
                    <Input
                      label="Harga Beli"
                      placeholder="0"
                      type="number"
                      variant="bordered"
                      value={formData.harga_beli.toString()}
                      onValueChange={(v) =>
                        setFormData({ ...formData, harga_beli: Number(v) })
                      }
                      startContent={
                        <div className="pointer-events-none flex items-center">
                          <span className="text-default-400 text-small">
                            Rp
                          </span>
                        </div>
                      }
                    />
                    <Input
                      label="Harga Jual"
                      placeholder="0"
                      type="number"
                      variant="bordered"
                      value={formData.harga_jual.toString()}
                      onValueChange={(v) =>
                        setFormData({ ...formData, harga_jual: Number(v) })
                      }
                      startContent={
                        <div className="pointer-events-none flex items-center">
                          <span className="text-default-400 text-small">
                            Rp
                          </span>
                        </div>
                      }
                    />
                  </div>
                  <Input
                    label="Satuan"
                    placeholder="pcs, rim, box"
                    variant="bordered"
                    value={formData.satuan}
                    onValueChange={(v) =>
                      setFormData({ ...formData, satuan: v })
                    }
                  />
                </ModalBody>
                <ModalFooter>
                  <Button color="danger" variant="flat" onPress={onClose}>
                    Batal
                  </Button>
                  <Button color="primary" onPress={handleSubmit}>
                    Simpan
                  </Button>
                </ModalFooter>
              </>
            )}
          </ModalContent>
        </Modal>

        {/* Modal Detail */}
        <Modal
          isOpen={isDetailOpen}
          onOpenChange={onDetailOpenChange}
          placement="top-center"
          backdrop="opaque">
          <ModalContent>
            {(onClose) => (
              <>
                <ModalHeader className="flex flex-col gap-1">
                  Detail Barang
                </ModalHeader>
                <ModalBody>
                  {isDetailLoading ? (
                    <div className="flex items-center justify-center py-8">
                      <Spinner size="md" />
                    </div>
                  ) : (
                    <div className="grid grid-cols-2 gap-4">
                      <div>
                        <p className="text-sm text-default-500">Kode Barang</p>
                        <p className="font-semibold">
                          {selectedBarang?.kode_barang || "-"}
                        </p>
                      </div>
                      <div>
                        <p className="text-sm text-default-500">Nama Barang</p>
                        <p className="font-semibold">
                          {selectedBarang?.nama_barang || "-"}
                        </p>
                      </div>
                      <div className="col-span-2">
                        <p className="text-sm text-default-500">Deskripsi</p>
                        <p className="text-sm">
                          {selectedBarang?.deskripsi || "-"}
                        </p>
                      </div>
                      <div>
                        <p className="text-sm text-default-500">Harga Beli</p>
                        <p className="font-semibold text-primary">
                          {formatCurrency(selectedBarang?.harga_beli || 0)}
                        </p>
                      </div>
                      <div>
                        <p className="text-sm text-default-500">Harga Jual</p>
                        <p className="font-semibold text-success">
                          {formatCurrency(selectedBarang?.harga_jual || 0)}
                        </p>
                      </div>
                      <div>
                        <p className="text-sm text-default-500">
                          Stok Saat Ini
                        </p>
                        <Chip
                          size="sm"
                          color={
                            (selectedBarang?.stok || 0) > 0
                              ? "success"
                              : "danger"
                          }
                          variant="flat">
                          {selectedBarang?.stok || 0} {selectedBarang?.satuan}
                        </Chip>
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

        {/* Modal Delete */}
        <Modal
          isOpen={isDeleteOpen}
          onOpenChange={onDeleteOpenChange}
          backdrop="opaque">
          <ModalContent>
            {(onClose) => (
              <>
                <ModalHeader className="flex flex-col gap-1">
                  Konfirmasi Hapus
                </ModalHeader>
                <ModalBody>
                  <p>
                    Apakah anda yakin ingin menghapus barang{" "}
                    <span className="font-bold">
                      {selectedBarang?.nama_barang}
                    </span>
                    ?
                  </p>
                  <p className="text-sm text-default-500">
                    Tindakan ini tidak dapat dibatalkan.
                  </p>
                </ModalBody>
                <ModalFooter>
                  <Button color="default" variant="light" onPress={onClose}>
                    Batal
                  </Button>
                  <Button color="danger" onPress={handleDelete}>
                    Hapus
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

"use client";

import { useEffect, useState, useCallback, Key } from "react";
import { useRouter } from "next/navigation";
import {
  Input,
  Button,
  Spinner,
  useDisclosure,
  Modal,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
  Chip,
  Select,
  SelectItem,
} from "@nextui-org/react";
import { PlusIcon } from "@heroicons/react/24/outline";
import { AppLayout } from "@/components/layout/app-layout";
import { userApi, authApi } from "@/lib/api";
import { User, RegisterRequest } from "@/lib/types";
import { useAuth } from "@/components/providers/auth-provider";
import { DataTable, Column } from "@/components/table/data-table";

export default function UserPage() {
  const { user, isLoading: authLoading } = useAuth();
  const router = useRouter();

  // Redirect if not admin
  useEffect(() => {
    if (!authLoading && (!user || user.role !== "admin")) {
      router.push("/dashboard");
    }
  }, [user, authLoading, router]);

  // State
  const [data, setData] = useState<User[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [search, setSearch] = useState("");
  const [page, setPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(10);

  // Modal State
  const {
    isOpen: isFormOpen,
    onOpen: onFormOpen,
    onOpenChange: onFormOpenChange,
    onClose: onFormClose,
  } = useDisclosure();

  const [formData, setFormData] = useState<RegisterRequest>({
    username: "",
    password: "",
    email: "",
    full_name: "",
    role: "staff",
  });

  // Fetch Data
  const fetchData = useCallback(async () => {
    setIsLoading(true);
    try {
      // Assuming GetAll returns all users without pagination params initially
      // Or filter client side if backend doesn't support pagination yet for users
      const users = await userApi.getAll();
      setData(users);
    } catch (error) {
      console.error("Failed to fetch users", error);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    if (user?.role === "admin") {
      fetchData();
    }
  }, [fetchData, user]);

  // Handlers
  const handleSearch = (value: string) => {
    setSearch(value);
    setPage(1);
  };

  const filteredData = data.filter(
    (item) =>
      item.username.toLowerCase().includes(search.toLowerCase()) ||
      item.full_name.toLowerCase().includes(search.toLowerCase()) ||
      item.email.toLowerCase().includes(search.toLowerCase()),
  );

  const totalPages = Math.ceil(filteredData.length / rowsPerPage) || 1;
  const items = filteredData.slice(
    (page - 1) * rowsPerPage,
    page * rowsPerPage,
  );

  const resetForm = () => {
    setFormData({
      username: "",
      password: "",
      email: "",
      full_name: "",
      role: "staff",
    });
  };

  const handleCreateOpen = () => {
    resetForm();
    onFormOpen();
  };

  const handleSubmit = async () => {
    try {
      await authApi.register(formData);
      fetchData();
      onFormClose();
    } catch (error) {
      console.error("Failed to create user", error);
      alert("Gagal menambahkan user. Pastikan username belum digunakan.");
    }
  };

  // Table Config
  const columns: Column[] = [
    { name: "ID", uid: "id" },
    { name: "USERNAME", uid: "username" },
    { name: "NAMA LENGKAP", uid: "full_name" },
    { name: "EMAIL", uid: "email" },
    { name: "ROLE", uid: "role" },
  ];

  const renderCell = useCallback((item: User, columnKey: Key) => {
    const cellValue = item[columnKey as keyof User];

    switch (columnKey) {
      case "role":
        return (
          <Chip
            className="capitalize"
            color={cellValue === "admin" ? "secondary" : "default"}
            size="sm"
            variant="flat">
            {cellValue}
          </Chip>
        );
      case "full_name":
        return (
          <div className="flex flex-col">
            <p className="text-bold text-sm capitalize">{item.full_name}</p>
          </div>
        );
      default:
        return cellValue;
    }
  }, []);

  if (authLoading || user?.role !== "admin") {
    return (
      <div className="flex h-screen items-center justify-center bg-gray-50">
        <Spinner label="Memuat..." />
      </div>
    );
  }

  return (
    <AppLayout>
      <div className="h-full w-full bg-white p-6">
        <div className="flex flex-col gap-4">
          {/* Header Section */}
          <div className="flex flex-col justify-between gap-4 sm:flex-row sm:items-end">
            <div className="flex flex-col gap-1">
              <h1 className="text-2xl font-bold tracking-tight text-gray-900">
                Manajemen Pengguna
              </h1>
              <p className="text-sm text-gray-500">
                Kelola data pengguna dan hak akses aplikasi.
              </p>
            </div>
          </div>

          {/* Table Section */}
          <DataTable
            columns={columns}
            data={items}
            isLoading={isLoading}
            emptyContent="Tidak ada data user."
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
            searchPlaceholder="Cari user..."
            renderCell={renderCell}
            topContent={
              <Button
                className="bg-[#005F02] text-white font-medium"
                size="md"
                endContent={<PlusIcon className="h-4 w-4 font-bold" />}
                onPress={handleCreateOpen}>
                Tambah User
              </Button>
            }
          />
        </div>

        {/* Modal Form */}
        <Modal
          isOpen={isFormOpen}
          onOpenChange={onFormOpenChange}
          placement="top-center"
          backdrop="opaque">
          <ModalContent>
            {(onClose) => (
              <>
                <ModalHeader className="flex flex-col gap-1">
                  Tambah User Baru
                </ModalHeader>
                <ModalBody>
                  <Input
                    autoFocus
                    label="Username"
                    placeholder="Contoh: admin01"
                    variant="bordered"
                    value={formData.username}
                    onValueChange={(v) =>
                      setFormData({ ...formData, username: v })
                    }
                  />
                  <Input
                    label="Password"
                    placeholder="Minimal 6 karakter"
                    type="password"
                    variant="bordered"
                    value={formData.password}
                    onValueChange={(v) =>
                      setFormData({ ...formData, password: v })
                    }
                  />
                  <Input
                    label="Nama Lengkap"
                    placeholder="Contoh: Budi Santoso"
                    variant="bordered"
                    value={formData.full_name}
                    onValueChange={(v) =>
                      setFormData({ ...formData, full_name: v })
                    }
                  />
                  <Input
                    label="Email"
                    placeholder="budi@example.com"
                    type="email"
                    variant="bordered"
                    value={formData.email}
                    onValueChange={(v) =>
                      setFormData({ ...formData, email: v })
                    }
                  />
                  <Select
                    label="Role"
                    placeholder="Pilih Role"
                    selectedKeys={[formData.role]}
                    variant="bordered"
                    onSelectionChange={(keys: any) => {
                      const selected = Array.from(keys)[0] as "admin" | "staff";
                      if (selected) {
                        setFormData({ ...formData, role: selected });
                      }
                    }}>
                    <SelectItem key="admin" value="admin">
                      Admin
                    </SelectItem>
                    <SelectItem key="staff" value="staff">
                      Staff
                    </SelectItem>
                  </Select>
                </ModalBody>
                <ModalFooter>
                  <Button color="danger" variant="flat" onPress={onClose}>
                    Batal
                  </Button>
                  <Button
                    className="bg-[#005F02] text-white"
                    onPress={handleSubmit}>
                    Simpan
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

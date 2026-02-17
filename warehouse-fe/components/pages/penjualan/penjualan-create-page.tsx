"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import {
  Input,
  Button,
  Card,
  CardHeader,
  CardBody,
  CardFooter,
  Autocomplete,
  AutocompleteItem,
  Table,
  TableHeader,
  TableColumn,
  TableBody,
  TableRow,
  TableCell,
  Textarea,
} from "@nextui-org/react";
import {
  TrashIcon,
  PlusIcon,
  ShoppingCartIcon,
} from "@heroicons/react/24/outline";
import { AppLayout } from "@/components/layout/app-layout";
import { barangApi, penjualanApi } from "@/lib/api";
import { Barang, CreatePenjualanRequest } from "@/lib/types";
import { useAuth } from "@/components/providers/auth-provider";

export default function PenjualanCreatePage() {
  const router = useRouter();
  const { user } = useAuth();
  const [barangList, setBarangList] = useState<Barang[]>([]);
  const [cart, setCart] = useState<{ barang: Barang; qty: number }[]>([]);
  const [customer, setCustomer] = useState("");
  const [note, setNote] = useState("");
  const [loading, setLoading] = useState(false);
  const [selectedItemKey, setSelectedItemKey] = useState<
    string | number | null
  >(null);

  useEffect(() => {
    barangApi
      .getAllWithStok({ limit: 1000 })
      .then((res) => setBarangList(res.data || []));
  }, []);

  const addToCart = (key: string | number) => {
    const item = barangList.find((b) => b.id.toString() === key.toString());
    if (!item) return;

    setCart((prev) => {
      const existing = prev.find((p) => p.barang.id === item.id);
      if (existing) {
        return prev.map((p) =>
          p.barang.id === item.id ? { ...p, qty: p.qty + 1 } : p,
        );
      }
      return [...prev, { barang: item, qty: 1 }];
    });
    setSelectedItemKey(null);
  };

  const updateQty = (id: number, delta: number) => {
    setCart((prev) =>
      prev.map((item) => {
        if (item.barang.id === id) {
          const newQty = Math.max(1, item.qty + delta);
          // Check stock
          if (newQty > (item.barang.stok || 0)) {
            alert(`Stok tidak cukup! Tersedia: ${item.barang.stok}`);
            return item;
          }
          return { ...item, qty: newQty };
        }
        return item;
      }),
    );
  };

  const removeFromCart = (id: number) => {
    setCart((prev) => prev.filter((item) => item.barang.id !== id));
  };

  const calculateTotal = () => {
    return cart.reduce(
      (acc, item) => acc + item.barang.harga_jual * item.qty,
      0,
    );
  };

  const handleSubmit = async () => {
    if (cart.length === 0) {
      alert("Keranjang belanja kosong!");
      return;
    }
    if (!customer.trim()) {
      alert("Nama customer wajib diisi!");
      return;
    }

    setLoading(true);

    try {
      const payload: CreatePenjualanRequest = {
        no_faktur: `INV/OUT/${Date.now()}`,
        customer,
        user_id: user?.id || 1,
        details: cart.map((item) => ({
          barang_id: item.barang.id,
          qty: item.qty,
          harga: item.barang.harga_jual,
        })),
      };

      await penjualanApi.create(payload);
      router.push("/penjualan");
    } catch (error) {
      console.error("Failed to create transaction", error);
      alert("Gagal membuat transaksi. Silakan coba lagi.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <AppLayout>
      <div className="flex flex-col lg:flex-row gap-6 h-[calc(100vh-100px)]">
        {/* Left: Product Selection */}
        <div className="flex-1 flex flex-col gap-4">
          <Card className="flex-1">
            <CardHeader className="flex flex-col items-start gap-2">
              <h2 className="text-xl font-bold">POS - Penjualan Baru</h2>
              <p className="text-gray-500 text-sm">
                Cari barang dan tambahkan ke keranjang.
              </p>
            </CardHeader>
            <CardBody className="overflow-visible">
              <div className="flex gap-2 mb-4">
                <Autocomplete
                  label="Cari Barang"
                  placeholder="Ketik nama atau kode barang..."
                  className="max-w-xs"
                  defaultItems={barangList}
                  onSelectionChange={(key) => key && addToCart(key)}
                  selectedKey={selectedItemKey}>
                  {(item) => (
                    <AutocompleteItem
                      key={item.id}
                      textValue={item.nama_barang}>
                      <div className="flex justify-between items-center">
                        <span>{item.nama_barang}</span>
                        <span className="text-xs text-gray-500">
                          Stok: {item.stok}
                        </span>
                      </div>
                    </AutocompleteItem>
                  )}
                </Autocomplete>
              </div>

              <div className="grid grid-cols-2 md:grid-cols-3 gap-3 overflow-y-auto max-h-[500px]">
                {barangList.slice(0, 12).map((item) => (
                  <Card
                    isPressable
                    key={item.id}
                    onPress={() => addToCart(item.id)}
                    className="bg-white border border-gray-100 shadow-sm hover:border-green-500 hover:shadow-md transition-all">
                    <CardBody className="p-3 text-center">
                      <CubeIcon className="w-8 h-8 mx-auto text-green-600 mb-2" />
                      <h4 className="font-semibold text-sm truncate">
                        {item.nama_barang}
                      </h4>
                      <p className="text-xs text-gray-500">
                        Rp {item.harga_jual.toLocaleString()}
                      </p>
                      <p className="text-[10px] mt-1 text-gray-400">
                        Stok: {item.stok}
                      </p>
                    </CardBody>
                  </Card>
                ))}
              </div>
            </CardBody>
          </Card>
        </div>

        {/* Right: Cart Summary */}
        <div className="w-full lg:w-[400px] flex flex-col">
          <Card className="h-full flex flex-col">
            <CardHeader className="bg-gray-50 border-b">
              <div className="flex items-center gap-2">
                <ShoppingCartIcon className="w-5 h-5" />
                <h3 className="font-bold">Keranjang Belanja</h3>
              </div>
            </CardHeader>
            <CardBody className="flex-1 overflow-y-auto p-0">
              {cart.length === 0 ? (
                <div className="h-full flex flex-col items-center justify-center text-gray-400">
                  <ShoppingCartIcon className="w-16 h-16 mb-2 opacity-20" />
                  <p>Keranjang kosong</p>
                </div>
              ) : (
                <div className="divide-y">
                  {cart.map((item, idx) => (
                    <div
                      key={idx}
                      className="p-4 flex justify-between items-center hover:bg-gray-50 dark:hover:bg-zinc-800">
                      <div className="flex-1">
                        <h4 className="font-medium text-sm">
                          {item.barang.nama_barang}
                        </h4>
                        <p className="text-xs text-green-600 font-semibold">
                          Rp {item.barang.harga_jual.toLocaleString()}
                        </p>
                      </div>
                      <div className="flex items-center gap-3">
                        <div className="flex items-center border rounded-lg">
                          <button
                            className="px-2 py-1 text-gray-500 hover:bg-gray-200"
                            onClick={() => updateQty(item.barang.id, -1)}>
                            -
                          </button>
                          <span className="px-2 text-sm font-medium">
                            {item.qty}
                          </span>
                          <button
                            className="px-2 py-1 text-gray-500 hover:bg-gray-200"
                            onClick={() => updateQty(item.barang.id, 1)}>
                            +
                          </button>
                        </div>
                        <button
                          className="text-red-500 hover:text-red-700"
                          onClick={() => removeFromCart(item.barang.id)}>
                          <TrashIcon className="w-4 h-4" />
                        </button>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </CardBody>
            <CardFooter className="flex-col gap-4 bg-white border-t p-4">
              <div className="w-full space-y-3">
                <Input
                  label="Nama Customer"
                  size="sm"
                  variant="bordered"
                  value={customer}
                  onValueChange={setCustomer}
                  isRequired
                />
                <Textarea
                  label="Catatan"
                  size="sm"
                  variant="bordered"
                  minRows={1}
                  value={note}
                  onValueChange={setNote}
                />
              </div>
              <div className="w-full flex justify-between items-center pt-2 border-t">
                <span className="text-lg font-bold">Total</span>
                <span className="text-xl font-bold text-green-600">
                  Rp {calculateTotal().toLocaleString()}
                </span>
              </div>
              <Button
                className="w-full text-white font-bold"
                color="success"
                size="md"
                onPress={handleSubmit}
                isLoading={loading}
                isDisabled={cart.length === 0}>
                Proses Pembayaran
              </Button>
            </CardFooter>
          </Card>
        </div>
      </div>
    </AppLayout>
  );
}

function CubeIcon(props: any) {
  return (
    <svg
      {...props}
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 24 24"
      strokeWidth={1.5}
      stroke="currentColor">
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M21 7.5l-9-5.25L3 7.5m18 0l-9 5.25m9-5.25v9l-9 5.25M3 7.5l9 5.25M3 7.5v9l-9 5.25m9-5.25l9-5.25"
      />
    </svg>
  );
}

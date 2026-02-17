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
  Textarea,
} from "@nextui-org/react";
import { TrashIcon, ShoppingCartIcon } from "@heroicons/react/24/outline";
import { AppLayout } from "@/components/layout/app-layout";
import { barangApi, pembelianApi } from "@/lib/api";
import { Barang, CreatePembelianRequest } from "@/lib/types";
import { useAuth } from "@/components/providers/auth-provider";

export default function PembelianCreatePage() {
  const router = useRouter();
  const { user } = useAuth();
  const [barangList, setBarangList] = useState<Barang[]>([]);
  const [cart, setCart] = useState<
    { barang: Barang; qty: number; harga: number }[]
  >([]);
  const [supplier, setSupplier] = useState("");
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
      // Default to stored buy price, but allow editing later if needed
      return [...prev, { barang: item, qty: 1, harga: item.harga_beli }];
    });
    setSelectedItemKey(null);
  };

  const updateQty = (id: number, delta: number) => {
    setCart((prev) =>
      prev.map((item) => {
        if (item.barang.id === id) {
          const newQty = Math.max(1, item.qty + delta);
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
    return cart.reduce((acc, item) => acc + item.harga * item.qty, 0);
  };

  const handleSubmit = async () => {
    if (cart.length === 0) return alert("Keranjang kosong!");
    if (!supplier) return alert("Nama supplier wajib diisi!");

    setLoading(true);
    try {
      const payload: CreatePembelianRequest = {
        no_faktur: `INV/IN/${Date.now()}`,
        supplier,
        user_id: user?.id || 1,
        details: cart.map((item) => ({
          barang_id: item.barang.id,
          qty: item.qty,
          harga: item.harga,
        })),
      };

      await pembelianApi.create(payload);
      router.push("/pembelian");
    } catch (error) {
      console.error(error);
      alert("Gagal membuat transaksi pembelian");
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
              <h2 className="text-xl font-bold">Pembelian Baru (Restock)</h2>
              <p className="text-gray-500 text-sm">
                Pilih barang untuk restock stok gudang.
              </p>
            </CardHeader>
            <CardBody className="overflow-visible">
              <div className="flex gap-2 mb-4">
                <Autocomplete
                  label="Cari Barang"
                  placeholder="Ketik nama atau kode barang..."
                  className="max-w-md"
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

              <div className="mt-4">
                <h3 className="font-semibold mb-2">Item di Keranjang</h3>
                {cart.length === 0 ? (
                  <p className="text-gray-400 text-sm">
                    Belum ada item dipilih.
                  </p>
                ) : (
                  <div className="space-y-2">
                    {cart.map((item) => (
                      <Card
                        key={item.barang.id}
                        className="border border-gray-100 shadow-sm">
                        <CardBody className="flex flex-row items-center justify-between py-2">
                          <div className="flex-1">
                            <div className="font-medium">
                              {item.barang.nama_barang}
                            </div>
                            <div className="text-xs text-gray-500">
                              {item.barang.kode_barang}
                            </div>
                          </div>
                          <div className="flex items-center gap-4">
                            <div className="w-24">
                              <Input
                                type="number"
                                label="Harga Beli"
                                size="sm"
                                value={item.harga.toString()}
                                onValueChange={(val) => {
                                  const newPrice = Number(val);
                                  setCart((prev) =>
                                    prev.map((p) =>
                                      p.barang.id === item.barang.id
                                        ? { ...p, harga: newPrice }
                                        : p,
                                    ),
                                  );
                                }}
                                startContent={
                                  <span className="text-xs text-gray-500">
                                    Rp
                                  </span>
                                }
                              />
                            </div>
                            <div className="flex items-center border rounded-lg h-10">
                              <button
                                className="px-3 hover:bg-gray-100"
                                onClick={() => updateQty(item.barang.id, -1)}>
                                -
                              </button>
                              <span className="px-3 font-medium">
                                {item.qty}
                              </span>
                              <button
                                className="px-3 hover:bg-gray-100"
                                onClick={() => updateQty(item.barang.id, 1)}>
                                +
                              </button>
                            </div>
                            <div className="font-semibold w-24 text-right">
                              Rp {(item.harga * item.qty).toLocaleString()}
                            </div>
                            <button
                              onClick={() => removeFromCart(item.barang.id)}
                              className="text-red-500 hover:bg-red-50 p-2 rounded-full">
                              <TrashIcon className="w-5 h-5" />
                            </button>
                          </div>
                        </CardBody>
                      </Card>
                    ))}
                  </div>
                )}
              </div>
            </CardBody>
          </Card>
        </div>

        {/* Right: Summary */}
        <div className="w-full lg:w-[350px]">
          <Card className="sticky top-4">
            <CardHeader className="border-b bg-gray-50">
              <div className="flex items-center gap-2">
                <ShoppingCartIcon className="w-5 h-5" />
                <h3 className="font-bold">Ringkasan Pesanan</h3>
              </div>
            </CardHeader>
            <CardBody className="gap-4">
              <Input
                label="Nama Supplier"
                placeholder="PT. Example Supplier"
                labelPlacement="outside"
                value={supplier}
                onValueChange={setSupplier}
                isRequired
              />
              <Textarea
                label="Catatan"
                placeholder="Optional memo..."
                labelPlacement="outside"
                minRows={2}
                value={note}
                onValueChange={setNote}
              />

              <div className="border-t pt-4 mt-2 space-y-2">
                <div className="flex justify-between">
                  <span className="text-gray-500">Total Item</span>
                  <span className="font-semibold">
                    {cart.reduce((a, b) => a + b.qty, 0)} pcs
                  </span>
                </div>
                <div className="flex justify-between text-lg font-bold">
                  <span>Grand Total</span>
                  <span className="text-primary">
                    Rp {calculateTotal().toLocaleString()}
                  </span>
                </div>
              </div>
            </CardBody>
            <CardFooter>
              <Button
                className="w-full font-bold text-white"
                color="primary"
                size="md"
                onPress={handleSubmit}
                isLoading={loading}
                isDisabled={cart.length === 0}>
                Simpan Pembelian
              </Button>
            </CardFooter>
          </Card>
        </div>
      </div>
    </AppLayout>
  );
}

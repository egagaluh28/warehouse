"use client";

import type { BeliHeader, JualHeader } from "./types";

function formatRupiah(value: number) {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    minimumFractionDigits: 0,
  }).format(value || 0);
}

function formatDateTime(value?: string) {
  if (!value) return "-";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return String(value);
  return date.toLocaleString("id-ID");
}

function safeFilePart(value: string) {
  return value
    .trim()
    .replace(/[\\/:*?"<>|]+/g, "-")
    .replace(/\s+/g, "-")
    .slice(0, 80);
}

async function createDoc() {
  const jsPDFModule = await import("jspdf");
  const autoTableModule: any = await import("jspdf-autotable");

  const jsPDF = (jsPDFModule as any).jsPDF ?? (jsPDFModule as any).default;
  const autoTable = autoTableModule.default ?? autoTableModule;

  const doc = new jsPDF({ unit: "pt", format: "a4" });
  return { doc, autoTable } as const;
}

export async function downloadPembelianInvoice(tx: BeliHeader) {
  const { doc, autoTable } = await createDoc();

  const marginX = 40;
  const titleY = 48;

  doc.setFontSize(16);
  doc.text("INVOICE PEMBELIAN", marginX, titleY);

  doc.setFontSize(10);
  doc.text(`No Faktur: ${tx.no_faktur || "-"}`, marginX, titleY + 22);
  doc.text(`Tanggal: ${formatDateTime(tx.created_at)}`, marginX, titleY + 38);
  doc.text(`Supplier: ${tx.supplier || "-"}`, marginX, titleY + 54);

  const rows = (tx.details || []).map((d, idx) => [
    String(idx + 1),
    d.barang?.kode_barang || String(d.barang_id ?? "-"),
    d.barang?.nama_barang || "-",
    String(d.qty ?? 0),
    formatRupiah(Number(d.harga || 0)),
    formatRupiah(Number(d.subtotal || 0)),
  ]);

  autoTable(doc, {
    startY: titleY + 80,
    head: [["No", "Kode", "Barang", "Qty", "Harga", "Subtotal"]],
    body: rows,
    styles: { fontSize: 9, cellPadding: 4 },
    headStyles: { fillColor: [245, 245, 245], textColor: 60 },
    columnStyles: {
      0: { halign: "right", cellWidth: 30 },
      1: { cellWidth: 70 },
      2: { cellWidth: 200 },
      3: { halign: "right", cellWidth: 40 },
      4: { halign: "right", cellWidth: 80 },
      5: { halign: "right", cellWidth: 90 },
    },
  });

  const finalY = (doc as any).lastAutoTable?.finalY;
  const totalY = (typeof finalY === "number" ? finalY : titleY + 120) + 20;

  doc.setFontSize(11);
  doc.text("Total:", marginX + 280, totalY);
  doc.text(formatRupiah(Number(tx.total || 0)), marginX + 420, totalY, {
    align: "right",
  });

  doc.save(`invoice-pembelian-${safeFilePart(tx.no_faktur || "")}.pdf`);
}

export async function downloadPenjualanInvoice(tx: JualHeader) {
  const { doc, autoTable } = await createDoc();

  const marginX = 40;
  const titleY = 48;

  doc.setFontSize(16);
  doc.text("INVOICE PENJUALAN", marginX, titleY);

  doc.setFontSize(10);
  doc.text(`No Faktur: ${tx.no_faktur || "-"}`, marginX, titleY + 22);
  doc.text(`Tanggal: ${formatDateTime(tx.created_at)}`, marginX, titleY + 38);
  doc.text(`Customer: ${tx.customer || "-"}`, marginX, titleY + 54);

  const rows = (tx.details || []).map((d, idx) => [
    String(idx + 1),
    d.barang?.kode_barang || String(d.barang_id ?? "-"),
    d.barang?.nama_barang || "-",
    String(d.qty ?? 0),
    formatRupiah(Number(d.harga || 0)),
    formatRupiah(Number(d.subtotal || 0)),
  ]);

  autoTable(doc, {
    startY: titleY + 80,
    head: [["No", "Kode", "Barang", "Qty", "Harga", "Subtotal"]],
    body: rows,
    styles: { fontSize: 9, cellPadding: 4 },
    headStyles: { fillColor: [245, 245, 245], textColor: 60 },
    columnStyles: {
      0: { halign: "right", cellWidth: 30 },
      1: { cellWidth: 70 },
      2: { cellWidth: 200 },
      3: { halign: "right", cellWidth: 40 },
      4: { halign: "right", cellWidth: 80 },
      5: { halign: "right", cellWidth: 90 },
    },
  });

  const finalY = (doc as any).lastAutoTable?.finalY;
  const totalY = (typeof finalY === "number" ? finalY : titleY + 120) + 20;

  doc.setFontSize(11);
  doc.text("Total:", marginX + 280, totalY);
  doc.text(formatRupiah(Number(tx.total || 0)), marginX + 420, totalY, {
    align: "right",
  });

  doc.save(`invoice-penjualan-${safeFilePart(tx.no_faktur || "")}.pdf`);
}

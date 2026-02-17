"use client";

import {
  Input,
  Pagination,
  Select,
  SelectItem,
  SortDescriptor,
  Spinner,
  Table,
  TableBody,
  TableCell,
  TableColumn,
  TableHeader,
  TableRow,
} from "@nextui-org/react";
import { MagnifyingGlassIcon } from "@heroicons/react/24/outline";
import { Key, ReactNode, useCallback, useMemo } from "react";

export interface Column {
  uid: string;
  name: string;
  sortable?: boolean;
  align?: "start" | "center" | "end";
}

export interface DataTableProps<T> {
  columns: Column[];
  data: T[];
  isLoading?: boolean;
  emptyContent?: string;

  getRowKey?: (item: T) => Key;

  // Pagination
  page?: number;
  totalPages?: number;
  onPageChange?: (page: number) => void;
  rowsPerPage?: number;
  onRowsPerPageChange?: (value: number) => void;

  // Search
  search?: string;
  onSearchChange?: (value: string) => void;
  searchPlaceholder?: string;

  // Sorting
  sortDescriptor?: SortDescriptor;
  onSortChange?: (descriptor: SortDescriptor) => void;

  // Visibility
  showTopBar?: boolean;
  showRowsPerPage?: boolean;
  showSearch?: boolean;

  // Rendering
  renderCell: (item: T, columnKey: Key) => ReactNode;

  // Custom Content
  topContent?: ReactNode;
  filtersContent?: ReactNode;
  bottomContent?: ReactNode;
}

export function DataTable<T>({
  columns,
  data,
  isLoading = false,
  emptyContent = "Data tidak ditemukan.",
  getRowKey,

  page = 1,
  totalPages = 1,
  onPageChange,
  rowsPerPage = 10,
  onRowsPerPageChange,

  search,
  onSearchChange,
  searchPlaceholder = "Cari...",

  sortDescriptor,
  onSortChange,

  showTopBar = true,
  showRowsPerPage = true,
  showSearch = true,

  renderCell,
  topContent,
  filtersContent,
  bottomContent,
}: DataTableProps<T>) {
  const resolveRowKey = useCallback(
    (item: T) => {
      if (getRowKey) return getRowKey(item);
      return (item as any)?.id;
    },
    [getRowKey],
  );

  const topBar = useMemo(() => {
    if (!showTopBar) return null;

    return (
      <div className="flex flex-col gap-3 mb-4">
        <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-3">
          {showSearch || filtersContent || showRowsPerPage ? (
            <div className="flex flex-col gap-2 w-full sm:w-auto">
              {showSearch || filtersContent ? (
                <div className="flex flex-col sm:flex-row sm:items-end gap-2 w-full">
                  {showSearch ? (
                    <Input
                      isClearable
                      className="w-full sm:w-64"
                      placeholder={searchPlaceholder}
                      startContent={
                        <MagnifyingGlassIcon className="h-4 w-4 text-default-400" />
                      }
                      value={search}
                      onValueChange={onSearchChange}
                      size="md"
                      variant="bordered"
                      color="primary"
                    />
                  ) : null}

                  {filtersContent ? (
                    <div className="flex flex-col sm:flex-row gap-2 w-full sm:w-auto">
                      {filtersContent}
                    </div>
                  ) : null}
                </div>
              ) : null}

              {showRowsPerPage ? (
                <div className="flex items-center gap-2 text-default-400 text-small">
                  <span>Show</span>
                  <Select
                    className="min-w-[10px] w-20"
                    selectedKeys={[String(rowsPerPage)]}
                    size="md"
                    variant="bordered"
                    onChange={(e) =>
                      onRowsPerPageChange?.(Number(e.target.value))
                    }
                    aria-label="Rows per page">
                    <SelectItem key="5" value="5">
                      5
                    </SelectItem>
                    <SelectItem key="10" value="10">
                      10
                    </SelectItem>
                    <SelectItem key="20" value="20">
                      20
                    </SelectItem>
                    <SelectItem key="50" value="50">
                      50
                    </SelectItem>
                  </Select>
                  <span>entries</span>
                </div>
              ) : null}
            </div>
          ) : (
            <span />
          )}

          {topContent ? (
            <div className="flex justify-end w-full sm:w-auto flex-shrink-0">
              {topContent}
            </div>
          ) : null}
        </div>
      </div>
    );
  }, [
    showTopBar,
    showSearch,
    search,
    onSearchChange,
    searchPlaceholder,
    filtersContent,
    showRowsPerPage,
    rowsPerPage,
    onRowsPerPageChange,
    topContent,
  ]);

  const bottomBar = useMemo(() => {
    return (
      <div className="py-2 px-2 flex items-center gap-3">
        <div className="flex-1">{bottomContent}</div>

        {totalPages > 1 ? (
          <Pagination
            isCompact
            showControls
            showShadow
            color="primary"
            page={page}
            total={totalPages}
            onChange={onPageChange}
            className="ml-auto"
          />
        ) : null}
      </div>
    );
  }, [bottomContent, onPageChange, page, totalPages]);

  return (
    <div className="p-4 bg-white rounded-xl shadow-sm border border-slate-100">
      <Table
        aria-label="Data Table"
        bottomContent={bottomBar}
        topContent={topBar}
        topContentPlacement="outside"
        bottomContentPlacement="outside"
        sortDescriptor={sortDescriptor}
        onSortChange={onSortChange}
        classNames={{
          wrapper: "shadow-none p-0",
          th: "bg-slate-50 text-slate-600 font-semibold border-b border-slate-100",
          td: "py-3 border-b border-slate-50 last:border-b-0",
        }}>
        <TableHeader columns={columns}>
          {(column) => (
            <TableColumn
              key={column.uid}
              align={column.align || "start"}
              allowsSorting={column.sortable}>
              {column.name}
            </TableColumn>
          )}
        </TableHeader>
        <TableBody
          items={data}
          emptyContent={emptyContent}
          isLoading={isLoading}
          loadingContent={<Spinner label="Memuat data..." />}>
          {(item) => {
            const key = resolveRowKey(item);
            const safeKey = key ?? JSON.stringify(item);

            return (
              <TableRow
                key={safeKey}
                className="hover:bg-slate-50/50 transition-colors">
                {(columnKey) => (
                  <TableCell>{renderCell(item, columnKey)}</TableCell>
                )}
              </TableRow>
            );
          }}
        </TableBody>
      </Table>
    </div>
  );
}

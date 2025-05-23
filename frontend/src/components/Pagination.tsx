import { type FC } from "react";

interface PaginationProps {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
  className?: string;
}

export const Pagination: FC<PaginationProps> = ({
  currentPage,
  totalPages,
  onPageChange,
  className = "",
}) => {
  const getPageNumbers = () => {
    const pageNumbers = [];

    const maxVisibleButtons = 5;

    if (totalPages <= maxVisibleButtons) {
      for (let i = 1; i <= totalPages; i++) {
        pageNumbers.push(i);
      }
    } else {
      let startPage = Math.max(1, currentPage - Math.floor(maxVisibleButtons / 2));
      let endPage = startPage + maxVisibleButtons - 1;

      if (endPage > totalPages) {
        endPage = totalPages;
        startPage = Math.max(1, endPage - maxVisibleButtons + 1);
      }

      if (startPage > 1) {
        pageNumbers.push(1);
        if (startPage > 2) {
          pageNumbers.push("...");
        }
      }

      for (let i = startPage; i <= endPage; i++) {
        pageNumbers.push(i);
      }

      if (endPage < totalPages) {
        if (endPage < totalPages - 1) {
          pageNumbers.push("...");
        }
        pageNumbers.push(totalPages);
      }
    }

    return pageNumbers;
  };

  const goToPreviousPage = () => {
    if (currentPage > 1) {
      onPageChange(currentPage - 1);
    }
  };

  const goToNextPage = () => {
    if (currentPage < totalPages) {
      onPageChange(currentPage + 1);
    }
  };

  const getButtonStyles = (page: number | string) => {
    if (page === "...") {
      return "px-3 py-2 text-gray-400 cursor-pointer";
    }

    return page === currentPage
      ? "px-3 py-2 bg-blue-600 text-white rounded cursor-pointer"
      : "px-3 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded cursor-pointer";
  };

  return (
    <div className={`items-center justify-center space-x-1 my-4 ${className}`}>
      <button
        onClick={goToPreviousPage}
        disabled={currentPage === 1}
        className={`px-3 py-2 rounded ${
          currentPage === 1
            ? "bg-gray-800 text-gray-500 cursor-not-allowed"
            : "bg-gray-700 hover:bg-gray-600 text-white cursor-pointer "
        }`}
        aria-label="Previous page"
      >
        &lt;
      </button>

      {getPageNumbers().map((page, index) => (
        <button
          key={index}
          onClick={() => (typeof page === "number" ? onPageChange(page) : null)}
          className={getButtonStyles(page)}
          disabled={page === "..."}
        >
          {page}
        </button>
      ))}

      <button
        onClick={goToNextPage}
        disabled={currentPage === totalPages}
        className={`px-3 py-2 rounded ${
          currentPage === totalPages
            ? "bg-gray-800 text-gray-500 cursor-not-allowed"
            : "bg-gray-700 hover:bg-gray-600 text-white cursor-pointer "
        }`}
        aria-label="Next page"
      >
        &gt;
      </button>
    </div>
  );
};

import { SERVER } from "../../constants";
import { useState, useEffect, type MouseEvent, useRef } from "react";
import { Settings2, Edit2, Trash2 } from "lucide-react";
import { axiosInstance } from "../../utils";
import { type Ifile } from "./types";

const FileSettings = ({
  fileId,
  originalName,
  close,
  refetch,
}: {
  fileId: number;
  originalName: string;
  close: () => void;
  refetch: () => void;
}) => {
  const [isRenaming, setIsRenaming] = useState(false);
  const [newName, setNewName] = useState("");
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (isRenaming && inputRef.current) {
      const fileName = originalName;
      const lastDotIndex = fileName.lastIndexOf(".");
      if (lastDotIndex > 0) {
        inputRef.current.setSelectionRange(0, lastDotIndex);
      } else {
        inputRef.current.select();
      }
    }
  }, [isRenaming, originalName]);

  const handleRename = (ev: MouseEvent<HTMLButtonElement>) => {
    ev.stopPropagation();
    setIsRenaming(true);
    setNewName(originalName);
  };

  const handleRenameSubmit = async () => {
    if (newName.trim().length === 0) return;

    try {
      await axiosInstance.put(`${SERVER}/files/${fileId}`, {
        newName: newName.trim(),
      });
      refetch();
      close();
    } catch (err) {
      console.error("rename failed:", err);
    }
  };

  const handleRenameCancel = () => {
    setIsRenaming(false);
    setNewName("");
  };

  const handleDelete = async (ev: MouseEvent<HTMLButtonElement>) => {
    ev.stopPropagation();

    try {
      await axiosInstance.delete(`${SERVER}/files/${fileId}`);
      refetch();
      close();
    } catch (err) {
      console.error("delete failed:", err);
    }
  };

  if (isRenaming) {
    return (
      <div className="absolute top-20 left-10 bg-white dark:bg-gray-800 rounded-lg shadow-lg border border-gray-200 dark:border-gray-600 p-4 z-20">
        <div className="mb-3">
          <label className="block text-sm font-medium mb-2">New filename:</label>
          <input
            type="text"
            ref={inputRef}
            value={newName}
            onChange={(e) => setNewName(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter") handleRenameSubmit();
              if (e.key === "Escape") handleRenameCancel();
            }}
            className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            onFocus={(e) => {
              const fileName = e.target.value;
              const lastDotIndex = fileName.lastIndexOf(".");
              if (lastDotIndex > 0) {
                e.target.setSelectionRange(0, lastDotIndex);
              } else {
                e.target.select();
              }
            }}
            autoFocus
          />
        </div>
        <div className="flex gap-2">
          <button
            onClick={handleRenameSubmit}
            className="flex-1 px-3 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600"
          >
            Save
          </button>
          <button
            onClick={handleRenameCancel}
            className="flex-1 px-3 py-2 bg-gray-500 text-white rounded-lg hover:bg-gray-600"
          >
            Cancel
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="absolute top-20 left-10 bg-white dark:bg-gray-800 rounded-lg shadow-lg border border-gray-200 dark:border-gray-600 p-2 z-20">
      <button
        className="w-full px-4 py-2 text-left hover:bg-gray-100 dark:hover:bg-gray-700 flex items-center gap-2 mb-1"
        onClick={(ev) => handleRename(ev)}
      >
        <Edit2 className="w-4 h-4" />
        Rename
      </button>
      <button
        className="w-full px-4 py-2 text-left hover:bg-red-50 dark:hover:bg-red-900/20 text-red-600 dark:text-red-400 flex items-center gap-2 mb-1"
        onClick={(ev) => handleDelete(ev)}
      >
        <Trash2 className="w-4 h-4" />
        Delete
      </button>
    </div>
  );
};

export const LightBox = ({
  idx,
  files,
  onClose,
  refetch,
}: {
  idx: number;
  files: Ifile[];
  onClose: () => void;
  refetch: () => void;
}) => {
  const [currentIdx, setCurrentIdx] = useState(idx);
  const [openSettings, setOpenSettings] = useState(false);

  useEffect(() => {
    setCurrentIdx(idx);
  }, [idx]);

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        onClose();
      } else if (e.key === "ArrowLeft") {
        setCurrentIdx((prev) => {
          const newIdx = Math.max(0, prev - 1);
          return newIdx;
        });
      } else if (e.key === "ArrowRight") {
        setCurrentIdx((prev) => {
          const newIdx = Math.min(files.length - 1, prev + 1);
          return newIdx;
        });
      }
    };
    document.addEventListener("keydown", handleKeyDown);
    return () => document.removeEventListener("keydown", handleKeyDown);
  }, [onClose, files.length]);

  const currentFile = files[currentIdx];
  const canGoPrev = currentIdx > 0;
  const canGoNext = currentIdx < files.length - 1;

  if (!currentFile) {
    return null;
  }

  return (
    // <div
    //   className="fixed inset-0 bg-black bg-opacity-90 flex items-center justify-center z-50"
    //   onClick={onClose}
    // >
    <div className="fixed inset-0 bg-black bg-opacity-90 flex items-center justify-center z-50">
      <Settings2
        className="absolute top-10 left-10 cursor-pointer"
        onClick={(ev) => {
          ev.stopPropagation();
          setOpenSettings(!openSettings);
        }}
      />

      {openSettings && (
        <FileSettings
          fileId={currentFile.id}
          originalName={currentFile.originalName}
          close={onClose}
          refetch={refetch}
        />
      )}

      {canGoPrev && (
        <button
          onClick={(e) => {
            e.stopPropagation();
            setCurrentIdx((prev) => {
              const newIdx = Math.max(0, prev - 1);
              return newIdx;
            });
          }}
          className="absolute left-4 top-1/2 transform -translate-y-1/2 text-white text-4xl hover:text-gray-300 z-10 w-12 h-12 flex items-center justify-center bg-black bg-opacity-50 rounded-full"
        >
          ‹
        </button>
      )}

      {canGoNext && (
        <button
          onClick={(e) => {
            e.stopPropagation();
            setCurrentIdx((prev) => {
              const newIdx = Math.min(files.length - 1, prev + 1);
              return newIdx;
            });
          }}
          className="absolute right-4 top-1/2 transform -translate-y-1/2 text-white text-4xl hover:text-gray-300 z-10 w-12 h-12 flex items-center justify-center bg-black bg-opacity-50 rounded-full"
        >
          ›
        </button>
      )}

      <div className="relative max-w-[95vw] max-h-[95vh] flex flex-col">
        <button
          onClick={onClose}
          className="absolute top-4 right-4 text-white text-3xl hover:text-gray-300 z-10 w-8 h-8 flex items-center justify-center bg-black bg-opacity-50 rounded-full"
        >
          ×
        </button>

        {currentFile.mimeType.includes("image") && (
          <img
            src={`${SERVER}${currentFile.filePath}`}
            alt={currentFile.originalName}
            className="max-w-full max-h-[90vh] object-contain"
            onClick={(e) => e.stopPropagation()}
          />
        )}

        {currentFile.mimeType.includes("video") && (
          <video
            key={currentFile.filePath}
            className="max-w-full max-h-[90vh] object-contain"
            controls
            onClick={(e) => e.stopPropagation()}
          >
            <source src={`${SERVER}${currentFile.filePath}`} type={currentFile.mimeType} />
          </video>
        )}

        <div className="mt-4 text-white text-center">
          <p className="bg-black bg-opacity-50 px-4 py-2 rounded max-w-full truncate">
            {currentFile.originalName}
          </p>
          <p className="text-sm mt-2 opacity-75">
            {currentIdx + 1} / {files.length}
          </p>
        </div>
      </div>
    </div>
  );
};

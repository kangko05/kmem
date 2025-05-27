import { useState, type ChangeEvent } from "react";
import { PageLayout, FileInput } from "../components";
import { axiosInstance } from "../utils";
import { X } from "lucide-react";
import { useAuth } from "../hooks/useAuth";

interface FileCardProps {
  file: File;
  onRemove: () => void;
}

export const FileCard = ({ file, onRemove }: FileCardProps) => {
  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return "0 Bytes";
    const k = 1024;
    const sizes = ["Bytes", "KB", "MB", "GB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
  };

  const getFileIcon = (type: string) => {
    if (type.startsWith("image/")) return "🖼️";
    if (type.startsWith("video/")) return "🎥";
    return "📄";
  };

  return (
    <div className="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-700 rounded-lg border border-gray-200 dark:border-gray-600">
      <div className="flex items-center gap-3 flex-1 min-w-0">
        <span className="text-xl">{getFileIcon(file.type)}</span>
        <div className="flex-1 min-w-0">
          <p className="text-sm font-medium text-gray-900 dark:text-white truncate">{file.name}</p>
          <p className="text-xs text-gray-500 dark:text-gray-400">{formatFileSize(file.size)}</p>
        </div>
      </div>
      <button
        onClick={onRemove}
        className="p-1 text-gray-400 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 rounded transition-colors"
      >
        <X className="w-4 h-4" />
      </button>
    </div>
  );
};

export const UploadPage = () => {
  useAuth();

  const [files, setFiles] = useState<File[]>([]);
  const [uploading, setUploading] = useState(false);

  const handleFileChange = async (ev: ChangeEvent<HTMLInputElement>) => {
    const selectedfiles = ev.target.files;
    if (!selectedfiles) return;
    setFiles([...files, ...Array.from(selectedfiles)]);
  };

  const handleClick = async () => {
    if (files.length == 0) return;
    setUploading(true);
    for (let i = 0; i < files.length; i += 3) {
      const start = i;
      const end = Math.min(i + 3, files.length);
      const batch = files.slice(start, end);
      const promises = batch.map((file) => axiosInstance.post("/files/upload", file));
      await Promise.all(promises);
    }
    setUploading(false);
  };

  const handleRemove = (filename: string) => {
    const filtered = files.filter((file) => file.name != filename);
    setFiles([...filtered]);
  };

  return (
    <PageLayout>
      <div className="w-full max-w-2xl mx-auto p-6">
        {files.length > 0 && (
          <div className="mb-6">
            <h2 className="text-lg font-semibold mb-4 text-gray-700 dark:text-gray-300">
              selected {files.length} files
            </h2>
            <div className="space-y-3 max-h-64 overflow-y-auto">
              {files.map((file, idx) => (
                <FileCard key={idx} file={file} onRemove={() => handleRemove(file.name)} />
              ))}
            </div>
          </div>
        )}

        <div className="space-y-4">
          <FileInput onChange={handleFileChange} />

          {files.length > 0 && (
            <button
              className="w-full bg-blue-500 hover:bg-blue-600 disabled:bg-gray-400 
                        text-white font-semibold py-3 px-4 rounded-lg 
                        transition-colors duration-200 disabled:cursor-not-allowed"
              onClick={handleClick}
              disabled={uploading}
            >
              {uploading ? "Uploading..." : `Upload ${files.length} files`}
            </button>
          )}
        </div>
      </div>
    </PageLayout>
  );
};

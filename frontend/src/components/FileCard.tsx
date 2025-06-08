import { X, CheckCircle, XCircle, Loader2 } from "lucide-react";

export interface IFile {
  id: string;
  fileobj: File;
  status: "pending" | "uploading" | "success" | "error";
  msg?: string;
}

interface FileCardProps {
  file: IFile;
  progress?: number;
  onRemove: () => void;
}

export const FileCard = ({ file, progress, onRemove }: FileCardProps) => {
  const { fileobj, status, msg } = file;

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

  const getStatusIcon = () => {
    switch (status) {
      case "uploading":
        return (
          <div className="flex items-center gap-2">
            <Loader2 className="w-4 h-4 text-blue-500 animate-spin" />
            {progress && progress == 100 ? (
              <span className="text-xs animate-pulse">postprocessing...</span>
            ) : (
              <span className="text-xs">{progress}%</span>
            )}
          </div>
        );
      case "success":
        return <CheckCircle className="w-4 h-4 text-green-500" />;
      case "error":
        return <XCircle className="w-4 h-4 text-red-500" />;
      default:
        return null;
    }
  };

  const getStatusColor = () => {
    switch (status) {
      case "uploading":
        return "border-blue-200 bg-blue-50 dark:bg-blue-900/20 dark:border-blue-800";
      case "success":
        return "border-green-200 bg-green-50 dark:bg-green-900/20 dark:border-green-800";
      case "error":
        return "border-red-200 bg-red-50 dark:bg-red-900/20 dark:border-red-800";
      default:
        return "border-gray-200 bg-gray-50 dark:bg-gray-700 dark:border-gray-600";
    }
  };

  return (
    <div className={`flex items-center justify-between p-3 rounded-lg border ${getStatusColor()}`}>
      <div className="flex items-center gap-3 flex-1 min-w-0">
        <span className="text-xl">{getFileIcon(fileobj.type)}</span>
        <div className="flex-1 min-w-0">
          <p className="text-sm max-w-[90%] font-medium text-gray-900 dark:text-white truncate">
            {fileobj.name}
          </p>

          <div className="flex items-center gap-2">
            <p className="text-xs text-gray-500 dark:text-gray-400">
              {formatFileSize(fileobj.size)}
            </p>
            {msg && (
              <span className="text-xs px-2 py-1 bg-orange-100 dark:bg-orange-900/30 text-orange-700 dark:text-orange-300 rounded-full">
                {msg}
              </span>
            )}
          </div>
        </div>
      </div>

      <div className="flex items-center gap-2">
        {getStatusIcon()}
        <button
          onClick={onRemove}
          disabled={status == "uploading"}
          className="p-1 text-gray-400 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 rounded transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <X className="w-4 h-4" />
        </button>
      </div>
    </div>
  );
};

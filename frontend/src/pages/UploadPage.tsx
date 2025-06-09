import { useState, type ChangeEvent } from "react";
import { PageLayout, FileInput, FileCard, type IFile } from "../components";
import { axiosInstance } from "../utils";
import { useAuth } from "../hooks/useAuth";
import { AxiosError } from "axios";

export const UploadPage = () => {
  useAuth();

  const [files, setFiles] = useState<IFile[]>([]);
  const [uploading, setUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState<Map<string, number>>(new Map());

  const handleFileChange = (ev: ChangeEvent<HTMLInputElement>) => {
    const selectedFiles = ev.target.files;
    if (!selectedFiles) return;

    const newFiles = Array.from(selectedFiles).map(
      (file): IFile => ({
        id: `${Date.now()}-${Math.random()}`,
        fileobj: file,
        status: "pending",
      })
    );

    setFiles((prev) => [...prev, ...newFiles]);
  };

  const updateFileStatus = (fileId: string, status: IFile["status"], msg?: string) => {
    setFiles((prev) => prev.map((file) => (file.id == fileId ? { ...file, status, msg } : file)));
  };

  const updateFileProgress = (fileId: string, progress: number) => {
    setUploadProgress((prev) => new Map(prev.set(fileId, progress)));
  };

  const cleanupProgress = (fileId: string) => {
    setUploadProgress((prev) => {
      const newMap = new Map(prev);
      newMap.delete(fileId);
      return newMap;
    });
  };

  const uploadFile = async (file: IFile) => {
    try {
      updateFileStatus(file.id, "uploading");

      const encodedFilename = btoa(encodeURIComponent(file.fileobj.name));

      const resp = await axiosInstance.post(
        `/files/upload?filename=${encodedFilename}`,
        file.fileobj,
        {
          onUploadProgress: (ev) => {
            const progress = Math.round((ev.loaded * 100) / (ev.total as number));
            updateFileProgress(file.id, progress);
          },
        }
      );

      if (resp.status === 200) {
        updateFileStatus(file.id, "success");
      } else {
        updateFileStatus(file.id, "error");
      }
    } catch (err) {
      console.error("Upload failed:", file.fileobj.name, err);

      if (
        err instanceof AxiosError &&
        err?.response?.data?.message?.includes("file already exists")
      ) {
        updateFileStatus(file.id, "error", "file exists");
      } else {
        updateFileStatus(file.id, "error");
      }
    } finally {
      cleanupProgress(file.id);
    }
  };

  const handleClick = async () => {
    if (files.length == 0) return;

    setUploading(true);

    const batchSize = 3;
    const pendingFiles = files.filter((f) => f.status == "pending" || f.status == "error");

    for (let i = 0; i < pendingFiles.length; i += batchSize) {
      const batch = pendingFiles.slice(i, i + batchSize);
      const promises = batch.map((file) => uploadFile(file));

      await Promise.allSettled(promises);
    }

    setUploading(false);
  };

  const handleRemove = (fileId: string) => {
    setFiles((prev) => prev.filter((file) => file.id !== fileId));
    cleanupProgress(fileId);
  };

  const getUploadStats = () => {
    const total = files.length;
    const success = files.filter((f) => f.status == "success").length;
    const error = files.filter((f) => f.status == "error").length;
    const uploading = files.filter((f) => f.status == "uploading").length;

    return { total, success, error, uploading };
  };

  const stats = getUploadStats();
  const hasFailedFiles = stats.error > 0;
  const allCompleted = stats.success + stats.error === stats.total && stats.total > 0;

  return (
    <PageLayout>
      <div className="w-full max-w-2xl mx-auto p-6">
        {files.length > 0 && (
          <div className="mb-6">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-gray-700 dark:text-gray-300">
                Selected {stats.total} files
              </h2>

              <div className="flex items-center gap-4">
                {(uploading || allCompleted) && (
                  <div className="text-sm text-gray-600 dark:text-gray-400">
                    <span className="text-green-600 dark:text-green-400">✓ {stats.success}</span>
                    {stats.error > 0 && (
                      <span className="text-red-600 dark:text-red-400 ml-2">✗ {stats.error}</span>
                    )}
                    {stats.uploading > 0 && (
                      <span className="text-blue-600 dark:text-blue-400 ml-2">
                        ⏳ {stats.uploading}
                      </span>
                    )}
                  </div>
                )}
                <p onClick={() => setFiles([])} className="text-red-400 cursor-pointer">
                  clear
                </p>
              </div>
            </div>

            <div className="space-y-3 max-h-64 overflow-y-auto">
              {files.map((file) => (
                <FileCard
                  key={file.id}
                  file={file}
                  onRemove={() => handleRemove(file.id)}
                  progress={uploadProgress.get(file.id) || 0}
                />
              ))}
            </div>
          </div>
        )}

        <div className="space-y-4">
          <FileInput onChange={handleFileChange} />

          {files.length > 0 && (
            <div className="space-y-2">
              <button
                className="w-full bg-blue-500 hover:bg-blue-600 disabled:bg-gray-400 
                          text-white font-semibold py-3 px-4 rounded-lg 
                          transition-colors duration-200 disabled:cursor-not-allowed"
                onClick={handleClick}
                disabled={uploading}
              >
                {uploading
                  ? `Uploading... (${stats.success + stats.error}/${stats.total})`
                  : hasFailedFiles
                    ? `Retry Failed Files (${stats.error})`
                    : `Upload ${files.filter((f) => f.status == "pending").length} files`}
              </button>

              {allCompleted && stats.success > 0 && (
                <div className="text-center text-sm text-green-600 dark:text-green-400">
                  {stats.success} files uploaded successfully!
                </div>
              )}
            </div>
          )}
        </div>
      </div>
    </PageLayout>
  );
};

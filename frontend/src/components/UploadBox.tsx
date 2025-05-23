import { useState, type ChangeEvent } from "react";
import axios from "axios";
import SparkMD5 from "spark-md5";
import uploadIcon from "../assets/uploadIcon.svg";
import { SERVER } from "../constants";
import { useMutation } from "@tanstack/react-query";
import { UploadedList } from "./UploadedList";

/*
 *  TODO:
 *   - need better error handling
 */

export const UploadBox = () => {
  const [selectedFile, setSelectedFile] = useState<File[]>([]);
  // Add a separate state for tracking form data preparation
  const [isPreparingFiles, setIsPreparingFiles] = useState(false);

  const handleFileSelect = (ev: ChangeEvent<HTMLInputElement>) => {
    if (ev.target.files) {
      const files = Array.from(ev.target.files);
      setSelectedFile((prev) => [...prev, ...files]);
    }
  };

  const hashMD5 = async (blob: Blob): Promise<string> => {
    const buffer = await blob.arrayBuffer();
    const spark = new SparkMD5.ArrayBuffer();
    spark.append(buffer);
    return spark.end();
  };

  const buildFormData = async (): Promise<FormData> => {
    if (selectedFile.length === 0) return Promise.reject("files not found");

    const chunkSize = 3 * 1024 * 1024; // 3mb
    const formData = new FormData();

    for (const file of selectedFile) {
      const totalChunks = Math.ceil(file.size / chunkSize);

      for (let i = 0; i < totalChunks; i++) {
        const start = i * chunkSize;
        const end = Math.min(file.size, start + chunkSize);
        const chunk = file.slice(start, end);
        const hash = await hashMD5(chunk);

        const filename = btoa(encodeURIComponent(file.name)); // encoded
        formData.append(`file-${filename}-${hash}-${i}`, chunk);
      }
    }

    return formData;
  };

  const handleSubmit = async () => {
    if (selectedFile.length === 0 || isPreparingFiles || isPending) return;

    try {
      // Set preparing state to true before starting the process
      setIsPreparingFiles(true);

      const formData = await buildFormData();
      mutate(formData);
    } catch (error) {
      console.error("Error preparing upload:", error);
      // Could add error state handling here
    } finally {
      // Reset preparing state when done (only if mutation didn't start)
      if (!isPending) {
        setIsPreparingFiles(false);
      }
    }
  };

  const { mutate, isPending, isError, error, isSuccess } = useMutation({
    mutationKey: ["upload files", selectedFile],
    mutationFn: async (formData: FormData) =>
      await axios.post(`${SERVER}/files/upload`, formData, {
        withCredentials: true,
        headers: { "Content-Type": "multipart/form-data" },
      }),
    onSuccess: () => {
      // Clear file selection after successful upload
      setSelectedFile([]);
      setIsPreparingFiles(false);
    },
    onError: () => {
      setIsPreparingFiles(false);
    },
  });

  // Determine if any loading is happening (either preparing or uploading)
  const isLoading = isPreparingFiles || isPending;

  // Get the appropriate loading text
  const loadingText = isPreparingFiles ? "Preparing files..." : "Uploading...";

  return (
    <div className="h-[calc(100vh-100px)] flex flex-col justify-center items-center p-10">
      {/* Display upload status message */}
      {isError && (
        <div className="w-full max-w-3xl mb-4 p-3 bg-red-100 text-red-700 rounded-lg">
          Upload failed: {error instanceof Error ? error.message : "Unknown error"}
        </div>
      )}

      {isSuccess && (
        <div className="w-full max-w-3xl mb-4 p-3 bg-green-100 text-green-700 rounded-lg">
          Files uploaded successfully!
        </div>
      )}

      {/* File list */}
      {selectedFile.length > 0 && <UploadedList list={selectedFile} setList={setSelectedFile} />}

      <div className="w-full max-w-3xl">
        {/* File upload area */}
        <label
          htmlFor="file-upload"
          className={`flex flex-col items-center justify-center w-full h-32 px-4 py-6 border-2 border-dashed 
          ${isLoading ? "border-gray-400 bg-gray-100 cursor-not-allowed" : "border-gray-300 bg-gray-50 hover:bg-gray-100 cursor-pointer"} 
          rounded-lg transition-colors`}
        >
          <div className="flex flex-col items-center justify-center">
            {isLoading ? (
              <div className="flex flex-col items-center">
                <div className="w-8 h-8 mb-3">
                  <svg
                    className="animate-spin h-8 w-8 text-blue-500"
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                  >
                    <circle
                      className="opacity-25"
                      cx="12"
                      cy="12"
                      r="10"
                      stroke="currentColor"
                      strokeWidth="4"
                    ></circle>
                    <path
                      className="opacity-75"
                      fill="currentColor"
                      d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                    ></path>
                  </svg>
                </div>
                <p className="text-sm text-gray-500">{loadingText}</p>
              </div>
            ) : (
              <>
                <img src={uploadIcon} className="w-8 h-8 mb-3 text-gray-500" alt="Upload icon" />
                <p className="mb-2 text-sm text-gray-500">
                  <span className="font-semibold">Select or drag files</span>
                </p>
              </>
            )}
          </div>
          <input
            id="file-upload"
            type="file"
            multiple
            className="hidden"
            onChange={handleFileSelect}
            disabled={isLoading}
          />
        </label>

        {/* Upload button */}
        <button
          className={`w-full py-2 rounded-lg mt-5 flex items-center justify-center transition-colors ${
            isLoading || selectedFile.length === 0
              ? "bg-gray-400 cursor-not-allowed text-white"
              : "bg-blue-500 hover:bg-blue-600 text-white cursor-pointer"
          }`}
          type="button"
          onClick={handleSubmit}
          disabled={isLoading || selectedFile.length === 0}
        >
          {isLoading ? (
            <>
              <svg
                className="animate-spin -ml-1 mr-2 h-4 w-4 text-white"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
              >
                <circle
                  className="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  strokeWidth="4"
                ></circle>
                <path
                  className="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                ></path>
              </svg>
              {loadingText}
            </>
          ) : (
            "Upload"
          )}
        </button>
      </div>
    </div>
  );
};

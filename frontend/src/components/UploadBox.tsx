import { useState, type ChangeEvent } from "react";
import axios from "axios";
import SparkMD5 from "spark-md5";
import uploadIcon from "../assets/uploadIcon.svg";
import { SERVER } from "../constants";
import { useMutation } from "@tanstack/react-query";

/*
 *  TODO:
 *   - need better error handling
 *   - need UI for upload progress
 */

export const UploadBox = () => {
  const [selectedFile, setSelectedFile] = useState<File[]>([]);

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

  const { mutate } = useMutation({
    mutationKey: ["upload files", selectedFile],
    mutationFn: async (formData: FormData) =>
      await axios.post(`${SERVER}/files/upload`, formData, {
        withCredentials: true,
        headers: { "Content-Type": "multipart/form-data" },
      }),
  });

  return (
    <>
      {/* TODO: consider extracting thumbnails for list of previews */}
      {selectedFile.length > 0 && (
        <div className="flex flex-col gap-3">
          {selectedFile.map((file, idx) => (
            <p key={idx}>{file.name}</p>
          ))}
        </div>
      )}

      <div className="flex flex-col items-center justify-center w-full mt-5">
        <label
          htmlFor="file-upload"
          className="flex flex-col items-center justify-center w-full h-32 px-4 py-6 border-2 border-dashed border-gray-300 rounded-lg cursor-pointer bg-gray-50 hover:bg-gray-100 transition-colors"
        >
          <div className="flex flex-col items-center justify-center">
            <img src={uploadIcon} className="w-8 h-8 mb-3 text-gray-500" />
            <p className="mb-2 text-sm text-gray-500">
              <span className="font-semibold">Select or drag files</span>
            </p>
          </div>
          <input
            id="file-upload"
            type="file"
            multiple
            className="hidden"
            onChange={handleFileSelect}
          />
        </label>

        <button
          className="btn mt-5"
          type="submit"
          onClick={async () => mutate(await buildFormData())}
        >
          Upload
        </button>
      </div>
    </>
  );
};

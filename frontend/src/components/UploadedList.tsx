import { type Dispatch, type SetStateAction, useMemo } from "react";
import { formatFileSize } from "../utils";

const FileIcon = ({ type }: { type: string }) => {
  if (type.startsWith("image/")) {
    return (
      <svg
        className="w-5 h-5 text-blue-500"
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth={2}
          d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
        />
      </svg>
    );
  } else if (type.startsWith("video/")) {
    return (
      <svg
        className="w-5 h-5 text-red-500"
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth={2}
          d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"
        />
      </svg>
    );
  } else if (type.startsWith("audio/")) {
    return (
      <svg
        className="w-5 h-5 text-purple-500"
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth={2}
          d="M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3"
        />
      </svg>
    );
  } else if (
    type === "application/pdf" ||
    type === "text/plain" ||
    type === "application/msword" ||
    type.includes("document")
  ) {
    return (
      <svg
        className="w-5 h-5 text-yellow-500"
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth={2}
          d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
        />
      </svg>
    );
  } else if (
    type.includes("zip") ||
    type.includes("compressed") ||
    type === "application/x-rar-compressed"
  ) {
    return (
      <svg
        className="w-5 h-5 text-green-500"
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth={2}
          d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4"
        />
      </svg>
    );
  } else {
    return (
      <svg
        className="w-5 h-5 text-gray-500"
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth={2}
          d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z"
        />
      </svg>
    );
  }
};

const ImagePreview = ({ file }: { file: File }) => {
  if (!file.type.startsWith("image/")) {
    return null;
  }

  const imageUrl = URL.createObjectURL(file);

  return (
    <div className="relative w-8 h-8 overflow-hidden rounded mr-3">
      <img
        src={imageUrl}
        className="object-cover w-full h-full"
        onLoad={() => URL.revokeObjectURL(imageUrl)}
        alt={file.name}
      />
    </div>
  );
};

export const UploadedList = ({
  list,
  setList,
}: {
  list: File[];
  setList: Dispatch<SetStateAction<File[]>>;
}) => {
  const handleClick = (filename: string) => {
    const filtered = list.filter((file) => file.name !== filename);
    setList([...filtered]);
  };

  // Calculate total file size
  const totalSize = useMemo(() => {
    return list.reduce((total, file) => total + file.size, 0);
  }, [list]);

  return (
    <div className="w-full max-w-3xl mb-6">
      <div className="flex justify-between items-center mb-3">
        <div>
          <h3 className="text-lg font-semibold text-blue-100">
            Files Selected <span className="text-blue-300">({list.length})</span>
          </h3>
          {list.length > 0 && (
            <p className="text-sm text-gray-400">
              Total size: <span className="text-blue-300">{formatFileSize(totalSize)}</span>
            </p>
          )}
        </div>
        {list.length > 0 && (
          <button
            className="text-sm text-red-400 hover:text-red-300 cursor-pointer transition-colors"
            onClick={() => setList([])}
          >
            Clear All
          </button>
        )}
      </div>

      {list.length > 0 ? (
        <div className="bg-gray-800 bg-opacity-50 rounded-lg p-2 max-h-40 overflow-y-auto">
          {list.map((file, idx) => (
            <div
              key={idx}
              className="flex items-center p-2 rounded-md hover:bg-gray-700 transition-colors mb-1 group"
            >
              <div className="flex items-center flex-1 min-w-0">
                {file.type.startsWith("image/") && <ImagePreview file={file} />}

                <FileIcon type={file.type} />

                <div className="ml-3 flex flex-col min-w-0">
                  <span className="text-sm text-white truncate">{file.name}</span>
                  <span className="text-xs text-gray-400">{formatFileSize(file.size)}</span>
                </div>
              </div>

              <button
                className="ml-2 p-1.5 rounded-full text-gray-400 text-sm opacity-0 group-hover:opacity-100 hover:text-white transition-all cursor-pointer"
                onClick={() => handleClick(file.name)}
                title="Remove file"
              >
                <svg
                  className="w-4 h-4"
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M6 18L18 6M6 6l12 12"
                  />
                </svg>
              </button>
            </div>
          ))}
        </div>
      ) : (
        <div className="flex items-center justify-center bg-gray-800 bg-opacity-30 rounded-lg p-6 border border-dashed border-gray-600">
          <p className="text-gray-400 text-center">No files selected yet</p>
        </div>
      )}
    </div>
  );
};

import { useState, useEffect, type SetStateAction, type Dispatch } from "react";
import { SERVER } from "../constants";
import { MenuIcon, X, Download, Rotate3D, ZoomIn, ZoomOut, Share2 } from "lucide-react";
import { formatFileSize } from "../utils";

export interface Iitem {
  filename: string;
  contenttype: string;
  uploadedat: string;
  size: number;
}

export interface IfileMetadata {
  totalpages: number;
  totalitems: number;
  items: Iitem[];
}

const ViewFileModal = ({
  src,
  contentType,
  fileName,
  setShowModal,
}: {
  src: string;
  contentType: string;
  fileName: string;
  setShowModal: Dispatch<SetStateAction<boolean>>;
}) => {
  const [openMenu, setOpenMenu] = useState(false);
  const [zoomLevel, setZoomLevel] = useState(1);
  const [rotation, setRotation] = useState(0);

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") setShowModal(false);
    };

    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [setShowModal]);

  const handleRotate = () => {
    setRotation((prev) => (prev + 90) % 360);
  };

  const zoomIn = () => setZoomLevel((prev) => Math.min(prev + 0.25, 3));
  const zoomOut = () => setZoomLevel((prev) => Math.max(prev - 0.25, 0.5));

  const handleDownload = () => {
    const link = document.createElement("a");
    link.href = src;
    link.download = fileName || "download";
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };

  return (
    <div className="fixed inset-0 z-50 bg-black bg-opacity-95 flex flex-col">
      <div className="flex h-16 justify-between items-center p-4 bg-gray-900 bg-opacity-60">
        <div className="flex items-center">
          <button
            onClick={() => setOpenMenu(!openMenu)}
            className="p-2 rounded-full hover:bg-gray-700 text-white transition-colors"
            aria-label="Toggle menu"
          >
            <MenuIcon size={24} />
          </button>
          <h4 className="text-white ml-4 mb-6 truncate max-w-[200px] md:max-w-[400px]">
            {fileName?.split("/").pop() || "File viewer"}
          </h4>
        </div>

        <button
          onClick={() => setShowModal(false)}
          className="p-2 rounded-full hover:bg-gray-700 text-white transition-colors"
          aria-label="Close modal"
        >
          <X size={24} />
        </button>
      </div>

      {openMenu && (
        <div className="absolute top-16 left-5 w-60 bg-gray-800 shadow-lg rounded-tr-lg rounded-br-lg z-10 py-2">
          <button
            className="w-full text-left px-4 py-3 hover:bg-gray-700 text-white flex items-center"
            onClick={zoomIn}
          >
            <ZoomIn size={20} className="mr-3" /> Zoom in
          </button>
          <button
            className="w-full text-left px-4 py-3 hover:bg-gray-700 text-white flex items-center"
            onClick={zoomOut}
          >
            <ZoomOut size={20} className="mr-3" /> Zoom out
          </button>
          <button
            className="w-full text-left px-4 py-3 hover:bg-gray-700 text-white flex items-center"
            onClick={handleRotate}
          >
            <Rotate3D size={20} className="mr-3" /> rotate
          </button>
          <button
            className="w-full text-left px-4 py-3 hover:bg-gray-700 text-white flex items-center"
            onClick={handleDownload}
          >
            <Download size={20} className="mr-3" /> download
          </button>
          <button className="w-full text-left px-4 py-3 hover:bg-gray-700 text-white flex items-center">
            <Share2 size={20} className="mr-3" /> share
          </button>
        </div>
      )}

      <div
        className="flex-1 overflow-hidden flex items-center justify-center p-4"
        onClick={() => openMenu && setOpenMenu(false)}
      >
        {contentType === "image" ? (
          <img
            className="max-h-full max-w-full object-contain transition-transform"
            src={src}
            alt={fileName}
            style={{
              transform: `scale(${zoomLevel}) rotate(${rotation}deg)`,
              transition: "transform 0.3s ease",
            }}
          />
        ) : contentType === "video" ? (
          <video className="max-h-full max-w-full" src={src} controls />
        ) : (
          <div className="text-white text-center">
            <div className="text-6xl mb-4">📄</div>
            <p>file format unsupported</p>
            <button
              className="mt-4 px-4 py-2 bg-blue-600 rounded hover:bg-blue-700 transition-colors"
              onClick={handleDownload}
            >
              download
            </button>
          </div>
        )}
      </div>

      {/* bottom control - shows only if content type is image */}
      {contentType === "image" && (
        <div className="flex justify-center p-4 bg-gray-900 bg-opacity-60">
          <div className="flex space-x-4">
            <button
              onClick={zoomOut}
              className="p-2 rounded-full hover:bg-gray-700 text-white transition-colors"
              aria-label="Zoom out"
            >
              <ZoomOut size={20} />
            </button>
            <span className="text-white flex items-center">{Math.round(zoomLevel * 100)}%</span>
            <button
              onClick={zoomIn}
              className="p-2 rounded-full hover:bg-gray-700 text-white transition-colors"
              aria-label="Zoom in"
            >
              <ZoomIn size={20} />
            </button>
            <button
              onClick={handleRotate}
              className="p-2 rounded-full hover:bg-gray-700 text-white transition-colors"
              aria-label="Rotate"
            >
              <Rotate3D size={20} />
            </button>
            <button
              onClick={handleDownload}
              className="p-2 rounded-full hover:bg-gray-700 text-white transition-colors"
              aria-label="Download"
            >
              <Download size={20} />
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

export const FileCard = ({ file }: { file: Iitem }) => {
  const [showModal, setShowModal] = useState(false);

  return (
    <>
      <div
        className="flex-shrink-0 w-full bg-gray-800 rounded-lg overflow-hidden shadow-lg mr-4 cursor-pointer transform transition-transform hover:scale-105"
        onClick={() => setShowModal(true)}
      >
        <div className="aspect-square bg-gray-900 overflow-hidden">
          {file.contenttype === "image" ? (
            <img
              src={`${SERVER}${file.filename}`}
              alt={file.filename}
              className="w-full h-full object-cover"
            />
          ) : file.contenttype === "video" ? (
            <div className="relative w-full h-full">
              <video src={`${SERVER}${file.filename}`} className="w-full h-full object-cover" />
              <div className="absolute inset-0 flex items-center justify-center">
                <div className="w-12 h-12 bg-black bg-opacity-50 rounded-full flex items-center justify-center">
                  <div className="w-0 h-0 border-t-8 border-t-transparent border-l-12 border-l-white border-b-8 border-b-transparent ml-1"></div>
                </div>
              </div>
            </div>
          ) : (
            <div className="w-full h-full flex items-center justify-center text-4xl">📁</div>
          )}
        </div>
        <div className="p-3">
          <p className="font-medium text-white truncate">{file.filename.split("/").pop()}</p>
          <div className="flex justify-between mt-2 text-xs text-gray-300">
            <span>{new Date(file.uploadedat).toLocaleDateString()}</span>
            <span>{formatFileSize(file.size || 0)}</span>
          </div>
        </div>
      </div>

      {showModal && (
        <ViewFileModal
          src={`${SERVER}${file.filename}`}
          contentType={file.contenttype}
          fileName={file.filename}
          setShowModal={setShowModal}
        />
      )}
    </>
  );
};

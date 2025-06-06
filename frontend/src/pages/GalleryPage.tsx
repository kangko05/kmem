import { PageLayout } from "../components";
import { axiosInstance } from "../utils";
import { SERVER } from "../constants";
import { useAuth } from "../hooks/useAuth";
import { useInfiniteQuery } from "react-query";
import { BeatLoader } from "react-spinners";
import { useState, useEffect, type ChangeEvent, type ReactNode } from "react";
import { Settings2 } from "lucide-react";

type Tsort = "date" | "name";
type Ttype = "all" | "image" | "video";

interface Ifile {
  originalName: string;
  mimeType: string;
  filePath: string;
}

export const Lightbox = ({
  idx,
  files,
  onClose,
}: {
  idx: number;
  files: Ifile[];
  onClose: () => void;
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
    <div
      className="fixed inset-0 bg-black bg-opacity-90 flex items-center justify-center z-50"
      onClick={onClose}
    >
      <Settings2
        className="absolute top-10 left-10 cursor-pointer"
        onClick={(ev) => {
          ev.stopPropagation();
          setOpenSettings(!openSettings);
        }}
      />

      {openSettings && (
        <ul className="absolute top-20 left-10 bg-red-100">
          <li>Delete</li>
        </ul>
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
            src={`${SERVER}/static${currentFile.filePath}`}
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
            <source src={`${SERVER}/static${currentFile.filePath}`} />
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

const RenderContents = ({
  idx,
  mimetype,
  src,
  onImageClick,
}: {
  idx: number;
  mimetype: string;
  src: string;
  name: string;
  onImageClick: (idx: number) => void;
}) => {
  const contentSrc = `${SERVER}/static${src}`;

  if (mimetype.includes("image")) {
    return (
      <img
        key={idx}
        src={contentSrc}
        className="w-full h-full object-cover"
        onClick={() => onImageClick(idx)}
      />
    );
  }

  if (mimetype.includes("video")) {
    return (
      <video
        muted
        className="w-full h-full object-cover cursor-pointer"
        onClick={() => onImageClick(idx)}
      >
        <source key={idx} src={contentSrc}></source>
      </video>
    );
  }
};

const GalleryItem = ({
  it,
  idx,
  onImageClick,
}: {
  it: Ifile;
  idx: number;
  onImageClick: (idx: number) => void;
}) => {
  return (
    <div
      key={idx}
      className="group relative aspect-square bg-gray-100 dark:bg-gray-800 rounded-lg overflow-hidden hover:shadow-lg transition-shadow duration-200 cursor-pointer"
    >
      <div className="w-full h-full">
        {RenderContents({
          idx: idx,
          mimetype: it.mimeType,
          src: it.filePath,
          name: it.originalName,
          onImageClick: onImageClick,
        })}
      </div>

      <div className="absolute bottom-0 left-0 right-0 bg-black bg-opacity-50 text-white p-2 opacity-0 group-hover:opacity-100 transition-opacity duration-200">
        <p className="truncate">{it.originalName}</p>
      </div>
    </div>
  );
};

const GalleryView = ({
  fetchNextPage,
  isFetchingNextPage,
  hasNextPage,
  files,
  onImageClick,
}: {
  fetchNextPage: () => void;
  isFetchingNextPage: boolean;
  hasNextPage: boolean;
  files: Ifile[];
  onImageClick: (idx: number) => void;
}) => {
  return (
    <div className="w-full h-[75%] overflow-auto max-w-7xl mx-auto p-6">
      <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4">
        {files?.map((it, idx) => (
          <GalleryItem key={idx} it={it} idx={idx} onImageClick={onImageClick} />
        ))}
      </div>

      {files.length == 0 && (
        <div className="text-center py-12">
          <p className="text-gray-500 dark:text-gray-400">No files uploaded yet</p>
        </div>
      )}

      {hasNextPage && (
        <div className="text-center mt-8">
          {isFetchingNextPage ? (
            <div className="flex justify-center items-center py-4">
              <BeatLoader color="#3b82f6" size={10} />
            </div>
          ) : (
            <button
              onClick={() => fetchNextPage()}
              className="px-8 py-3 bg-blue-500 hover:bg-blue-600 text-white font-medium rounded-lg transition-colors duration-200 shadow-md hover:shadow-lg"
            >
              Load More
            </button>
          )}
        </div>
      )}
    </div>
  );
};

const DropDown = ({
  name,
  onChange,
  children,
}: {
  name: string;
  onChange: (ev: ChangeEvent<HTMLSelectElement>) => void;
  children: ReactNode;
}) => {
  return (
    <select
      name={name}
      onChange={(ev) => onChange(ev)}
      className="w-fit px-3 py-2 border border-gray-300 rounded-lg bg-white dark:bg-gray-700 dark:border-gray-600 text-gray-700 dark:text-gray-300"
    >
      {children}
    </select>
  );
};

export const GalleryPage = () => {
  useAuth();

  const [sort, setSort] = useState<Tsort>("date");
  const [type, setType] = useState<Ttype>("all");
  const [lightboxIdx, setLightboxIdx] = useState<number | null>(null);

  const fetchPage = async (pageParam: number) => {
    const resp = await axiosInstance.get(
      `${SERVER}/files?limit=20&page=${pageParam}&sort=${sort}&type=${type}`
    );
    return resp.data;
  };

  const {
    data,
    fetchNextPage,
    hasNextPage,
    isLoading: isFetchingNextPage,
  } = useInfiniteQuery({
    queryKey: ["gallery", sort, type],
    queryFn: ({ pageParam = 0 }) => fetchPage(pageParam),
    getNextPageParam: (lastPage) => (lastPage.hasNext ? lastPage.nextPage : undefined),
  });

  const files = data?.pages.flatMap((page) => page.files) || [];

  const handleImageClick = (idx: number) => {
    setLightboxIdx(idx);
  };

  const closeLightbox = () => {
    setLightboxIdx(null);
  };

  return (
    <PageLayout>
      <div className="w-full h-[80%] flex flex-col gap-5">
        <div className="w-full max-w-7xl mx-auto mt-5 flex justify-end gap-3">
          <DropDown name="type" onChange={(ev) => setType(ev.currentTarget.value as Ttype)}>
            <option value="all">all</option>
            <option value="image">image</option>
            <option value="video">video</option>
          </DropDown>

          <DropDown name="sort" onChange={(ev) => setSort(ev.currentTarget.value as Tsort)}>
            <option value="date">Latest First</option>
            <option value="name">Name</option>
          </DropDown>

          <div className="mr-3" />
        </div>

        <GalleryView
          files={files}
          fetchNextPage={() => fetchNextPage()}
          hasNextPage={hasNextPage || false}
          isFetchingNextPage={isFetchingNextPage}
          onImageClick={handleImageClick}
        />
      </div>

      {lightboxIdx != null && <Lightbox idx={lightboxIdx} files={files} onClose={closeLightbox} />}
    </PageLayout>
  );
};

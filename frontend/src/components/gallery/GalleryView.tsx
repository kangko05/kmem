import { GalleryItem } from "./GalleryItem";
import { BeatLoader } from "react-spinners";
import { type Ifile } from "./types";

export const GalleryView = ({
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

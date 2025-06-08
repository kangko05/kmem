import { SERVER } from "../../constants";
import { type Ifile } from "./types";

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
  const contentSrc = `${SERVER}${src}`;

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
    if (src.endsWith("jpg")) {
      return (
        <img
          key={idx}
          src={contentSrc}
          className="w-full h-full object-cover"
          onClick={() => onImageClick(idx)}
        />
      );
    }

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

export const GalleryItem = ({
  it,
  idx,
  onImageClick,
}: {
  it: Ifile;
  idx: number;
  onImageClick: (idx: number) => void;
}) => {
  let src = it.filePath;

  if (it.thumbnails?.small) {
    src = it.thumbnails.small.filePath;
  } else if (it.thumbnails?.medium) {
    src = it.thumbnails.medium.filePath;
  }

  return (
    <div
      key={idx}
      className="group relative aspect-square bg-gray-100 dark:bg-gray-800 rounded-lg overflow-hidden hover:shadow-lg transition-shadow duration-200 cursor-pointer"
    >
      <div className="w-full h-full">
        {RenderContents({
          idx: idx,
          mimetype: it.mimeType,
          src: src,
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

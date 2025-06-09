import { PageLayout } from "../components";
import { useAuth } from "../hooks/useAuth";
import { StorageBar, type Ifile, GalleryItem } from "../components";
import { axiosInstance } from "../utils";
import { SERVER } from "../constants";
import { useQuery } from "react-query";

export const HomePage = () => {
  useAuth();

  const getUsage = async () => {
    const resp = await axiosInstance.get(`${SERVER}/stats/usage`);
    return resp.data;
  };

  const getRecentUploads = async () => {
    const resp = await axiosInstance.get(`${SERVER}/files?limit=8&sort=date`);
    return resp.data;
  };

  const { data: usage } = useQuery({ queryKey: ["usage"], queryFn: getUsage });
  const { data: recent } = useQuery({ queryKey: ["recent uploads"], queryFn: getRecentUploads });

  return (
    <PageLayout>
      <div className="flex-center flex-col max-w-6xl mx-auto p-6 space-y-8">
        <div className="mt-20 sm:hidden" />
        <div className="text-center">
          <h1 className="text-3xl font-bold text-gray-800 dark:text-white">
            Hello <span className="text-blue-300">{usage?.data.username || "there"}</span>
          </h1>
          <p className="text-gray-600 dark:text-gray-400 mt-2">
            Total {usage?.data.count || 0} files stored
          </p>
        </div>
        <StorageBar
          percentage={(usage?.data.size / 1e12) * 100}
          used={usage?.data.readableSize}
          total={"1 TB"}
        />
        <div className="w-full">
          <h2 className="text-xl font-semibold mb-4 text-gray-800 dark:text-white">
            Recent uploads
          </h2>
          <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 xl:grid-cols-8 gap-4">
            {recent?.data.files?.length > 0 ? (
              recent.data.files.map((it: Ifile, idx: number) => (
                <GalleryItem it={it} idx={idx} key={idx} onImageClick={() => {}} />
              ))
            ) : (
              <div className="col-span-full text-center text-gray-500 py-8">
                No files uploaded yet
              </div>
            )}
          </div>
        </div>
      </div>
    </PageLayout>
  );
};

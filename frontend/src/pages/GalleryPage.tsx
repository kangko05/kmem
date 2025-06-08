import { PageLayout, LightBox, GalleryView, SearchBar, DropDown } from "../components";
import { axiosInstance } from "../utils";
import { SERVER } from "../constants";
import { useAuth } from "../hooks/useAuth";
import { useInfiniteQuery } from "react-query";
import { useState, useEffect } from "react";

type Tsort = "date" | "name";
type Ttype = "all" | "image" | "video";

export const GalleryPage = () => {
  useAuth();

  const [sort, setSort] = useState<Tsort>("date");
  const [type, setType] = useState<Ttype>("all");
  const [search, setSearch] = useState("");
  const [debouncedSearch, setDebouncedSearch] = useState("");
  const [lightboxIdx, setLightboxIdx] = useState<number | null>(null);

  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedSearch(search);
    }, 300);
    return () => clearTimeout(timer);
  }, [search]);

  const fetchPage = async (pageParam: number) => {
    let reqUrl = `${SERVER}/files?limit=20&page=${pageParam}&sort=${sort}&type=${type}`;

    if (search.length > 2) reqUrl += `&search=${debouncedSearch}`;

    const resp = await axiosInstance.get(reqUrl);

    return resp.data;
  };

  const {
    data,
    refetch,
    fetchNextPage,
    hasNextPage,
    isLoading: isFetchingNextPage,
  } = useInfiniteQuery({
    queryKey: ["gallery", sort, type, debouncedSearch],
    queryFn: ({ pageParam = 0 }) => fetchPage(pageParam),
    getNextPageParam: (lastPage) => (lastPage.data.hasNext ? lastPage.data.nextPage : undefined),
  });

  const files = data?.pages.flatMap((page) => page.data.files || []) || [];

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
          <div className="ml-1" />

          <SearchBar search={search} setSearch={setSearch} />

          <DropDown name="type" onChange={(ev) => setType(ev.currentTarget.value as Ttype)}>
            <option value="all">all</option>
            <option value="image">image</option>
            <option value="video">video</option>
          </DropDown>

          <DropDown name="sort" onChange={(ev) => setSort(ev.currentTarget.value as Tsort)}>
            <option value="date">Latest First</option>
            <option value="name">Name</option>
          </DropDown>

          <div className="sm:mr-3" />
        </div>

        <GalleryView
          files={files}
          fetchNextPage={() => fetchNextPage()}
          hasNextPage={hasNextPage || false}
          isFetchingNextPage={isFetchingNextPage}
          onImageClick={handleImageClick}
        />
      </div>

      {lightboxIdx != null && (
        <LightBox
          idx={lightboxIdx}
          files={files}
          onClose={closeLightbox}
          refetch={() => refetch()}
        />
      )}
    </PageLayout>
  );
};

import { useEffect, useRef, useState, type ChangeEvent } from "react";
import { axiosInstance } from "../utils/AxiosIntstance";
import { useAuthCheck } from "../hooks";
import { useQuery, useInfiniteQuery } from "@tanstack/react-query";
import { PageLayout, Spinner, Pagination } from "../components";
import { FileCard, type Iitem } from "../components/FileCard";
import { Link } from "react-router";
import upperArrow from "../assets/upperArrow.svg";

type sortOption = "date" | "name";

export const HomePage = () => {
  useAuthCheck();
  const divRef = useRef<HTMLDivElement | null>(null);
  const [sortBy, setSortBy] = useState<sortOption>("date");
  const [currentPage, setCurrentPage] = useState<number>(1);
  const [itemsPerPage, setItemsPerPage] = useState<number>(12);
  const [isMobile, setIsMobile] = useState<boolean>(false);

  const loadMoreRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    const checkMobile = () => {
      setIsMobile(window.innerWidth < 768);
    };

    checkMobile();

    window.addEventListener("resize", checkMobile);
    return () => window.removeEventListener("resize", checkMobile);
  }, []);

  const calculateItemsPerPage = () => {
    if (!divRef.current) return 12;

    const containerWidth = divRef.current.clientWidth;

    if (containerWidth < 640) return 6;
    else if (containerWidth < 768) return 8;
    else if (containerWidth < 1024) return 6;
    else if (containerWidth < 1280) return 8;
    else return 10;
  };

  useEffect(() => {
    if (!divRef.current) return;

    const newItemsPerPage = calculateItemsPerPage();
    if (newItemsPerPage !== itemsPerPage) {
      setItemsPerPage(newItemsPerPage);
    }

    const handleResize = () => {
      const newItemsPerPage = calculateItemsPerPage();
      if (newItemsPerPage !== itemsPerPage) {
        setItemsPerPage(newItemsPerPage);
      }
    };

    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, [divRef, itemsPerPage]);

  const getRecentItems = async (sort: sortOption, page: number, perPage: number) => {
    const resp = await axiosInstance.get(
      `/files/items?sort=${sort}&page=${page}&itemsPerPage=${perPage}`
    );

    return resp.data;
  };

  const getInfiniteItems = async ({ pageParam = 1 }) => {
    const resp = await axiosInstance.get(
      `/files/items?sort=${sortBy}&page=${pageParam}&itemsPerPage=${itemsPerPage}`
    );
    return {
      items: resp.data.items,
      page: pageParam,
      totalPages: resp.data.totalpages,
    };
  };

  const handleSortChange = (ev: ChangeEvent<HTMLSelectElement>) => {
    const newSort: sortOption = ev.target.value as sortOption;
    setSortBy(newSort);
    setCurrentPage(1);
    if (isMobile) {
      infiniteQueryRefetch();
    }
  };

  const {
    data: pageData,
    isSuccess: pageSuccess,
    isLoading: pageLoading,
    isError: pageError,
    error: pageErrorData,
    refetch: pageRefetch,
  } = useQuery({
    queryKey: ["recent items", sortBy, currentPage, itemsPerPage],
    queryFn: () => getRecentItems(sortBy, currentPage, itemsPerPage),
    enabled: !isMobile && itemsPerPage > 0,
  });

  const {
    data: infiniteData,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isSuccess: infiniteSuccess,
    isLoading: infiniteLoading,
    isError: infiniteError,
    error: infiniteErrorData,
    refetch: infiniteQueryRefetch,
  } = useInfiniteQuery({
    queryKey: ["infinite items", sortBy, itemsPerPage],
    queryFn: getInfiniteItems,
    getNextPageParam: (lastPage) => {
      return lastPage.page < lastPage.totalPages ? lastPage.page + 1 : undefined;
    },
    enabled: isMobile && itemsPerPage > 0,
  });

  useEffect(() => {
    if (!isMobile || !loadMoreRef.current) return;

    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && hasNextPage && !isFetchingNextPage) {
          fetchNextPage();
        }
      },
      { threshold: 0.1 }
    );

    observer.observe(loadMoreRef.current);
    return () => observer.disconnect();
  }, [isMobile, hasNextPage, isFetchingNextPage, fetchNextPage]);

  useEffect(() => {
    if (itemsPerPage > 0) {
      if (isMobile) {
        infiniteQueryRefetch();
      } else {
        pageRefetch();
      }
    }
  }, [itemsPerPage, isMobile, pageRefetch, infiniteQueryRefetch]);

  const infiniteItems = infiniteData?.pages.flatMap((page) => page.items) || [];

  const isLoading = isMobile ? infiniteLoading : pageLoading;

  const isError = isMobile ? infiniteError : pageError;
  const error = isMobile ? infiniteErrorData : pageErrorData;

  const isSuccess = isMobile ? infiniteSuccess : pageSuccess;

  return (
    <PageLayout>
      <div ref={divRef} className="p-4 w-full">
        <div className="flex justify-between items-center mb-4">
          {/* sort option */}
          <div className="flex gap-2">
            <select
              className="px-2 py-1 border rounded bg-gray-700 text-white"
              value={sortBy}
              onChange={(ev) => handleSortChange(ev)}
            >
              <option value="date">date</option>
              <option value="name">name</option>
            </select>
          </div>
        </div>

        {/* loading */}
        {isLoading && !isFetchingNextPage && (
          <div className="flex justify-center my-10">
            <Spinner loading={true} />
          </div>
        )}

        {/* error */}
        {isError && (
          <div className="p-3 bg-red-900 bg-opacity-30 border border-red-800 text-red-300 rounded mb-4">
            {error?.message}
          </div>
        )}

        {/* no data  */}
        {isSuccess &&
          ((isMobile && infiniteItems.length === 0) ||
            (!isMobile && (!pageData?.items || pageData.items.length === 0))) && (
            <div className="text-center py-10">
              <p className="text-gray-400">no fils have been uploaded yet</p>
              <Link
                className="mt-4 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
                to="/upload"
              >
                go to upload
              </Link>
            </div>
          )}

        {/* grid+pagination for desktop */}
        {!isMobile && isSuccess && pageData?.items && pageData.items.length > 0 && (
          <>
            <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-3">
              {pageData.items.map((file: Iitem, idx: number) => (
                <FileCard key={idx} file={file} />
              ))}
            </div>

            {/* pagination */}
            {pageData.totalpages > 1 && (
              <Pagination
                currentPage={currentPage}
                totalPages={pageData.totalpages}
                onPageChange={setCurrentPage}
                className="md: flex mt-4"
              />
            )}
          </>
        )}

        {/* infinite scroll for mobile */}
        {isMobile && isSuccess && infiniteItems.length > 0 && (
          <>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
              {infiniteItems.map((file: Iitem, idx: number) => (
                <FileCard key={idx} file={file} />
              ))}
            </div>

            <div ref={loadMoreRef} className="flex justify-center items-center py-4 mt-2">
              {isFetchingNextPage && <Spinner loading={true} />}
              {!hasNextPage && infiniteItems.length > 0 && (
                <p className="text-gray-400 text-sm">no more files</p>
              )}
            </div>

            <button
              onClick={() => {
                if (divRef.current) {
                  divRef.current.scrollIntoView({ behavior: "smooth" });
                }
              }}
              className="fixed bottom-6 right-6 bg-blue-600 text-white rounded-full w-12 h-12 flex items-center justify-center shadow-lg hover:bg-blue-700 transition-colors cursor-pointer"
              aria-label="Scroll to top"
            >
              <img src={upperArrow} className="invert" />
            </button>
          </>
        )}
      </div>
    </PageLayout>
  );
};

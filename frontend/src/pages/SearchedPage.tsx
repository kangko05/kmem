import { useEffect, useRef, useState } from "react";
import { axiosInstance } from "../utils/AxiosIntstance";
import { useAuthCheck } from "../hooks";
import { useQuery, useInfiniteQuery } from "@tanstack/react-query";
import { PageLayout, Spinner, Pagination } from "../components";
import { FileCard, type Iitem } from "../components/FileCard";
import { useLocation, useNavigate } from "react-router";
import upperArrow from "../assets/upperArrow.svg";

// Helper function to check if the search query is for a content type
const isContentTypeSearch = (query: string): boolean => {
  const contentTypeKeywords = [
    "image",
    "photo",
    "picture",
    "video",
    "movie",
    "audio",
    "music",
    "document",
    "doc",
    "text",
    "txt",
    "zip",
    "archive",
  ];

  return contentTypeKeywords.includes(query.toLowerCase());
};

export const SearchedPage = () => {
  useAuthCheck();
  const divRef = useRef<HTMLDivElement | null>(null);
  const location = useLocation();
  const navigate = useNavigate();
  const [searchQuery, setSearchQuery] = useState<string>("");
  const [currentPage, setCurrentPage] = useState<number>(1);
  const [itemsPerPage, setItemsPerPage] = useState<number>(12);
  const [isMobile, setIsMobile] = useState<boolean>(false);
  const loadMoreRef = useRef<HTMLDivElement | null>(null);

  // Get search query from URL params
  useEffect(() => {
    const params = new URLSearchParams(location.search);
    const query = params.get("q");
    if (query) {
      setSearchQuery(query);
    }
  }, [location.search]);

  // Detect mobile device
  useEffect(() => {
    const checkMobile = () => {
      setIsMobile(window.innerWidth < 768);
    };

    checkMobile();

    window.addEventListener("resize", checkMobile);
    return () => window.removeEventListener("resize", checkMobile);
  }, []);

  useEffect(() => {
    setCurrentPage(1); // set current page to default value whenever search query changes
  }, [searchQuery]);

  // Calculate number of items to display based on container width
  useEffect(() => {
    if (!divRef.current) return;

    const calculateItemsPerPage = () => {
      const containerWidth = divRef.current?.clientWidth || 0;

      if (containerWidth < 640) return 6;
      else if (containerWidth < 768) return 8;
      else if (containerWidth < 1024) return 6;
      else if (containerWidth < 1280) return 8;
      else return 10;
    };

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

  // Search files query function for regular pagination (desktop)
  const searchFiles = async (query: string) => {
    if (!query.trim()) {
      return { items: [], totalpages: 0 };
    }

    const resp = await axiosInstance.get(
      `/files/search?search=${query}&page=${currentPage}&itemsPerPage=${itemsPerPage}`
    );
    return resp.data;
  };

  // Search files query function for infinite query (mobile)
  const getInfiniteSearchResults = async ({ pageParam = 1 }) => {
    if (!searchQuery.trim()) {
      return {
        items: [],
        page: pageParam,
        totalPages: 0,
      };
    }

    const resp = await axiosInstance.get(
      `/files/search?search=${searchQuery}&page=${pageParam}&itemsPerPage=${itemsPerPage}`
    );

    return {
      items: resp.data.items,
      page: pageParam,
      totalPages: resp.data.totalpages,
    };
  };

  // React Query for search results (desktop)
  const {
    data: searchResults,
    isSuccess: pageSuccess,
    isLoading: pageLoading,
    isError: pageError,
    error: pageErrorData,
    refetch: pageRefetch,
  } = useQuery({
    queryKey: ["search-files", searchQuery, currentPage, itemsPerPage],
    queryFn: () => searchFiles(searchQuery),
    enabled: searchQuery.length > 0 && !isMobile,
  });

  // React Query for infinite search results (mobile)
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
    queryKey: ["infinite-search", searchQuery, itemsPerPage],
    queryFn: getInfiniteSearchResults,
    getNextPageParam: (lastPage) => {
      return lastPage.page < lastPage.totalPages ? lastPage.page + 1 : undefined;
    },
    enabled: searchQuery.length > 0 && isMobile,
  });

  // Effect to refetch when mobile/desktop view changes or items per page changes
  useEffect(() => {
    if (searchQuery.length > 0) {
      if (isMobile) {
        infiniteQueryRefetch();
      } else {
        pageRefetch();
      }
    }
  }, [isMobile, itemsPerPage, searchQuery, pageRefetch, infiniteQueryRefetch]);

  // Effect to handle infinite scroll for mobile
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

  // Handle page change for pagination
  const handlePageChange = (newPage: number) => {
    setCurrentPage(newPage);
  };

  // Prepare data for rendering
  const infiniteItems = infiniteData?.pages.flatMap((page) => page.items) || [];
  const isLoading = isMobile ? infiniteLoading : pageLoading;
  const isError = isMobile ? infiniteError : pageError;
  const error = isMobile ? infiniteErrorData : pageErrorData;
  const isSuccess = isMobile ? infiniteSuccess : pageSuccess;

  // Handle when search input is empty
  if (!searchQuery.trim()) {
    return (
      <PageLayout>
        <div ref={divRef} className="p-4 w-full">
          <div className="text-center py-10">
            <p className="text-gray-400">Enter a search term in the search box to find files</p>
          </div>
        </div>
      </PageLayout>
    );
  }

  return (
    <PageLayout>
      <div ref={divRef} className="p-4 w-full">
        <div className="mb-4">
          <h1 className="text-xl font-semibold text-white">
            Search results for: <span className="text-blue-400">{searchQuery}</span>
            {isContentTypeSearch(searchQuery) && (
              <span className="ml-2 text-sm text-gray-400">(searching by file type)</span>
            )}
          </h1>
        </div>

        {/* Loading state */}
        {isLoading && !isFetchingNextPage && (
          <div className="flex justify-center my-10">
            <Spinner loading={true} />
          </div>
        )}

        {/* Error state */}
        {isError && (
          <div className="p-3 bg-red-900 bg-opacity-30 border border-red-800 text-red-300 rounded mb-4">
            {error instanceof Error ? error.message : "An error occurred during search"}
          </div>
        )}

        {/* No results state */}
        {isSuccess &&
          ((isMobile && infiniteItems.length === 0) ||
            (!isMobile && (!searchResults?.items || searchResults.items.length === 0))) && (
            <div className="text-center py-10">
              <p className="text-gray-400">No files found matching your search</p>
              <button
                className="mt-4 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
                onClick={() => navigate("/home")}
              >
                Return to home
              </button>
            </div>
          )}

        {/* Desktop search results with pagination */}
        {!isMobile && isSuccess && searchResults?.items && searchResults.items.length > 0 && (
          <>
            <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-3">
              {searchResults.items.map((file: Iitem, idx: number) => (
                <FileCard key={idx} file={file} />
              ))}
            </div>

            {searchResults.totalpages > 1 && (
              <Pagination
                currentPage={currentPage}
                totalPages={searchResults.totalpages}
                onPageChange={handlePageChange}
                className="flex mt-4 justify-center"
              />
            )}

            <div className="mt-4 text-center text-gray-400">
              Found {searchResults.totalitems} file{searchResults.items.length !== 1 ? "s" : ""}{" "}
              matching your search
            </div>
          </>
        )}

        {/* Mobile search results with infinite scroll */}
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
                <p className="text-gray-400 text-sm">No more results</p>
              )}
            </div>

            <div className="mt-4 text-center text-gray-400">
              Found {infiniteItems.length} file{infiniteItems.length !== 1 ? "s" : ""} matching your
              search
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
              <img src={upperArrow} alt="Scroll to top" className="invert" />
            </button>
          </>
        )}
      </div>
    </PageLayout>
  );
};

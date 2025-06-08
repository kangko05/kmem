import { type Dispatch, type SetStateAction } from "react";
import { X } from "lucide-react";

export const SearchBar = ({
  search,
  setSearch,
}: {
  search: string;
  setSearch: Dispatch<SetStateAction<string>>;
}) => {
  return (
    <div className="relative flex-1 max-w-sm">
      <input
        type="text"
        value={search}
        onChange={(ev) => setSearch(ev.currentTarget.value)}
        placeholder="Search files..."
        className="w-full px-4 py-2 pr-10 border border-gray-300 dark:border-gray-600 
                   rounded-lg bg-white dark:bg-gray-700 
                   text-gray-700 dark:text-gray-300 text-xs sm:text-md
                   placeholder-gray-500 dark:placeholder-gray-400
                   focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent
                   transition-all duration-200"
      />
      {search && (
        <X
          onClick={() => setSearch("")}
          className="absolute right-3 top-1/2 transform -translate-y-1/2 
                     w-4 h-4 text-gray-400 hover:text-gray-600 cursor-pointer
                     transition-colors duration-200"
        />
      )}
    </div>
  );
};

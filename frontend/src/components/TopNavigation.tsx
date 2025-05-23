import { Link, useNavigate } from "react-router";
import { useState, useRef, useEffect, type Dispatch, type SetStateAction } from "react";
import { axiosInstance } from "../utils";
import { LAST_VISITED } from "../constants";
import { Search, X, Menu, LogOut, Upload, Home } from "lucide-react";

interface Inavlink {
  label: string;
  to: string;
  icon: React.ReactNode;
  setIsOpen?: Dispatch<SetStateAction<boolean>>;
}

const Navlink = ({ label, to, icon, setIsOpen }: Inavlink) => {
  return (
    <Link
      to={to}
      className="flex items-center gap-2 text-lg px-4 py-2 rounded-md font-medium hover:bg-gray-700 transition-colors text-white no-underline"
      onClick={() => setIsOpen && setIsOpen(false)}
    >
      {icon}
      <span>{label}</span>
    </Link>
  );
};

const SearchBar = () => {
  const [showSearchBar, setShowSearchBar] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (showSearchBar && inputRef.current) {
      inputRef.current.focus();
    }
  }, [showSearchBar]);

  return (
    <div className="relative flex items-center">
      {!showSearchBar ? (
        <button
          onClick={() => setShowSearchBar(true)}
          className="p-2 rounded-full hover:bg-gray-700 transition-colors"
          aria-label="Show search"
        >
          <Search size={20} className="m-auto" />
        </button>
      ) : (
        <form className="flex items-center">
          <div className="relative flex items-center">
            <input
              ref={inputRef}
              type="text"
              value={searchQuery}
              onChange={(ev) => setSearchQuery(ev.target.value)}
              className="w-50 md:w-60 h-9 pl-10 pr-8 bg-gray-700 text-white text-sm border border-gray-600 rounded-full focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="search..."
            />
            <button
              type="button"
              onClick={() => {
                setShowSearchBar(false);
                setSearchQuery("");
              }}
              className="absolute left-0 p-1 rounded-full hover:bg-gray-600 transition-colors"
              aria-label="Clear search"
            >
              <X size={16} className="text-gray-400 m-auto" />
            </button>
            <Link
              className="absolute right-0 p-1 rounded-full hover:bg-gray-600 transition-colors"
              to={`/search?q=${searchQuery}`}
            >
              <Search size={16} className="text-gray-400 m-auto" />
            </Link>
          </div>
        </form>
      )}
    </div>
  );
};

export const TopNavigation = () => {
  const [isOpen, setIsOpen] = useState(false);
  const navigate = useNavigate();

  const handleLogout = async () => {
    try {
      await axiosInstance.get("/auth/logout");
      const lastVisited = localStorage.getItem(LAST_VISITED);
      if (lastVisited && lastVisited != "/") {
        localStorage.removeItem(LAST_VISITED);
      }
      navigate("/");
    } catch (error) {
      console.error("Logout failed:", error);
    }
  };

  const navItems: Inavlink[] = [
    { label: "Home", to: "/home", icon: <Home size={18} /> },
    { label: "Upload", to: "/upload", icon: <Upload size={18} /> },
  ];

  return (
    <nav className="fixed top-0 w-full max-w-7xl h-16 bg-gray-800 text-white shadow-md z-10">
      <div className="w-full h-full mx-auto px-4">
        <div className="flex items-center justify-between h-full">
          <div className="flex items-center">
            {/* Desktop Navigation */}
            <div className="hidden md:flex space-x-1">
              {navItems.map((item, index) => (
                <Navlink key={index} {...item} />
              ))}
            </div>
          </div>

          {/* Search & Logout - Desktop */}
          <div className="hidden md:flex items-center space-x-4">
            <SearchBar />

            <button
              onClick={handleLogout}
              className="flex items-center gap-2 px-4 py-2 rounded-md text-sm font-medium bg-gray-700 hover:bg-gray-600 transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <LogOut size={16} />
              <span>Logout</span>
            </button>
          </div>

          {/* Mobile Controls */}
          <div className="flex w-full md:hidden justify-between items-center space-x-2">
            <div className="flex w-full">
              <button
                onClick={() => setIsOpen(!isOpen)}
                className="p-2 rounded-md hover:bg-gray-700 transition-colors"
                aria-label="Toggle menu"
              >
                {isOpen ? <X size={20} /> : <Menu size={20} />}
              </button>

              <SearchBar />
            </div>

            <button
              onClick={handleLogout}
              className="flex flex-center p-2 rounded-md bg-gray-700 hover:bg-gray-600 transition-colors"
              aria-label="Logout"
            >
              <LogOut size={16} />
            </button>
          </div>
        </div>
      </div>

      {/* Mobile Navigation Menu */}
      {isOpen && (
        <div className="absolute top-16 left-0 right-0 bg-gray-800 border-t border-gray-700 md:hidden shadow-lg py-2 z-20">
          {navItems.map((item, index) => (
            <div key={index} className="px-4 py-1">
              <Navlink {...item} setIsOpen={setIsOpen} />
            </div>
          ))}
        </div>
      )}
    </nav>
  );
};

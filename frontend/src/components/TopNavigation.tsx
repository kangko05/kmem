import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { LogOut, Menu, X } from "lucide-react";
import { axiosInstance } from "../utils";

const NavLink = ({ to, label, isMobile }: { to: string; label: string; isMobile?: boolean }) => {
  const mobileCss =
    "block px-4 py-3 text-gray-600 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-all duration-200 font-medium";
  const desktopCss =
    "block px-4 py-3 text-gray-600 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-all duration-200 font-medium";

  return (
    <Link to={to} className={isMobile ? mobileCss : desktopCss}>
      {label}
    </Link>
  );
};

const LogoutBtn = ({ isMobile }: { isMobile?: boolean }) => {
  const navigate = useNavigate();

  const mobileCss =
    "flex items-center gap-2 w-full px-4 py-3 text-gray-600 dark:text-gray-300 hover:text-red-600 dark:hover:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg transition-all duration-200 font-medium text-left";
  const desktopCss =
    "hidden md:flex items-center gap-2 px-4 py-2 text-gray-600 dark:text-gray-300 hover:text-red-600 dark:hover:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg transition-all duration-200 font-medium";

  const handleClick = async () => {
    await axiosInstance.get("/auth/logout");
    navigate("/login");
  };

  return (
    <button className={isMobile ? mobileCss : desktopCss} onClick={handleClick}>
      <LogOut className="w-4 h-4" />
      Logout
    </button>
  );
};

export const TopNavigation = () => {
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  return (
    <>
      <nav
        className="fixed top-0 left-0 w-full h-16 bg-white dark:bg-gray-800 
                     shadow-md z-50 border-b border-gray-200 dark:border-gray-700"
      >
        <div className="flex justify-between items-center h-full px-4 max-w-7xl mx-auto">
          <div className="flex items-center gap-6">
            {/* for desktop */}
            <div className="hidden md:flex items-center gap-1">
              <NavLink to="/home" label="Home" />
              <NavLink to="/gallery" label="Gallery" />
              <NavLink to="/upload" label="Upload" />
            </div>
          </div>

          {/* right section */}
          <div className="flex items-center gap-3">
            {/* logout (desktop) */}
            <LogoutBtn />

            {/* for mobile */}
            <button
              onClick={() => setIsMenuOpen(!isMenuOpen)}
              className="md:hidden p-2 text-gray-600 dark:text-gray-300 
                       hover:text-gray-800 dark:hover:text-white
                       hover:bg-gray-100 dark:hover:bg-gray-700
                       rounded-lg transition-all duration-200"
            >
              {isMenuOpen ? <X className="w-4 h-4" /> : <Menu className="w-4 h-4" />}
            </button>
          </div>
        </div>

        {/* mobile menu */}
        {isMenuOpen && (
          <div
            className="md:hidden absolute top-16 left-0 w-full 
                         bg-white dark:bg-gray-800 shadow-lg 
                         border-t border-gray-200 dark:border-gray-700"
          >
            <div className="px-4 py-2 space-y-1">
              <NavLink to="/home" label="Home" isMobile />
              <NavLink to="/gallery" label="Gallery" isMobile />
              <NavLink to="/upload" label="Upload" isMobile />

              <hr className="my-2 border-gray-200 dark:border-gray-700" />

              <LogoutBtn isMobile />
            </div>
          </div>
        )}
      </nav>

      <div className="h-16"></div>
    </>
  );
};

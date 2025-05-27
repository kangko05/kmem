import { useEffect } from "react";
import { axiosInstance } from "../utils";
import { useNavigate, useLocation } from "react-router-dom";
import { LAST_VISTIED } from "../constants";

export const useAuth = () => {
  const loc = useLocation();
  const navigate = useNavigate();

  useEffect(() => {
    (async () => {
      try {
        await axiosInstance.get("/auth/me");

        let lastPage = localStorage.getItem(LAST_VISTIED);

        if (!lastPage || lastPage == "/login") lastPage = "/home";

        navigate(lastPage); // TODO: this need to navigate to last visited page
      } catch {
        navigate("/login");
      }
    })();

    return () => {
      localStorage.setItem(LAST_VISTIED, loc.pathname);
    };
  }, []);
};

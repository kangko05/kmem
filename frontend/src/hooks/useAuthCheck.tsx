import { useEffect } from "react";
import { axiosInstance } from "../utils/AxiosIntstance";
import { useNavigate, useLocation } from "react-router";
import { LAST_VISITED } from "../constants";

export const useAuthCheck = () => {
  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    (async () => {
      try {
        const resp = await axiosInstance.get("/auth/me");

        if (resp.status == 200 && location.pathname == "/") {
          const prevPage = localStorage.getItem(LAST_VISITED) || "/home";
          navigate(prevPage == "/" ? "/home" : prevPage);
        }
      } catch {
        navigate("/");
      }
    })();

    return () => {
      if (location.pathname != "/") {
        localStorage.setItem(LAST_VISITED, location.pathname);
      }
    };
  }, []);
};

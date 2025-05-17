import { axiosInstance } from "../utils/AxiosIntstance";
import { useAuthCheck } from "../hooks";
import { LAST_VISITED } from "../constants";
import { Link, useNavigate } from "react-router";
import { UploadBox } from "../components";

export const HomePage = () => {
  const navigate = useNavigate();

  useAuthCheck();

  const test = async () => {
    try {
      const resp = await axiosInstance.get("/auth/me", {
        withCredentials: true,
      });
      console.log(resp);
    } catch (err) {
      console.error(err);
    }
  };

  const handleLogout = async () => {
    const resp = await axiosInstance.get("/auth/logout");

    console.log(resp.data);

    const lastVisited = localStorage.getItem(LAST_VISITED);
    if (lastVisited && lastVisited != "/") localStorage.removeItem(LAST_VISITED);

    navigate("/");
  };

  return (
    <>
      <h1>Home</h1>

      <button className="btn text-sm sm:text-base md:text-lg mt-3" onClick={test}>
        click
      </button>

      <button className="btn text-sm sm:text-base md:text-lg mt-3" onClick={handleLogout}>
        logout
      </button>

      <Link className="btn text-amber-300 text-sm sm:text-base md:text-lg mt-3" to="/more">
        More
      </Link>

      <UploadBox />
    </>
  );
};

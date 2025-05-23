import { Link } from "react-router";
import { TopNavigation } from "../components";
import { useAuthCheck } from "../hooks";

export const MorePage = () => {
  useAuthCheck();

  return (
    <>
      <TopNavigation />
      <h1>More!</h1>

      <Link className="btn text-neutral-100 text-sm sm:text-base md:text-lg mt-3" to="/home">
        Home
      </Link>
    </>
  );
};

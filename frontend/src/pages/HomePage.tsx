import { Link } from "react-router-dom";
import { PageLayout } from "../components";
import { useAuth } from "../hooks/useAuth";

export const HomePage = () => {
  useAuth();

  return (
    <PageLayout>
      <h1>Home!</h1>
      <Link to="/login" className="btn">
        Back to Login{" "}
      </Link>
    </PageLayout>
  );
};

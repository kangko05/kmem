import { useState } from "react";
import { PageLayout, TextInput } from "../components";
import { axiosInstance } from "../utils";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../hooks/useAuth";

export const LoginPage = () => {
  const navigate = useNavigate();

  useAuth();

  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [isError, setIsError] = useState(false);

  const handleClick = async () => {
    if (username.length < 4) {
      setIsError(true);
      return;
    }

    if (password.length < 8) {
      setIsError(true);
      return;
    }

    setUsername("");
    setPassword("");

    try {
      setIsLoading(true);
      const resp = await axiosInstance.post("/auth/login", {
        username: username,
        password: password,
      });

      if (resp.status == 200) navigate("/home");
    } catch (err) {
      console.error(err);
      setIsError(true);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <PageLayout>
      <div className="bg-white dark:bg-gray-800 p-8 rounded-lg shadow-lg w-[90%] sm:w-full max-w-md">
        <h2 className="text-2xl font-bold text-center mb-6 text-gray-800 dark:text-white">
          Sign In
        </h2>

        <div className="space-y-4">
          <TextInput
            placeholder="username"
            value={username}
            onChange={(ev) => {
              setIsError(false);
              setUsername(ev.currentTarget.value);
            }}
          />

          <TextInput
            placeholder="password"
            type="password"
            value={password}
            onChange={(ev) => {
              setIsError(false);
              setPassword(ev.currentTarget.value);
            }}
          />

          {isError && <p className="text-red-400">invalid username or password</p>}

          <button
            className="w-full bg-blue-500 hover:bg-blue-600 text-white font-semibold 
                   py-3 px-4 rounded-lg transition-colors duration-200
                   focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2
                   active:bg-blue-700"
            onClick={handleClick}
          >
            {isLoading ? "signing in..." : "sign in"}
          </button>
        </div>
      </div>
    </PageLayout>
  );
};

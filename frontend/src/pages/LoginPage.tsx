import { type FormEvent, useState } from "react";
import { LAST_VISITED, SERVER } from "../constants";
import { useMutation } from "@tanstack/react-query";
import { useNavigate } from "react-router";
import { axiosInstance } from "../utils/AxiosIntstance";
import { useAuthCheck } from "../hooks";

import axios from "axios";

export const LoginPage = () => {
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const navigate = useNavigate();

  useAuthCheck();

  const login = async (data: { username: string; password: string }) => {
    return await axiosInstance.post(`${SERVER}/auth/login`, data, { withCredentials: true });
  };

  const { mutate, isError, isPending } = useMutation({
    mutationFn: login,
    onError: (error) => {
      // setErrorMsg("invalid username or password");
      if (axios.isAxiosError(error)) {
        setErrorMsg(error.response?.data || "Invalid username or password");
      } else {
        setErrorMsg(error.message);
      }
    },
    onSuccess: () => {
      const prevPage = localStorage.getItem(LAST_VISITED) || "/home";
      navigate(prevPage == "/" ? "/home" : prevPage);
    },
  });

  const handleSubmit = (ev: FormEvent<HTMLFormElement>) => {
    ev.preventDefault();
    setErrorMsg(null);

    const formData = new FormData(ev.currentTarget);
    const username = formData.get("username")?.toString() || "";
    const password = formData.get("password")?.toString() || "";

    if (username.length < 4) {
      setErrorMsg("Username must be at least 4 characters");
      return;
    }

    if (password.length < 8) {
      setErrorMsg("Password must be at least 8 characters");
      return;
    }

    mutate({ username, password });
  };

  return (
    <div className="page-container flex-col w-dvw h-dvh">
      <form className="w-5/6 max-w-96 text-sm sm:text-base md:text-lg" onSubmit={handleSubmit}>
        <div className="flex flex-col gap-3">
          <label htmlFor="username">ID</label>
          <input
            className="text-sm sm:text-base md:text-lg"
            type="text"
            minLength={4}
            autoComplete="username"
            placeholder="enter your id"
            id="username"
            name="username"
          />
        </div>
        <div className="flex flex-col gap-3">
          <label htmlFor="password">Password</label>
          <input
            className="text-sm sm:text-base md:text-lg"
            type="password"
            minLength={8}
            autoComplete="off"
            placeholder="enter your password"
            id="password"
            name="password"
          />
        </div>
        <p className="block w-5/6 h-4 max-w-96 text-red-400 mb-1 mt-2">
          {(isError || errorMsg) && errorMsg}
        </p>
        <button type="submit" className="btn text-sm sm:text-base md:text-lg" disabled={isPending}>
          {isPending ? "Signing In..." : "Sign In"}
        </button>
      </form>
    </div>
  );
};

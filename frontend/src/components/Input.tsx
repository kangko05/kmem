import { type ChangeEvent } from "react";

export const TextInput = ({
  placeholder,
  type,
  value,
  onChange,
}: {
  placeholder: string;
  type?: "text" | "password";
  value: string;
  onChange: (ev: ChangeEvent<HTMLInputElement>) => void;
}) => {
  return (
    <input
      className="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 
                     rounded-lg bg-white dark:bg-gray-700 
                     text-gray-800 dark:text-white
                     placeholder-gray-500 dark:placeholder-gray-400
                     focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent
                     transition-all duration-200"
      type={type ? type : "text"}
      placeholder={placeholder}
      value={value}
      onChange={onChange}
      minLength={type == "text" ? 4 : 8}
    />
  );
};

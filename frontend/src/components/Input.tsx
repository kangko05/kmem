import { type ChangeEvent } from "react";
import { Upload } from "lucide-react";

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

export const FileInput = ({
  onChange,
}: {
  onChange: (ev: ChangeEvent<HTMLInputElement>) => void;
}) => {
  return (
    <div className="w-[90%] sm:w-full max-w-md mx-auto">
      <label className="flex flex-col items-center justify-center w-full h-32 sm:h-64 border-2 border-gray-300 border-dashed rounded-lg cursor-pointer bg-gray-50 dark:bg-gray-700 hover:bg-gray-100 dark:hover:bg-gray-600 dark:border-gray-600 transition-colors duration-200">
        <div className="flex flex-col items-center justify-center pt-5 pb-6">
          <Upload className="w-10 h-10 mb-3 text-gray-400" />
          <p className="mb-2 text-sm text-gray-500 dark:text-gray-400">
            <span className="font-semibold">Click to upload</span> or drag and drop
          </p>
          <p className="text-xs text-gray-500 dark:text-gray-400">PNG, JPG, MP4, AVI</p>
        </div>
        <input
          type="file"
          className="hidden"
          accept="image/*,video/*"
          multiple
          onChange={onChange}
        />
      </label>
    </div>
  );
};

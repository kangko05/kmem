import { type ChangeEvent, type ReactNode } from "react";

export const DropDown = ({
  name,
  onChange,
  children,
}: {
  name: string;
  onChange: (ev: ChangeEvent<HTMLSelectElement>) => void;
  children: ReactNode;
}) => {
  return (
    <select
      name={name}
      onChange={(ev) => onChange(ev)}
      className="w-fit px-3 py-2 border border-gray-300 rounded-lg bg-white dark:bg-gray-700 dark:border-gray-600 text-gray-700 dark:text-gray-300 text-xs sm:text-md"
    >
      {children}
    </select>
  );
};

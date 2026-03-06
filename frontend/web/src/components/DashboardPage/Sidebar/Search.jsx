import { FiSearch, FiCommand } from "react-icons/fi";
import { CommandMenu } from "./CommandMenu";
import { useState } from "react";

export const Search = () => {
  const [open, setOpen] = useState(false);

  return (
    <>
      <div className="bg-slate-700 mb-4 relative rounded flex items-center px-2 py-1.5 text-sm">
        <FiSearch className="mr-2" />
        <input
          onFocus={(e) => {
            e.target.blur();
            setOpen(true);
          }}
          type="text"
          placeholder="search"
          className="w-full bg-transparent placeholder:text-slate-400 focus:outline-none"
        />
        <span className="flex gap-0.5 p-1 text-xs items-center rounded absolute right-1.5 top-1/2 -translate-y-1/2">
          <FiCommand />K
        </span>
      </div>

      <CommandMenu open={open} setOpen={setOpen} />
    </>
  );
};

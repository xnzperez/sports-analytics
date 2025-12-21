import React from "react";
import { cn } from "../../lib/utils"; // Usamos la utilidad que creamos antes

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
}

// Usamos forwardRef para que funcione con react-hook-form
export const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ className, label, error, ...props }, ref) => {
    return (
      <div className="w-full space-y-2">
        {label && (
          <label className="text-sm font-medium text-slate-300 ml-1">
            {label}
          </label>
        )}
        <input
          ref={ref}
          className={cn(
            "flex h-10 w-full rounded-md border border-slate-700 bg-slate-800 px-3 py-2 text-sm text-slate-100 placeholder:text-slate-500",
            "focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-transparent transition-all",
            "disabled:cursor-not-allowed disabled:opacity-50",
            error && "border-red-500 focus:ring-red-500", // Borde rojo si hay error
            className
          )}
          {...props}
        />
        {error && (
          <p className="text-xs text-red-400 font-medium ml-1 animate-pulse">
            {error}
          </p>
        )}
      </div>
    );
  }
);

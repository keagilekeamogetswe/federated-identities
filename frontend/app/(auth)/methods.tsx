"use client";
import { useRef, useState } from "react";
import { AuthMethods } from "./method.config";
import { APIInterfaceModule } from "./provider/api.interface";

export function MethodsUIComponent() {
  const [loading_request, setLoading] = useState<boolean>(false);
  const current_loading_method = useRef<string | null>(null);
  const [error_message, setError] = useState<string | null>(null);

  return (
    <>
      {error_message && (
        <div className="w-full mb-4 p-3 rounded-xl bg-red-950/50 border border-red-800/60 text-red-300 text-sm font-medium text-left">
          {error_message}
        </div>
      )}

      <div className="space-y-3 w-full">
        {AuthMethods.map((method) => {
          const isLoading =
            loading_request && current_loading_method.current === method.name;

          return (
            <a
              key={method.name}
              href="#_" // Fallback href, prevented by onClick
              onClick={async (e) => {
                e.preventDefault(); // Stop native navigation so async fetch runs
                if (loading_request) return;

                // Clear any previous errors on a new attempt
                setError(null);
                current_loading_method.current = method.name;
                setLoading(true);

                try {
                  // Calls backend returning { url: "..." }
                  const redirect_url = await APIInterfaceModule.requestID(
                    method.name,
                  );
                  if (redirect_url) {
                    window.location.href = redirect_url;
                  }
                } catch (error: any) {
                  console.error("Failed to fetch auth URL:", error);
                  setError(
                    error?.message
                      ? "An unexpected error occured. Please try again."
                      : `Failed to connect to ${method.name}. Please try again.`,
                  );
                } finally {
                  current_loading_method.current = null;
                  setLoading(false);
                }
              }}
              className="w-full flex items-center gap-x-3 px-4 py-3 rounded-xl bg-neutral-800/60 hover:bg-neutral-800 border border-neutral-700/60 hover:border-neutral-600 text-neutral-200 font-medium text-sm transition-all duration-200 shadow-sm group cursor-pointer"
            >
              <div className="w-8 h-8 rounded-full bg-white/90 p-0.5 flex items-center justify-center shrink-0 shadow-sm">
                {isLoading ? (
                  <span className="animate-spin w-4 h-4 border-2 border-neutral-900 border-t-transparent rounded-full" />
                ) : (
                  <img
                    src={method.icon}
                    alt={method.name}
                    className="w-full h-full object-contain"
                  />
                )}
              </div>
              <span className="flex-1 text-left tracking-wide capitalize">
                Continue with {method.name}
              </span>
              <svg
                className="w-4 h-4 text-neutral-500 group-hover:text-neutral-300 transition-transform group-hover:translate-x-0.5"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M9 5l7 7-7 7"
                />
              </svg>
            </a>
          );
        })}
      </div>
    </>
  );
}

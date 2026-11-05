"use client";
import { useRef, useState } from "react";
import { AuthMethods } from "./method.config";
import { APIInterfaceModule } from "./provider/api.interface";

export function MethodsUIComponent() {
  const [loading_request, setLoading] = useState<boolean>(false);
  const current_loading_method = useRef<string | null>(null);

  return (
    <>
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
              } catch (error) {
                console.error("Failed to fetch auth URL:", error);
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
            <span className="flex-1 text-left tracking-wide">
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
                strokeLinejoin="round"
                strokeWidth={2}
                d="M9 5l7 7-7 7"
              />
            </svg>
          </a>
        );
      })}
    </>
  );
}

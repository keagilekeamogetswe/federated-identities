"use client";
import { AuthMethods } from "@/app/(auth)/method.config";
import { APIInterfaceModule } from "../../api.interface";
import { useEffect, useRef, useState } from "react";
import { useParams } from "next/navigation";

export default function IdentityProvider() {
  const params = useParams();
  const provider_name = typeof params.name === "string" ? params.name : "";

  const [loading_confirmation, setLoadingConfirmation] =
    useState<boolean>(true);
  const [error_message, setErrorMessage] = useState<string | null>(null);
  const [userData, setUserData] = useState<any>(null);

  // Guard ref to strictly prevent double execution in React Strict Mode / re-renders
  const hasFetched = useRef<boolean>(false);

  useEffect(() => {
    if (!provider_name || hasFetched.current) return;
    hasFetched.current = true; // Lock execution immediately

    async function makeBackendCall() {
      try {
        const searchString = window.location.search; // Gets ?code=...&state=...
        const req = await APIInterfaceModule.confirmID(
          provider_name,
          searchString,
        );
        console.log("data: ", req);

        if (!req) {
          setErrorMessage("Authentication confirmation failed.");
          setLoadingConfirmation(false);
          return;
        }

        setUserData(req);
      } catch (error: any) {
        console.error("Confirmation error:", error);
        setErrorMessage(
          error?.message || "An unexpected network error occurred.",
        );
      } finally {
        setLoadingConfirmation(false);
      }
    }

    makeBackendCall();
  }, [provider_name]);

  const current_method_config = AuthMethods.find(
    (item) => item.name.toLowerCase() === provider_name.toLowerCase(),
  );

  if (!current_method_config) {
    const formattedName = provider_name
      ? provider_name.charAt(0).toUpperCase() + provider_name.slice(1)
      : "Unknown";
    return (
      <main className="min-h-screen bg-black text-white flex flex-col items-center justify-center p-4 selection:bg-neutral-800">
        <div className="w-full max-w-md bg-neutral-900 border border-neutral-800 rounded-2xl p-8 shadow-2xl text-center space-y-6">
          <div className="space-y-2">
            <h1 className="text-xl font-semibold tracking-tight text-red-400">
              Unsupported Identity Provider
            </h1>
            <p className="text-sm text-neutral-400">
              <span className="capitalize text-white font-medium">
                {formattedName}
              </span>{" "}
              is not supported by this gateway.
            </p>
          </div>
          <a
            href="/"
            className="inline-flex w-full items-center justify-center rounded-xl bg-neutral-800 px-4 py-3 text-sm font-medium text-white hover:bg-neutral-700 transition border border-neutral-700/50"
          >
            Return Home
          </a>
        </div>
        <footer className="mt-8 text-xs tracking-widest text-neutral-600 uppercase font-mono">
          Secure Identity Gateway
        </footer>
      </main>
    );
  }

  // Fallback parsers for display name & email formatting safely
  const displayName =
    typeof userData?.name === "string"
      ? userData.name
      : typeof userData?.email === "object"
        ? userData?.email?.name
        : "Verified User";

  const emailAddress =
    typeof userData?.email === "string"
      ? userData.email
      : typeof userData?.email === "object"
        ? userData?.email?.address || ""
        : "";

  return (
    <main className="min-h-screen bg-black text-white flex flex-col items-center justify-center p-4 selection:bg-neutral-800">
      <div className="w-full max-w-md bg-neutral-900 border border-neutral-800 rounded-2xl p-8 shadow-2xl text-center space-y-6">
        <div className="space-y-2">
          <h1 className="text-xl font-semibold tracking-tight">
            {loading_confirmation
              ? "Authenticating Digital ID"
              : error_message
                ? "Authentication Failed"
                : "Authentication Complete"}
          </h1>
          <p className="text-sm text-neutral-400">
            {loading_confirmation ? (
              <>
                Communicating with{" "}
                <span className="capitalize text-white font-medium">
                  {current_method_config.name}
                </span>
                . <br />
                This might take a moment.
              </>
            ) : error_message ? (
              <span className="text-red-400 text-xs break-words block mt-1">
                {error_message}
              </span>
            ) : (
              <>
                Your identity has been successfully verified via{" "}
                <span className="capitalize text-white font-medium">
                  {current_method_config.name}
                </span>
                .
              </>
            )}
          </p>
        </div>

        {/* State feedback & profile display */}
        <div className="flex flex-col items-center justify-center py-2 space-y-4">
          {loading_confirmation ? (
            <div className="w-6 h-6 border-2 border-neutral-700 border-t-white rounded-full animate-spin"></div>
          ) : error_message ? (
            <a
              href="/"
              className="inline-flex w-full items-center justify-center rounded-xl bg-red-950/40 hover:bg-red-900/40 border border-red-800/60 px-4 py-2.5 text-xs font-medium text-red-200 transition"
            >
              Back to Sign In
            </a>
          ) : (
            <div className="w-full space-y-5">
              {/* Status Badge */}
              <div className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 text-xs font-medium">
                <span>Verified Successfully</span>
                <svg
                  className="w-3.5 h-3.5"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M5 13l4 4L19 7"
                  />
                </svg>
              </div>

              {/* User Profile Card - Horizontal Layout */}
              {userData && (
                <div className="flex items-center gap-4 bg-neutral-950/50 border border-neutral-800/80 rounded-xl p-4 shadow-inner w-full text-left">
                  {userData.picture ? (
                    <img
                      src={userData.picture}
                      alt={displayName}
                      referrerPolicy="no-referrer"
                      className="w-14 h-14 rounded-full object-cover border-2 border-neutral-700/60 shadow-md ring-2 ring-emerald-500/20 shrink-0"
                    />
                  ) : (
                    <div className="w-14 h-14 rounded-full bg-neutral-800 border-2 border-neutral-700 flex items-center justify-center text-lg font-semibold text-neutral-300 shrink-0">
                      {displayName.charAt(0).toUpperCase()}
                    </div>
                  )}

                  <div className="min-w-0 flex-1">
                    <h2 className="text-base font-semibold text-white truncate">
                      {displayName}
                    </h2>
                    {emailAddress && (
                      <p className="text-xs text-neutral-400 font-mono truncate mt-0.5">
                        {emailAddress}
                      </p>
                    )}
                  </div>
                </div>
              )}
            </div>
          )}
        </div>
      </div>

      <footer className="mt-8 text-xs tracking-widest text-neutral-600 uppercase font-mono">
        Secure Identity Gateway
      </footer>
    </main>
  );
}

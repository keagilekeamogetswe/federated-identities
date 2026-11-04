import { AuthMethods } from "./method.config.ts";

export function MethodsUIComponent() {
  return AuthMethods.map((method) => {
    return (
      <a
        key={method.name}
        href={method.endpoint}
        className="w-full flex items-center gap-x-3 px-4 py-3 rounded-xl bg-neutral-800/60 hover:bg-neutral-800 border border-neutral-700/60 hover:border-neutral-600 text-neutral-200 font-medium text-sm transition-all duration-200 shadow-sm group"
      >
        <div className="w-8 h-8 rounded-full bg-white/90 p-1.5 flex items-center justify-center shrink-0 shadow-sm">
          <img
            src={method.icon}
            alt={method.name}
            className="w-full h-full object-contain"
          />
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
            strokeWidth={2}
            d="M9 5l7 7-7 7"
          />
        </svg>
      </a>
    );
  });
}

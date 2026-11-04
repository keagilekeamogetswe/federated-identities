import { MethodsUIComponent } from "./(auth)/methods";

export default function Home() {
  return (
    <main className="min-h-screen bg-neutral-950 text-neutral-100 flex flex-col items-center justify-center p-6 selection:bg-neutral-800">
      <div className="w-full max-w-md bg-neutral-900/80 border border-neutral-800/80 rounded-2xl p-8 shadow-2xl backdrop-blur-xl flex flex-col items-center text-center">
        <h1 className="text-2xl font-semibold tracking-tight text-white mb-3">
          Welcome to my Experiment
        </h1>

        <p className="text-sm text-neutral-400 leading-relaxed mb-8">
          Usernames and passwords are not allowed here. We only trust our
          identity providers to determine if you are authorised to go beyond.
        </p>

        <div className="w-full flex flex-col gap-y-3">
          <MethodsUIComponent />
        </div>
      </div>

      <footer className="mt-8 text-xs text-neutral-600 tracking-wider uppercase font-medium">
        Secure Identity Gateway
      </footer>
    </main>
  );
}

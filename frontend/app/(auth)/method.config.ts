interface Method {
  name: string;
  icon: string;
  endpoint: string;
}

export const AuthMethods: Method[] = [
  {
    name: "Google",
    icon: "/assets/google.svg",
    endpoint: "api/provider/google",
  },
  {
    name: "Github",
    icon: "/assets/github.svg",
    endpoint: "api/provider/github",
  },
];

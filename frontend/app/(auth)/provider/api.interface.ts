export const APIInterfaceModule = (() => {
  function buildURL(provider: string, action: string): string {
    const provider_url_base = `/api/provider/`;
    return `${provider_url_base}/${provider}/${action}`;
  }
  return {
    async requestID(provider: string) {
      const url = buildURL(provider, "request");
      const req = await fetch(url, { method: "GET", credentials: "include" });
      if (!req.ok) return null;
      const data = await req.json();
      return data.url;
    },
    async confirmID(provider: string, params: string) {
      const url = buildURL(provider, "confirm");
      const req = await fetch(url + params, {
        method: "GET",
        credentials: "include",
      });
      console.log(req);
      if (!req.ok) return null;
      const data = await req.json();
      return data;
    },
  };
})();

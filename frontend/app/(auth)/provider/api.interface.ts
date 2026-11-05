export const APIInterfaceModule = (() => {
  function buildURL(provider: string, action: string): string {
    const provider_url_base = `/api/provider/`;
    return `${provider_url_base}/${provider}/${action}`;
  }
  return {
    async requestID(provider: string) {
      const url = buildURL(provider, "request");
      const req = await fetch(url, { method: "GET" });
      if (!req.ok) return null;
      const data = await req.json();
      return data.url;
    },
  };
})();

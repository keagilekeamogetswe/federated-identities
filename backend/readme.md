## Architectural Reference: Google OIDC vs. GitHub OAuth 2.0 Integration

This guide documents the core architectural, protocol, and implementation differences between integrating **Google OpenID Connect (OIDC)** and **GitHub OAuth 2.0** authentication flows in Go backend services.

---

### Core Protocol Comparison

| Feature / Dimension        | Google OIDC Authentication                                                                                                                        | GitHub OAuth 2.0 Authentication                                                                                                    |
| :------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------ | :--------------------------------------------------------------------------------------------------------------------------------- |
| **Protocol Specification** | OpenID Connect (OIDC) built on top of OAuth 2.0                                                                                                   | Standard OAuth 2.0 Authorization Code Grant                                                                                        |
| **Token Payload**          | Returns an **`id_token`** (signed JWT containing identity claims) alongside the `access_token`.                                                   | Returns strictly an **`access_token`** (opaque or bearer token). No ID token is issued.                                            |
| **Identity Resolution**    | **Stateless local verification:** The `id_token` can be cryptographically verified offline/locally using public keys (JWKS) via Google's library. | **Stateful remote API call:** Requires an explicit outbound REST API call to `https://api.github.com/user` using the access token. |
| **Email Privacy & Scopes** | Bundled directly into standard scopes (`openid email profile`) and validated within the `id_token`.                                               | Requires explicit scope (`user:email`). Private emails require a secondary fallback API query (`/user/emails`).                    |
| **Token Format Request**   | Standard URL-encoded form POST; responds with JSON automatically.                                                                                 | Requires explicit header `Accept: application/json`, otherwise responds with a URL-encoded query string.                           |

---

### Deep Dive into Key Differences

#### 1. Identity Standard & Token Verification

- **Google (OIDC):** Because Google implements OIDC, it returns a cryptographically signed JSON Web Token (`id_token`). Your backend validates the signature, audience (Client ID), issuer, and expiration time completely in-memory using Google's public key infrastructure (`google.golang.org/api/idtoken`).
- **GitHub (OAuth 2.0):** GitHub is not a native OpenID provider for web apps. After exchanging your authorization code for an `access_token`, your server must act as an HTTP client, attach the token as a Bearer header, and query GitHub's REST API (`GET https://api.github.com/user`) to discover who just logged in.

#### 2. Handling Private Emails

- **Google:** Verified user emails are inherently guaranteed and surfaced through the authenticated OIDC claims.
- **GitHub:** Many developers configure their GitHub accounts to keep their email addresses private (`@users.noreply.github.com`). If a standard `/user` request returns a blank or private email string, your service must handle a secondary back-channel request to `https://api.github.com/user/emails`, parsing the array to filter for `Primary: true` and `Verified: true`.

#### 3. Content Negotiation on Token Exchange

- **Google:** The token endpoint `/token` natively expects a standard `POST` request with form parameters and responds with a strict JSON payload.
- **GitHub:** The endpoint `https://github.com/login/oauth/access_token` requires you to explicitly set the HTTP header `Accept: application/json`. If omitted, GitHub responds with a URL-encoded string (`access_token=...&scope=...&token_type=bearer`), which will fail standard JSON decoders.

---

### Implementation Blueprint Summary

- **Google Flow:**
  `Redirect to Google` $\rightarrow$ `Receive Auth Code` $\rightarrow$ `POST /token (Get ID Token + Access Token)` $\rightarrow$ `Verify ID Token Locally (JWT)` $\rightarrow$ `Extract Claims`.
- **GitHub Flow:**
  `Redirect to GitHub` $\rightarrow$ `Receive Auth Code` $\rightarrow$ `POST /access_token (Accept: application/json)` $\rightarrow$ `GET /user (Bearer Token)` $\rightarrow$ `Optional: GET /user/emails (If Primary Email is Private)`.## Architectural Reference: Google OIDC vs. GitHub OAuth 2.0 Integration

This guide documents the core architectural, protocol, and implementation differences between integrating **Google OpenID Connect (OIDC)** and **GitHub OAuth 2.0** authentication flows in Go backend services.

---

### Core Protocol Comparison

| Feature / Dimension        | Google OIDC Authentication                                                                                                                        | GitHub OAuth 2.0 Authentication                                                                                                    |
| :------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------ | :--------------------------------------------------------------------------------------------------------------------------------- |
| **Protocol Specification** | OpenID Connect (OIDC) built on top of OAuth 2.0                                                                                                   | Standard OAuth 2.0 Authorization Code Grant                                                                                        |
| **Token Payload**          | Returns an **`id_token`** (signed JWT containing identity claims) alongside the `access_token`.                                                   | Returns strictly an **`access_token`** (opaque or bearer token). No ID token is issued.                                            |
| **Identity Resolution**    | **Stateless local verification:** The `id_token` can be cryptographically verified offline/locally using public keys (JWKS) via Google's library. | **Stateful remote API call:** Requires an explicit outbound REST API call to `https://api.github.com/user` using the access token. |
| **Email Privacy & Scopes** | Bundled directly into standard scopes (`openid email profile`) and validated within the `id_token`.                                               | Requires explicit scope (`user:email`). Private emails require a secondary fallback API query (`/user/emails`).                    |
| **Token Format Request**   | Standard URL-encoded form POST; responds with JSON automatically.                                                                                 | Requires explicit header `Accept: application/json`, otherwise responds with a URL-encoded query string.                           |

---

### Deep Dive into Key Differences

#### 1. Identity Standard & Token Verification

- **Google (OIDC):** Because Google implements OIDC, it returns a cryptographically signed JSON Web Token (`id_token`). Your backend validates the signature, audience (Client ID), issuer, and expiration time completely in-memory using Google's public key infrastructure (`google.golang.org/api/idtoken`).
- **GitHub (OAuth 2.0):** GitHub is not a native OpenID provider for web apps. After exchanging your authorization code for an `access_token`, your server must act as an HTTP client, attach the token as a Bearer header, and query GitHub's REST API (`GET https://api.github.com/user`) to discover who just logged in.

#### 2. Handling Private Emails

- **Google:** Verified user emails are inherently guaranteed and surfaced through the authenticated OIDC claims.
- **GitHub:** Many developers configure their GitHub accounts to keep their email addresses private (`@users.noreply.github.com`). If a standard `/user` request returns a blank or private email string, your service must handle a secondary back-channel request to `https://api.github.com/user/emails`, parsing the array to filter for `Primary: true` and `Verified: true`.

#### 3. Content Negotiation on Token Exchange

- **Google:** The token endpoint `/token` natively expects a standard `POST` request with form parameters and responds with a strict JSON payload.
- **GitHub:** The endpoint `https://github.com/login/oauth/access_token` requires you to explicitly set the HTTP header `Accept: application/json`. If omitted, GitHub responds with a URL-encoded string (`access_token=...&scope=...&token_type=bearer`), which will fail standard JSON decoders.

---

### Implementation Blueprint Summary

- **Google Flow:**
  `Redirect to Google` $\rightarrow$ `Receive Auth Code` $\rightarrow$ `POST /token (Get ID Token + Access Token)` $\rightarrow$ `Verify ID Token Locally (JWT)` $\rightarrow$ `Extract Claims`.
- **GitHub Flow:**
  `Redirect to GitHub` $\rightarrow$ `Receive Auth Code` $\rightarrow$ `POST /access_token (Accept: application/json)` $\rightarrow$ `GET /user (Bearer Token)` $\rightarrow$ `Optional: GET /user/emails (If Primary Email is Private)`.## Architectural Reference: Google OIDC vs. GitHub OAuth 2.0 Integration

This guide documents the core architectural, protocol, and implementation differences between integrating **Google OpenID Connect (OIDC)** and **GitHub OAuth 2.0** authentication flows in Go backend services.

---

### Core Protocol Comparison

| Feature / Dimension        | Google OIDC Authentication                                                                                                                        | GitHub OAuth 2.0 Authentication                                                                                                    |
| :------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------ | :--------------------------------------------------------------------------------------------------------------------------------- |
| **Protocol Specification** | OpenID Connect (OIDC) built on top of OAuth 2.0                                                                                                   | Standard OAuth 2.0 Authorization Code Grant                                                                                        |
| **Token Payload**          | Returns an **`id_token`** (signed JWT containing identity claims) alongside the `access_token`.                                                   | Returns strictly an **`access_token`** (opaque or bearer token). No ID token is issued.                                            |
| **Identity Resolution**    | **Stateless local verification:** The `id_token` can be cryptographically verified offline/locally using public keys (JWKS) via Google's library. | **Stateful remote API call:** Requires an explicit outbound REST API call to `https://api.github.com/user` using the access token. |
| **Email Privacy & Scopes** | Bundled directly into standard scopes (`openid email profile`) and validated within the `id_token`.                                               | Requires explicit scope (`user:email`). Private emails require a secondary fallback API query (`/user/emails`).                    |
| **Token Format Request**   | Standard URL-encoded form POST; responds with JSON automatically.                                                                                 | Requires explicit header `Accept: application/json`, otherwise responds with a URL-encoded query string.                           |

---

### Deep Dive into Key Differences

#### 1. Identity Standard & Token Verification

- **Google (OIDC):** Because Google implements OIDC, it returns a cryptographically signed JSON Web Token (`id_token`). Your backend validates the signature, audience (Client ID), issuer, and expiration time completely in-memory using Google's public key infrastructure (`google.golang.org/api/idtoken`).
- **GitHub (OAuth 2.0):** GitHub is not a native OpenID provider for web apps. After exchanging your authorization code for an `access_token`, your server must act as an HTTP client, attach the token as a Bearer header, and query GitHub's REST API (`GET https://api.github.com/user`) to discover who just logged in.

#### 2. Handling Private Emails

- **Google:** Verified user emails are inherently guaranteed and surfaced through the authenticated OIDC claims.
- **GitHub:** Many developers configure their GitHub accounts to keep their email addresses private (`@users.noreply.github.com`). If a standard `/user` request returns a blank or private email string, your service must handle a secondary back-channel request to `https://api.github.com/user/emails`, parsing the array to filter for `Primary: true` and `Verified: true`.

#### 3. Content Negotiation on Token Exchange

- **Google:** The token endpoint `/token` natively expects a standard `POST` request with form parameters and responds with a strict JSON payload.
- **GitHub:** The endpoint `https://github.com/login/oauth/access_token` requires you to explicitly set the HTTP header `Accept: application/json`. If omitted, GitHub responds with a URL-encoded string (`access_token=...&scope=...&token_type=bearer`), which will fail standard JSON decoders.

---

### Implementation Blueprint Summary

- **Google Flow:**
  `Redirect to Google` $\rightarrow$ `Receive Auth Code` $\rightarrow$ `POST /token (Get ID Token + Access Token)` $\rightarrow$ `Verify ID Token Locally (JWT)` $\rightarrow$ `Extract Claims`.
- **GitHub Flow:**
  `Redirect to GitHub` $\rightarrow$ `Receive Auth Code` $\rightarrow$ `POST /access_token (Accept: application/json)` $\rightarrow$ `GET /user (Bearer Token)` $\rightarrow$ `Optional: GET /user/emails (If Primary Email is Private)`.

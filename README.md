
# Federated Identities

An experimental workspace exploring cross-domain identity federation, OAuth/OIDC flows, and secure token exchanges. This repository pairs a Go-based OIDC and authentication backend with a Next.js TypeScript frontend.

---

## Architecture Overview


```

[ Next.js Frontend ] <--- (OAuth/OIDC Redirect & State CSRF) ---> [ Go OIDC Backend ]
|                                                                  
+--- Playwright End-to-End Tests                               

```

* **`backend/`**: A Go backend handling authentication workflows, session management, and Google OIDC OAuth2 integrations.
* **`frontend/`**: A Next.js (App Router) user interface handling login screens, callbacks, and client-side views.
* **`test/`**: Integration and endpoint tests syncing up with the current API contracts.

---

## Prerequisites

Before setting up your local environment, ensure you have installed:
* **Go** (v1.22 or higher)
* **Node.js** (v18+ or npm-compatible runtime)
* **Git Bash** or a compatible POSIX shell environment (for executing local `.sh` startup scripts)

---

## Environment Variable Configuration

Both the backend and frontend rely on localized environment files to manage secrets. 

### 1. Backend Environment (`backend/.env` or root configuration)
Create a `.env` file inside the `backend` directory (or reference the environment variables expected by your server startup scripts). Based on the GitHub Actions workflow configurations, the core required backend variables include:

| Variable Name | Description | Example / Source |
| :--- | :--- | :--- |
| `GOOGLE_CLIENT_ID` | OAuth Client ID from Google Cloud Console | `your-google-client-id.apps.googleusercontent.com` |
| `GOOGLE_CLIENT_SECRET` | OAuth Client Secret from Google Cloud Console | `GOCSPX-your-secret-string` |
| `GITHUB_CLIENT_ID` | OAuth Client ID from GitHub Developer Settings | `Iv1.yourgithubclientid` |
| `GITHUB_CLIENT_SECRET` | OAuth Client Secret from GitHub Developer Settings | `your-github-client-secret` |


## Obtaining OAuth Credentials

### Getting Google Client ID & Client Secret

1. Go to the [Google Cloud Console](https://console.cloud.google.com/).
2. Create a new project or select an existing one.
3. Navigate to **APIs & Services** > **Credentials**.
4. Click **Create Credentials** > **OAuth client ID**.
5. Select **Web application** as the application type.
6. Add authorized redirect URIs matching your local server routing `http://localhost:3000/provider/google/confirm` 
7. Copy the generated **Client ID** and **Client Secret**.

### Getting GitHub Client ID & Client Secret

1. Go to your [GitHub Developer Settings](https://github.com/settings/developers).
2. Under **OAuth Apps**, click **New OAuth App**.
3. Fill out your application details:
* **Homepage URL:** `http://localhost:3000`
* **Authorization callback URL:** `http://localhost:3000/provider/github/confirm`


4. Click **Register application**.
5. Generate a **Client Secret** and copy both the **Client ID** and **Client Secret**.

---

## Running Locally

The repository includes convenient shell scripts at the root level to fire up services locally:

* **Start the Backend:**
```bash
./start.backend.sh

```


* **Start the Frontend:**
```bash
./start.frontend.sh

```


* **Run Playwright Tests:**
```bash
./start.playwright.sh

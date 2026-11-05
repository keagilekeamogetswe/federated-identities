import { test, expect } from "@playwright/test";

test("load provider confirmation page, mock confirm API, and render user profile", async ({
  page,
}) => {
  // Request phase
  // Navigate to the local page
  await page.goto("http://localhost:3000/");

  // Locate and click the Google button/link
  const googleLink = page.getByRole("button", { name: /google/i });
  await expect(googleLink).toBeVisible();

  // Wait for the navigation to Google's sign-in page
  await Promise.all([
    page.waitForURL(/accounts\.google\.com/, { timeout: 10000 }),
    googleLink.click(),
  ]);
  // Confirmation phase
  // Intercept the background fetch request made by APIInterfaceModule.confirmID
  await page.route("**/api/provider/google/confirm*", async (route) => {
    await route.fulfill({
      status: 200,
      contentType: "application/json",
      json: {
        email: "mrkkeagile@gmail.com",
        google_id: "106186491976426521338",
        name: "Keamogetswe Keagile",
        picture:
          "https://lh3.googleusercontent.com/a/ACg8ocKw99FQAs82VgRhJJ06S_wSRnbjG_UI4art0CSa5DMQNf03m4U=s96-c",
        status: "authenticated",
      },
    });
  });

  // 2. Navigate to the actual Next.js page route (just like your normal browser)
  await page.goto(
    "http://localhost:3000/provider/google/confirm?code=mock_auth_code&state=mock_state",
  );

  // 3. Verify the React component successfully loads, fires the fetch, and renders your UI
  await expect(page.getByText("Authentication Complete")).toBeVisible();
  await expect(page.getByText("Keamogetswe Keagile")).toBeVisible();
  await expect(page.getByText("mrkkeagile@gmail.com")).toBeVisible();
});

import { describe, it, expect } from "vitest";

describe.skip("Testing the initial set-up of Federated Routes Responses", () => {
  it("should validate github provider route response", async () => {
    const provider = "github";
    const url = `http://localhost:8080/provider/${provider}`;
    const expected_res_pattern = /selected\s+github\s+provider/i;
    const post_res = await fetch(url, {
      method: "POST",
    });

    const post_body = await post_res.text();

    expect(post_body).toMatch(expected_res_pattern);
    const get_res = await fetch(url, {
      method: "GET",
    });

    const get_body = await get_res.text();
    expect(get_body).toMatch(expected_res_pattern);
  });
  it("should validate Google provider route response", async () => {
    const provider = "google";
    const url = `http://localhost:8080/provider/${provider}`;
    const expected_res_pattern = /selected\s+google\s+provider/i;
    const post_res = await fetch(url, {
      method: "POST",
    });

    const post_body = await post_res.text();

    expect(post_body).toMatch(expected_res_pattern);
    const get_res = await fetch(url, {
      method: "GET",
    });

    const get_body = await get_res.text();
    expect(get_body).toMatch(expected_res_pattern);
  });
  it("should validate unknown provider route response", async () => {
    const provider = "bad-provider";
    const url = `http://localhost:8080/provider/${provider}`;
    const expected_res_pattern = /Unknown\s+provider/i;
    const post_res = await fetch(url, {
      method: "POST",
    });

    const post_body = await post_res.text();

    expect(post_body).toMatch(expected_res_pattern);
    const get_res = await fetch(url, {
      method: "GET",
    });

    const get_body = await get_res.text();
    expect(get_body).toMatch(expected_res_pattern);
  });
});

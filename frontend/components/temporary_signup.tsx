"use client";

import { SubmitEvent } from "react";

const onSubmit = async (
  window: Window,
  event: SubmitEvent<HTMLFormElement>,
) => {
  event.preventDefault();
  const formData = new FormData(event.currentTarget);
  const response = await fetch("/auth/v1/tmp/signup/", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(Object.fromEntries(formData)),
  });
  const data = await response.json();
  if (response.ok) {
    window.location.href = "/auth/signup/?token=" + data.signup_token;
  }
};

export const Signup = () => {
  return (
    <form onSubmit={(e) => onSubmit(window, e)}>
      <h1 className="text-2xl font-bold mb-4">Sign Up</h1>
      <div className="mb-4">
        <label className="block mb-2" htmlFor="email">
          Email
        </label>
        <input
          className="w-full p-2 border border-gray-300 rounded"
          type="email"
          id="email"
          name="email"
          required
        />
      </div>
      <div className="mb-4">
        <label className="block mb-2" htmlFor="password">
          Password
        </label>
        <input
          className="w-full p-2 border border-gray-300 rounded"
          type="password"
          id="password"
          name="password"
          required
        />
      </div>
      <button
        className="w-full bg-blue-500 text-white p-2 rounded"
        type="submit"
      >
        Sign Up
      </button>
    </form>
  );
};

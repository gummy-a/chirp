"use client";

import { SubmitEvent, useState } from "react";
import { postApiAuthV1Login } from "@/lib/client/auth/v1/sdk.gen";
import { useRouter } from "next/navigation";

const onSubmit = async (event: SubmitEvent<HTMLFormElement>) => {
  event.preventDefault();
  const formData = new FormData(event.currentTarget);

  const ret = await postApiAuthV1Login({
    body: {
      email: formData.get("email") as string,
      password: formData.get("password") as string,
    },
    throwOnError: false,
    credentials: "include",
  });

  return ret;
};

export const LoginForm = () => {
  const [msg, setMsg] = useState(<></>);
  const router = useRouter();

  const submitProcess = async (e: SubmitEvent<HTMLFormElement>) => {
    const ret = await onSubmit(e);

    if (ret.response.ok) {
      router.push("/");
    } else {
      setMsg(<div className="text-red-500">Login failed!</div>);
    }
  };

  return (
    <form onSubmit={(e) => submitProcess(e)}>
      <h1 className="text-2xl font-bold mb-4">Login</h1>
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
      <div className="mb-4">{msg}</div>
      <button
        className="w-full bg-blue-500 text-white p-2 rounded"
        type="submit"
      >
        Login
      </button>
    </form>
  );
};

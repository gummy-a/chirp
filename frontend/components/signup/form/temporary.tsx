"use client";

import { SubmitEvent, useRef, useState } from "react";
import { postApiAuthV1TmpSignup } from "@/lib/client/signup/sdk.gen";
import { useRouter } from "next/navigation";

const onSubmit = async (event: SubmitEvent<HTMLFormElement>) => {
  event.preventDefault();
  const formData = new FormData(event.currentTarget);

  const ret = await postApiAuthV1TmpSignup({
    body: {
      email: formData.get("email") as string,
      password: formData.get("password") as string,
    },
    throwOnError: false,
  });

  return ret;
};

export const TemporarySignup = () => {
  const ref = useRef<HTMLFormElement>(null);
  const [msg, setMsg] = useState(<></>);
  const router = useRouter();

  const submitProcess = async (e: SubmitEvent<HTMLFormElement>) => {
    const ret = await onSubmit(e);

    if (ret.response.ok) {
      router.push(`/signup/?token=${ret.data?.signup_token}`);
      ref.current?.reset();
    } else {
      setMsg(<div className="text-red-500">Signup failed!</div>);
    }
  };

  return (
    <form ref={ref} onSubmit={(e) => submitProcess(e)}>
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
      <div className="mb-4">{msg}</div>
      <button
        className="w-full bg-blue-500 text-white p-2 rounded"
        type="submit"
      >
        Sign Up
      </button>
    </form>
  );
};

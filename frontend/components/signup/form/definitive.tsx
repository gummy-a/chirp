"use client";

import { SubmitEvent, useRef, useState } from "react";
import { postApiAuthV1Signup } from "@/lib/client/signup/sdk.gen";
import { useSearchParams } from "next/navigation";
import { useRouter } from "next/navigation";

const onSubmit = async (token: string, event: SubmitEvent<HTMLFormElement>) => {
  event.preventDefault();
  const formData = new FormData(event.currentTarget);

  const ret = await postApiAuthV1Signup({
    body: {
      signup_token: token,
      number_code: parseInt(formData.get("number_code") as string),
    },
    throwOnError: false,
    credentials: "include",
  });

  return ret;
};

export const Signup = () => {
  const ref = useRef<HTMLFormElement>(null);
  const [msg, setMsg] = useState(<></>);
  const router = useRouter();
  const param = useSearchParams();
  const token = param.get("token") || "";

  const submitProcess = async (e: SubmitEvent<HTMLFormElement>) => {
    const ret = await onSubmit(token, e);

    if (ret.response.ok) {
      router.push("/");
    } else {
      setMsg(<div className="text-red-500">Signup failed!</div>);
      ref.current?.reset();
    }
  };

  return (
    <form ref={ref} onSubmit={(e) => submitProcess(e)}>
      <h1 className="text-2xl font-bold mb-4">Confirm</h1>
      <div className="mb-4">
        <label className="block mb-2" htmlFor="number_code">
          input 6 number code
        </label>
        <input
          className="w-full p-2 border border-gray-300 rounded"
          type="number"
          id="number_code"
          name="number_code"
          min={100000}
          max={999999}
          required
        />
      </div>
      <div className="mb-4">{msg}</div>
      <button
        className="w-full bg-blue-500 text-white p-2 rounded"
        type="submit"
      >
        Register
      </button>
    </form>
  );
};

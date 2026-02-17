"use client";

import { TemporarySignup } from "./form/temporary";
import { Signup } from "./form/definitive";
import { useEffect, useState } from "react";
import { getApiAuthV1TmpAccountById } from "@/lib/client/auth/v1/sdk.gen";
import { useRouter, useSearchParams } from "next/navigation";

export const ShowSignupForm = () => {
  const [element, setElement] = useState(<></>);
  const param = useSearchParams();
  const token = param.get("token") || "";
  const router = useRouter();

  useEffect(() => {
    (async () => {
      try {
        if (localStorage.getItem("jwt_token")) {
          router.push("/");
        }

        const ret = await getApiAuthV1TmpAccountById({
          path: {
            id: token,
          },
          throwOnError: false,
          baseUrl: process.env.NEXT_PUBLIC_API_BASE_URL,  
        });

        if (ret.response.ok && token) {
          setElement(<Signup />);
        } else {
          setElement(<TemporarySignup />);
        }
      } catch {
        setElement(<TemporarySignup />);
      }
    })();
  }, [token, router]);

  return element;
};

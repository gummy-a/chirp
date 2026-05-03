"use client";

import { TemporarySignup } from "./form/temporary";
import { Signup } from "./form/definitive";
import { useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";

export const ShowSignupForm = () => {
  const [element, setElement] = useState(<></>);
  const param = useSearchParams();
  const token = param.get("token") || "";
  const router = useRouter();

  useEffect(() => {
    (async () => {
      try {
        const ret = await fetch(`/api/auth/v1/account/tmp/${token}/`, {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
          },
        });

        if (ret.ok && token) {
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

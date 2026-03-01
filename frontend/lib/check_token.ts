import { cookies } from "next/headers";
import { ConstKeySession } from "./constant_variable";
import { redirect } from "next/navigation";

export const RedirectLoginIfSessionHasNotSet = async () => {
  const cookie = await cookies();
  if (!cookie.get(ConstKeySession)) {
    redirect("/login/");
  }
};

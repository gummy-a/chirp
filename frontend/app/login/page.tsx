import { LoginForm } from "@/components/login/form";
import { ConstKeySession } from "@/lib/constant_variable";
import { cookies } from "next/headers";
import { redirect } from "next/navigation";

export default async function Home() {
  const cookie = await cookies();
  if (cookie.get(ConstKeySession)) {
    redirect("/");
  }

  return <LoginForm />;
}

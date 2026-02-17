import { LoginForm } from "@/components/login/form";
import { cookies } from "next/headers";
import { redirect } from "next/navigation";

export default async function Home() {
  const cookie = await cookies();
  if (cookie.get("session")) {
    redirect("/");
  }

  return <LoginForm />;
}

import { ShowSignupForm } from "@/components/signup/show_form";
import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { Suspense } from "react";

export default async function Home() {
  const cookie = await cookies();
  if (cookie.get("session")) {
    redirect("/");
  }

  return (
    <Suspense>
      <ShowSignupForm />
    </Suspense>
  );
}

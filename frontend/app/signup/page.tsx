import { ShowSignupForm } from "@/components/signup/show_form";
import { Suspense } from "react";

export default function Home() {
  return (
    <Suspense>
      <ShowSignupForm />
    </Suspense>
  );
}

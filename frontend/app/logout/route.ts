import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { postApiAuthV1Logout } from "@/lib/client/auth/v1/sdk.gen";

export async function GET() {
  const cookieStore = await cookies();
  const token = cookieStore.get("session");

  if (!token) {
    return;
  }

  const ret = await postApiAuthV1Logout({
    body: {
      session: token.value,
    },
    throwOnError: false,
    baseUrl: process.env.NEXT_PUBLIC_API_BASE_URL,
  });

  if (ret.response.ok) {
    cookieStore.delete("session");
    redirect("/");
  }
}

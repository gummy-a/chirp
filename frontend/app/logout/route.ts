import { ConstKeySession } from "@/lib/constant_variable";
import { cookies } from "next/headers";
import { redirect } from "next/navigation";

export async function GET() {
  const cookieStore = await cookies();
  cookieStore.delete(ConstKeySession);
  redirect("/");
}

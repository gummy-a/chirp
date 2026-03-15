import { ConstKeySession } from "@/lib/constant_variable";
import { cookies } from "next/headers";

export async function GET(
  req: Request,
  { params }: { params: { job_id: string } },
) {
  try {
    /*
    Server-Send-Event request can't set authorization header, so
    we get authorization token via cookie, not from http header
    */
    const cookie = await cookies();
    const authorization = cookie.get(ConstKeySession);
    if (!authorization) {
      throw new Error("authorization header is required");
    }
    const headers = new Headers();
    headers.set("authorization", "Bearer " + authorization.value);

    const param = await params;
    const ret = await fetch(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/media/v1/upload/status/${param.job_id}/`,
      {
        method: "GET",
        headers: headers,
      },
    );

    return new Response(ret.body, {
      status: ret.ok ? 200 : 400,
      headers: ret.headers,
    });
  } catch (e) {
    console.error(e);
    return new Response(JSON.stringify({ error: e }), {
      status: 400,
      headers: {
        "Content-Type": "application/json",
      },
    });
  }
}

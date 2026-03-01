import { postApiMediaV1Upload } from "@/lib/client/media/v1/sdk.gen";

export async function POST(req: Request) {
  const formData = await req.formData();
  const files = formData.getAll("file") as File[];
  const authorizationHeader = req.headers.get("Authorization");

  try {
    const ret = await postApiMediaV1Upload({
      body: {
        files: files,
      },
      headers: {
        Authorization: authorizationHeader,
      },
      throwOnError: false,
      baseUrl: process.env.NEXT_PUBLIC_API_BASE_URL,
    });

    if (ret.response.ok) {
      return new Response(JSON.stringify(ret.data), {
        status: 200,
        headers: {
          "Content-Type": "application/json",
        },
      });
    }

    return new Response(await ret.response.json(), {
      status: 400,
      headers: {
        "Content-Type": "application/json",
      },
    });
  } catch (e) {
    return new Response(JSON.stringify({ error: e }), {
      status: 400,
      headers: {
        "Content-Type": "application/json",
      },
    });
  }
}

export async function POST(req: Request) {
  const formData = await req.formData();

  try {
    const headers = new Headers();
    const auth = req.headers.get("authorization");
    if (auth) {
      headers.set("authorization", auth);
    } else {
      throw new Error("authorization header is required");
    }

    const ret = await fetch(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/media/v1/upload/`,
      {
        method: "POST",
        body: formData,
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

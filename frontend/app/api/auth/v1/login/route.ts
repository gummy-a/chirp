export async function POST(req: Request) {
  try {
    const formData = await req.formData();
    const ret = await fetch(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/auth/v1/login/`,
      {
        method: "POST",
        body: formData,
        credentials: "include",
      },
    );

    return new Response(ret.body, {
      status: ret.ok ? 200 : 400,
      headers: ret.headers,
    });
  } catch (e) {
    console.error(e);
    return new Response(JSON.stringify({ error: "error" }), {
      status: 400,
      headers: {
        "Content-Type": "application/json",
      },
    });
  }
}

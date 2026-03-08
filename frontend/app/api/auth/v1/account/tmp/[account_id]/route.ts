export async function GET(
  req: Request,
  { params }: { params: { account_id: string } },
) {
  try {
    const param = await params;
    const ret = await fetch(
      `${process.env.NEXT_PUBLIC_API_BASE_URL}/api/auth/v1/account/tmp/${param.account_id}/`,
      {
        method: "GET",
      },
    );

    return new Response(ret.body, {
      status: ret.ok ? 200 : 400,
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

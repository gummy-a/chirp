import { cookies } from 'next/headers';
import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';
import { ConstKeySession } from './lib/constant_variable';

export async function proxy(request: NextRequest) {
  const requestHeaders = new Headers(request.headers);
  
  const cookie = await cookies();
  const token = cookie.get(ConstKeySession);
  requestHeaders.set('Authorization', `Bearer ${token?.value}`);

  return NextResponse.next({
    request: {
      headers: requestHeaders,
    },
  });
}

export const config = {
  matcher: [
    '/api/media/:path*',
  ],
};
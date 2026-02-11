"use client";

import { client } from "@/lib/client/signup/client.gen";

export default function Initializer() {
  client.setConfig({
    baseUrl: process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080",
  });

  return null;
}

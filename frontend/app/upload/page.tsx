import { UploadForm } from "@/components/upload/form";
import { RedirectLoginIfSessionHasNotSet } from "@/lib/check_token";

export default async function Home() {
  await RedirectLoginIfSessionHasNotSet();

  return <UploadForm />;
}

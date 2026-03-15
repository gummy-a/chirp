import { useEffect, useState } from "react";
import { ApiResponse } from "./form";

type encodeProcess = {
  msg: string;
  media_id: string | undefined;
  original_file_name: string | undefined;
};

type Props = {
  json: ApiResponse[];
};

type Progress = {
  message: string;
  original_file_name: string;
};

export const Progress = ({ json }: Props) => {
  const [msg, setMsg] = useState<Progress[]>([]);

  useEffect(() => {
    if (json.length === 0) {
      return;
    }

    const sources: EventSource[] = [];
    const progress: Progress[] = [];

    for (const j of json) {
      const es = new EventSource(`/api/media/v1/status/${j.job_id}/`);
      sources.push(es);
      progress.push({
        message: `uploading ${j.original_file_name} ...`,
        original_file_name: j.original_file_name,
      });

      es.onmessage = (e) => {
        const data = JSON.parse(e.data) as encodeProcess;

        setMsg((prev) =>
          prev.map((e) => {
            if (e.original_file_name === data.original_file_name) {
              e.message = `upload finished ${data.original_file_name}.`;
            }
            return e;
          }),
        );
      };
    }
    setMsg(progress);

    return () => {
      sources.forEach((s) => s.close());
    };
  }, [json]);

  return (
    <>
      {msg.map((m) => (
        <div key={m.original_file_name}>{m.message}</div>
      ))}
    </>
  );
};

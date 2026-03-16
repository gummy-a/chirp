import { useEffect, useState } from "react";
import { ApiResponse } from "./form";

type SSEResponse = {
  msg: string;
  job_id: string;
  media_id: string | undefined;
};

type Props = {
  json: ApiResponse[];
};

type ProgressMessage = {
  msg: string;
  job_id: string;
};

export const SSEProgressMessage = ({ json }: Props) => {
  const [msg, setMsg] = useState<ProgressMessage[]>([]);

  useEffect(() => {
    if (json.length === 0) {
      return;
    }

    const es = new EventSource(`/api/media/v1/status/`);

    for (const j of json) {
      setMsg((prev) => [
        ...prev,
        {
          msg: `uploading ${j.original_file_name} ...`,
          job_id: j.job_id,
        },
      ]);
    }

    es.onerror = () => {
      setMsg([{msg: "something went wrong.", job_id: ""}])
      es.close();
    }

    es.onmessage = (e) => {
      const data = JSON.parse(e.data) as SSEResponse;

      setMsg((prev) =>
        prev.map((v) => {
          if (v.job_id === data.job_id) {
            return {
              ...v,
              msg: data.msg,
            };
          }
          return v;
        }),
      );

      if (data.media_id !== undefined) {
        es.close();
      }
    };

    return () => {
      es.close();
    };
  }, [json]);

  return (
    <>
      {msg.map((m) => (
        <div key={m.job_id}>{m.msg}</div>
      ))}
    </>
  );
};

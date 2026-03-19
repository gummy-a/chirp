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

const sseCloseConnectionMessage = "--- end SSE connection---";

export const SSEProgressMessage = ({ json }: Props) => {
  const [msg, setMsg] = useState<ProgressMessage[]>([]);
  let finishedCount = 0;

  useEffect(() => {
    if (json.length === 0) {
      return;
    }

    // create job message list
    for (const j of json) {
      setMsg((prev) => [
        ...prev,
        {
          msg: `uploading ${j.original_file_name} ...`,
          job_id: j.job_id,
        },
      ]);
    }

    const es = new EventSource(`/api/media/v1/status/`);

    es.onerror = () => {
      setMsg((prev) => [
        ...prev,
        { msg: "unexpected error occured.", job_id: "error" },
      ]);
      es.close();
    };

    es.onmessage = (e) => {
      const data = JSON.parse(e.data) as SSEResponse;

      // each jobs will send a close connection message
      if (data.msg === sseCloseConnectionMessage) {
        finishedCount++;
        if (finishedCount === json.length) {
          setMsg((prev) => [
            ...prev,
            { msg: "connection closed.", job_id: "finished" },
          ]);
          es.close();
        }
        return;
      } else {
        // update job message
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

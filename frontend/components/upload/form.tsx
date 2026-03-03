"use client";

import "./button.css";
import { SubmitEvent, useRef, useState } from "react";

const onSubmit = async (event: SubmitEvent<HTMLFormElement>) => {
  event.preventDefault();
  const formData = new FormData(event.currentTarget);

  const ret = await fetch("/api/media/v1/upload", {
    method: "POST",
    body: formData,
  });

  return await ret.json();
};

export const UploadForm = () => {
  const ref = useRef<HTMLFormElement>(null);
  const [disabled, setDisabled] = useState(false);

  const submit = async (e: SubmitEvent<HTMLFormElement>) => {
    setDisabled(true);
    await onSubmit(e);
    setDisabled(false);
    ref.current?.reset();
  };

  return (
    <form ref={ref} encType="multipart/form-data" onSubmit={(e) => submit(e)}>
      <h1 className="text-2xl font-bold mb-4">File Upload</h1>
      <div className="mb-4">
        <label htmlFor="file">select upload files</label>
      </div>
      <div className="mb-4">
        <input
          className="border rounded w-full p-2 cursor-pointer"
          type="file"
          id="file"
          name="file"
          multiple
        />
      </div>
      <div>
        <button
          className={`w-full ${disabled ? "bg-gray-500" : "bg-blue-500"} text-white p-2 rounded`}
          type="submit"
          disabled={disabled}
        >
          {disabled ? "uploading..." : "submit"}
        </button>
      </div>
    </form>
  );
};

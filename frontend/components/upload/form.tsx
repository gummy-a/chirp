"use client";

import "./button.css";
import { SubmitEvent } from "react";

const onSubmit = async (event: SubmitEvent<HTMLFormElement>) => {
  event.preventDefault();
  const formData = new FormData(event.currentTarget);

  const ret = await fetch("/api/media/v1/upload", {
    method: "POST",
    body: formData,
  });
  const data = await ret.json();
  console.log(data);
};

export const UploadForm = () => {
  return (
    <form encType="multipart/form-data" onSubmit={(e) => onSubmit(e)}>
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
          className="w-full bg-blue-500 text-white p-2 rounded"
          type="submit"
        >
          submit
        </button>
      </div>
    </form>
  );
};

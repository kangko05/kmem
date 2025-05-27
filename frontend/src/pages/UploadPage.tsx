import { type ChangeEvent } from "react";

import { PageLayout, FileInput } from "../components";
import { axiosInstance } from "../utils";

export const UploadPage = () => {
  const handleFileChange = async (ev: ChangeEvent<HTMLInputElement>) => {
    const files = ev.target.files;

    if (!files) return;

    for (const file of files) {
      const resp = await axiosInstance.post("/files/upload", file);
      console.log(file.name, ":", resp.status);
    }
  };

  return (
    <PageLayout>
      <FileInput onChange={handleFileChange} />
    </PageLayout>
  );
};

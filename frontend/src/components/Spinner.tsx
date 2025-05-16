import { BounceLoader } from "react-spinners";

export const Spinner = ({ loading }: { loading: boolean }) => {
  return <BounceLoader color={"oklch(0.809 0.105 251.813)"} loading={loading} />;
};

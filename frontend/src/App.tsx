import { createBrowserRouter, RouterProvider } from "react-router";
import { HomePage, LoginPage, MorePage, SearchedPage, UploadPage } from "./pages";

const router = createBrowserRouter([
  { path: "/", element: <LoginPage /> },
  { path: "/home", element: <HomePage /> },
  { path: "/more", element: <MorePage /> },
  { path: "/upload", element: <UploadPage /> },
  { path: "/search", element: <SearchedPage /> },
]);

const App = () => {
  return <RouterProvider router={router} />;
};

export default App;

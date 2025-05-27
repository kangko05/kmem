import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { LoginPage, HomePage, UploadPage } from "./pages";

const router = createBrowserRouter([
  { path: "/", element: <div>Hello World!</div> },
  { path: "/login", element: <LoginPage /> },
  { path: "/home", element: <HomePage /> },
  { path: "/upload", element: <UploadPage /> },
]);

function App() {
  return <RouterProvider router={router} />;
}

export default App;

import { createBrowserRouter, RouterProvider } from "react-router";
import { HomePage, LoginPage, MorePage } from "./pages";

const router = createBrowserRouter([
  { path: "/", element: <LoginPage /> },
  { path: "/home", element: <HomePage /> },
  { path: "/more", element: <MorePage /> },
]);

const App = () => {
  return <RouterProvider router={router} />;
};

export default App;

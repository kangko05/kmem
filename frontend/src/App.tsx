import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { LoginPage } from "./pages";
import { HomePage } from "./pages/HomePage";

const router = createBrowserRouter([
  { path: "/", element: <div>Hello World!</div> },
  { path: "/login", element: <LoginPage /> },
  { path: "/home", element: <HomePage /> },
]);

function App() {
  return <RouterProvider router={router} />;
}

export default App;

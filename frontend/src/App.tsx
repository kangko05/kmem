import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { LoginPage, HomePage, UploadPage, GalleryPage } from "./pages";
import { QueryClient, QueryClientProvider } from "react-query";

const queryClient = new QueryClient();

const router = createBrowserRouter([
  { path: "/", element: <div>Hello World!</div> },
  { path: "/login", element: <LoginPage /> },
  { path: "/home", element: <HomePage /> },
  { path: "/upload", element: <UploadPage /> },
  { path: "/gallery", element: <GalleryPage /> },
]);

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>
  );
}

export default App;

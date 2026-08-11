import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { createHashRouter, RouterProvider } from "react-router-dom";

import App from "./pages/App.jsx";
import Mode from "./pages/Mode.jsx"
import ExistingClient from "./pages/ExistingClient.jsx"
import NewClient from "./pages/NewClient.jsx"

// use react router dom  to render .jsx onto single page, SPA framework
const router = createHashRouter([
  {
    path: "/",
    element: <App />,
    children: [
      { index: true, element: <Mode /> },
      { path: "existing", element: <ExistingClient /> },
      { path: "new", element: <NewClient /> },
    ],
  },
]);

createRoot(document.getElementById("root")).render(
  <StrictMode>
    <RouterProvider router={router} />
  </StrictMode>,
);

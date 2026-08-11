import { Outlet } from "react-router-dom";

// supports SPA rendering, Outlet renders .jsx files accordingly
function App() {
  return (
    <div>
      <main>
        <Outlet />
      </main>
    </div>
  );
}

export default App;

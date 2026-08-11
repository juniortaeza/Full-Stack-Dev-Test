import { Link } from "react-router-dom";

import "../style/Mode.css";

// use Link to render .jsx onto page for efficiency
function Mode() {
  return (
    <div className="mode-container">
      <div className="mode-card">
        <h3>Choose an option</h3>
        <div className="mode-buttons">
          <Link to="/existing">Existing Client</Link>
          <Link to="/new">New Client</Link>
        </div>
      </div>
    </div>
  );
}

export default Mode;

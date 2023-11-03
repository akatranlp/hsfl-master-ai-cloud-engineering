import { Link } from "react-router-dom";

const App = () => {
  return (
    <div>
      Hallo App
      <div>
        <Link to="/books">Books</Link>
      </div>
    </div>
  );
};

export default App;

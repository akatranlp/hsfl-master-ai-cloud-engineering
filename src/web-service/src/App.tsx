import { Navigate } from "react-router-dom";

const App = () => {
  const shouldRedirect = true;

  if (shouldRedirect) {
    return <Navigate to="/books" replace />;
  }

  return <></>;
};

export default App;

import "./App.css";
import axios from "axios";
import { useEffect, useState } from "react";

function App() {
  const [data, setData] = useState("");
  useEffect(() => {
    axios
      .get("http://localhost:8080/")
      .then((response) => {
        console.log("Data fetched:", response);
        setData(response);
      })
      .catch((error) => {
        console.error("Error fetching data:", error);
      });
  }, []);
  return (
    <>
      <h1>App</h1>
      <pre>{JSON.stringify(data, null, 2)}</pre>
    </>
  );
}

export default App;

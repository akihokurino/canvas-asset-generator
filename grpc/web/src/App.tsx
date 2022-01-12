import React from "react";
import logo from "./logo.svg";
import "./App.css";
import { InternalAPIClient } from "./rpc/api_pb_service";
import { SignedGsUrlsRequest } from "./rpc/api_pb";

function App() {
  React.useEffect(() => {
    const client = new InternalAPIClient("https://grpc-dot-canvas-329810.an.r.appspot.com");
    const req = new SignedGsUrlsRequest();
    req.setGsUrlsList([
      "gs://canvas-329810-video/ColorRain.mp4"
    ])
    client.signedGsUrls(req, (err, res) => {
      console.log(err);
      res?.getUrlsList().forEach((url) => {
        console.log(url);
      })
    });
  }, []);

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Edit <code>src/App.tsx</code> and save to reload.
        </p>
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn React
        </a>
      </header>
    </div>
  );
}

export default App;

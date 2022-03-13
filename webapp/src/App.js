import "./site.scss";
import './App.css';
import {useEffect} from "react";
import Header from "./Header.jsx";
import Generator from "./Generator.jsx";
import Footer from "./Footer.jsx";
import GitHubCorners from '@uiw/react-github-corners';
import Container from "react-bootstrap/Container";

const pageTitle = "Docker Task Generator";

function App() {
  useEffect(() => {
    document.title = pageTitle;
  }, []);

  return (
    <div className="App">
      <GitHubCorners
        position="right"
        href="https://github.com/uiwjs/react-github-corners"
      />
      <Container>
        <Header title={pageTitle} />
        <Generator />
        <Footer />
      </Container>
    </div>
  );
}

export default App;

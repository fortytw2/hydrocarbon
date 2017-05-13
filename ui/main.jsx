import { Component, h, render } from "preact";
import { Route, Router } from "preact-router";
import { Store, history } from "./state/Store";

import Footer from "./components/Footer";
import Login from "./components/Login";
import Nav from "./components/Nav";
import NotFound from "./components/NotFound";
import Redux from "preact-redux";
import TextContent from "./components/TextContent";
import { initDevTools } from "./vendor/devtools";

initDevTools();

const App = function() {
  return (
    <Redux.Provider store={Store}>
      <div class="min-vh-100">
        <Nav />
        <Router history={history}>
          <Route
            path="/"
            component={TextContent}
            text="this is the home page"
          />
          <Route
            path="/about"
            component={TextContent}
            text="hi this is about page"
          />
          <Route path="/login" component={Login} />
          <Route component={NotFound} default />
        </Router>
        <Footer />
      </div>
    </Redux.Provider>
  );
};

render(<App />, document.body);

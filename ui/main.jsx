import { Component, h, render } from "preact";
import { History, Store } from "./state/store";
import { Route, Router } from "preact-router";

import Footer from "./components/footer";
import Login from "./components/login";
import LoginCallback from "./components/login_callback";
import Nav from "./components/nav";
import NotFound from "./components/not_found";
import NotificationWindow from "./components/notification_window";
import RehydrateProvider from "./components/rehydrate_provider";
import TextContent from "./components/text_content";
import { initDevTools } from "./vendor/devtools";

initDevTools();

const App = function() {
  return (
    <RehydrateProvider store={Store}>
      <div class="min-vh-100">
        <Nav />
        <NotificationWindow />
        <Router history={History}>
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
          <Route path="/login-callback" component={LoginCallback} />
          <Route component={NotFound} default />
        </Router>
        <Footer />
      </div>
    </RehydrateProvider>
  );
};

render(<App />, document.body);

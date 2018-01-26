import { h, Component } from "preact";
import { Router } from "preact-router";

import Layout from "./layout";
import Home from "../routes/home";
import Profile from "../routes/profile";
import NotFound from "../routes/notfound";

import Feed from "async!../routes/feed";
// import Profile from 'async!../routes/profile';

export default class App extends Component {
  /** Gets fired when the route changes.
   *	@param {Object} event		"change" event from [preact-router](http://git.io/preact-router)
   *	@param {string} event.url	The newly routed URL
   */
  handleRoute = e => {
    this.currentUrl = e.url;
  };

  render() {
    return (
      <div id="app">
        <Layout>
          <Router onChange={this.handleRoute}>
            <Home path="/" />
            <Profile path="/profile/" user="me" />
            <Profile path="/profile/:user" />
            <Feed path="/feed/" id="0" />
            <Feed path="/feed/:id" />
            <NotFound default />
          </Router>
        </Layout>
      </div>
    );
  }
}

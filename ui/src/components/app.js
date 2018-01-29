import { h, Component } from "preact";
import { Router } from "preact-router";

import Layout from "./layout";
import Home from "../routes/home";
import Profile from "../routes/profile";
import NotFound from "../routes/notfound";
import Login from "../routes/login";
import FolderList from "../routes/folderlist";
import Feed from "../routes/feed";
// import Profile from 'async!../routes/profile';

import "preact/devtools";

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
            <FolderList path="/folders/" id="0" />
            <FolderList path="/folders/:id" />
            <FolderList path="/folders/:id/:feedID" />
            <Login path="/login" />
            <Login path="/login-callback" callback={true} />
            <NotFound default />
          </Router>
        </Layout>
      </div>
    );
  }
}

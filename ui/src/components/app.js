import { h, Component } from "preact";
import { Router } from "preact-router";

import Layout from "./layout";
import Home from "../routes/home";
import Profile from "../routes/profile";
import NotFound from "../routes/notfound";
import Login from "../routes/login";
import Logout from "../routes/logout";
import FolderList from "../routes/folderlist";
// import Profile from 'async!../routes/profile';

import "preact/devtools";
import style from "./style";

export default class App extends Component {
  constructor(props) {
    super(props);

    window.baseURL = "";

    this.setState({ loggedIn: false });

    this.loginSwap.bind(this);
  }

  componentWillMount() {
    this.setState({ loggedIn: this.isLoggedIn() });
  }

  isLoggedIn = () => {
    try {
      let key = window.localStorage.getItem("hydrocarbon-key");
      return key !== null;
    } catch (e) {
      return false;
    }
  };

  loginSwap() {
    this.setState({ loggedIn: this.isLoggedIn() });
  }

  handleRoute = e => {
    this.currentUrl = e.url;
  };

  render({}, { loggedIn }) {
    return (
      <div id="app" class={style.themed}>
        <Layout loggedIn={loggedIn}>
          <Router onChange={this.handleRoute}>
            <Home path="/" />
            <Profile path="/profile/" user="me" />
            <Profile path="/profile/:user" />
            <FolderList path="/folders/" id="0" />
            <FolderList path="/folders/:id" />
            <FolderList path="/folders/:id/:feedID" />
            <Login path="/login" />
            <Login
              path="/login-callback"
              callback={true}
              loginSwapper={() => {
                this.loginSwap();
              }}
            />
            <Logout
              path="/logout"
              loginSwapper={() => {
                this.loginSwap();
              }}
            />
            <NotFound default />
          </Router>
        </Layout>
      </div>
    );
  }
}

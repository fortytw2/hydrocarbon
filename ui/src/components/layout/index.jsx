import { h, Component } from "preact";
import Router from "preact-router";

import Header from "@/components/navbar";
import Home from "@/routes/home";
import Feed from "@/routes/feed";
import Settings from "@/routes/settings";
import Login from "@/routes/login";
import Callback from "@/routes/callback";

import style from "./style.css";

require("preact/debug");

export default class Layout extends Component {
  constructor(props) {
    super(props);

    this.setState({
      email: null,
      loggedIn: false,
      apiKey: null
    });
  }

  render({}, { email, loggedIn, apiKey }) {
    return (
      <div class={style.layout}>
        <Header />
        <Router>
          <Home path="/" />
          <Feed path="/feed" />
          <Settings path="/settings" />
          <Login path="/login" />
          <Callback path="/callback" />
        </Router>
      </div>
    );
  }
}

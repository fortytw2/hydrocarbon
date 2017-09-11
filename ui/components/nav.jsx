import { Component, h } from "preact";

import { Link } from "preact-router/match";
import Logout from "./logout";
import Redux from "preact-redux";

class Nav extends Component {
  loggedIn() {
    if (this.props.login.email !== "") {
      return true;
    }
    return false;
  }
  render(props, state) {
    return (
      <nav class="pa1 pa2-ns">
        <Link
          class="link dim black b f6 f5-ns dib mr3"
          activeClassName="blue"
          href="/"
        >
          hydrocarbon
        </Link>
        {this.loggedIn()
          ? <div class="dib">
              <Link
                class="link dim gray f6 f5-ns dib mr3"
                activeClassName="blue"
                href="/feed"
              >
                {props.login.email}
              </Link>
              <Logout
                class="link dim gray f6 f5-ns dib mr3"
                apiKey={props.login.apiKey}
              />
            </div>
          : <Link
              class="link dim gray f6 f5-ns dib mr3"
              activeClassName="blue"
              href="/login"
            >
              login
            </Link>}

      </nav>
    );
  }
}

const mapStateToProps = state => {
  return {
    login: state.login
  };
};

export default Redux.connect(mapStateToProps)(Nav);

import { Component, h, render } from "preact";

import { Link } from "preact-router/match";

class Nav extends Component {
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
        <Link
          class="link dim gray f6 f5-ns dib mr3"
          activeClassName="blue"
          href="/login"
        >
          login
        </Link>
      </nav>
    );
  }
}

export default Nav;

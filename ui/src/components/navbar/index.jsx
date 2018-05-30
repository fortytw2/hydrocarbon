import { h, Component } from "preact";
import { Link } from "preact-router/match";
import style from "./style.css";

export default class Header extends Component {
  renderUserNav(email, logoutCallback) {
    if (!email) {
      return (
        <nav>
          <Link tabIndex="0" activeClassName={style.active} href="/login">
            Login
          </Link>
        </nav>
      );
    }

    return (
      <nav>
        <Link tabIndex="0" activeClassName={style.active} href="/feed">
          Feed
        </Link>
        <Link tabIndex="0" activeClassName={style.active} href="/settings">
          Settings
        </Link>
        <a class={style.logout} onClick={logoutCallback}>
          Logout {email}
        </a>
      </nav>
    );
  }

  render({ email, logoutCallback }, {}) {
    return (
      <header class={style.header}>
        <Link tabIndex="0" href="/">
          <h1>Hydrocarbon</h1>
        </Link>
        {this.renderUserNav(email, logoutCallback)}
      </header>
    );
  }
}

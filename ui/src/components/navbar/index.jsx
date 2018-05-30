import { h, Component } from "preact";
import { Link } from "preact-router/match";
import style from "./style.css";

export default class Header extends Component {
  renderUserNav(loggedIn, email, logoutCallback) {
    if (!loggedIn) {
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
          Logout
        </a>
      </nav>
    );
  }

  render({ loggedIn, email, logoutCallback }, {}) {
    return (
      <header class={style.header}>
        <Link tabIndex="0" href="/">
          <h1>Hydrocarbon</h1>
        </Link>
        {this.renderUserNav(loggedIn, email, logoutCallback)}
      </header>
    );
  }
}

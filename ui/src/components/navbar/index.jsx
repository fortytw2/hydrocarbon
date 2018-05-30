import { h, Component } from "preact";
import { Link } from "preact-router/match";
import style from "./style.css";

export default class Header extends Component {
  render() {
    return (
      <header class={style.header}>
        <Link tabIndex="0" href="/">
          <h1>Hydrocarbon</h1>
        </Link>
        <nav>
          <Link tabIndex="0" activeClassName={style.active} href="/feed">
            Feed
          </Link>
          <Link tabIndex="0" activeClassName={style.active} href="/settings">
            Settings
          </Link>
          <Link tabIndex="0" activeClassName={style.active} href="/logout">
            Logout
          </Link>
        </nav>
      </header>
    );
  }
}

import { h, Component } from "preact";
import { route } from "preact-router";
import { bind } from "decko";
import { createKey } from "@/http";

export default class Callback extends Component {
  @bind
  async componentDidMount() {
    if (this.props.token) {
      try {
        const { data } = await createKey({ token: this.props.token });
        const { key, email } = data;

        this.props.loginCallback(email, key);

        route("/");
      } catch (e) {
        console.warn(e);
        route("/login");
      }
    }
  }

  render({}, {}) {
    return null;
  }
}

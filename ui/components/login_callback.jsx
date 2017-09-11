import { Component, h, render } from "preact";

import { ActivateLoginToken } from "../http/user";
import { NOTIFICATION_LEVEL_WARNING } from "../state/notifications/types";
import Redux from "preact-redux";
import { addNotification } from "../state/notifications/actions";
import { route } from "preact-router";

class LoginCallback extends Component {
  componentDidMount() {
    var usp = new URLSearchParams(window.location.search);
    if (!usp.has("token")) {
      this.props.dispatch(
        addNotification(NOTIFICATION_LEVEL_WARNING, "no token, try again")
      );
    } else {
      ActivateLoginToken(this.props.dispatch, usp.get("token"));
      route("/");
    }
  }
  render(props, state) {
    return <div />;
  }
}

const mapStateToProps = state => {
  return {
    login: state.login
  };
};

export default Redux.connect(mapStateToProps)(LoginCallback);

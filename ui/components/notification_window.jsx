import { Component, h } from "preact";

import Notification from "./notification";
import Redux from "preact-redux";

class NotificationWindow extends Component {
  render(props, state) {
    if (props.notifications.length === 0) {
      return null;
    }

    return (
      <div>
        {props.notifications.map(n => (
          <Notification
            dispatch={props.dispatch}
            sKey={n.key}
            message={n.message}
            level={n.level}
          />
        ))}
      </div>
    );
  }
}

const mapStateToProps = state => {
  return {
    notifications: state.notifications
  };
};

export default Redux.connect(mapStateToProps)(NotificationWindow);

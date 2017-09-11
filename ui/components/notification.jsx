import { Component, h } from "preact";

import { NOTIFICATION_LEVEL_WARNING } from "../state/notifications/types";
import { removeNotification } from "../state/notifications/actions";

class Notification extends Component {
  constructor(props) {
    super(props);
    this.handleClick = this.handleClick.bind(this);
  }
  handleClick(e) {
    e.preventDefault();
    this.props.dispatch(removeNotification(this.props.sKey));
  }
  render(props, state) {
    var style = props.level === NOTIFICATION_LEVEL_WARNING
      ? "flex items-center justify-center pa2 " + "bg-washed-red"
      : "flex items-center justify-center pa2 " + "bg-lightest-blue";
    return (
      <div class={style} onClick={this.handleClick}>
        <span class="lh-title ml1">{props.message}</span>
      </div>
    );
  }
}

export default Notification;

import { h, Component } from "preact";
import { bind } from "decko";

import style from "./style.css";

export default class Modal extends Component {
  @bind
  calculateStyle(open) {
    if (open) {
      return "visibility: visible; opacity: 1;";
    }

    return "";
  }

  render({ children, open, close }, {}) {
    return (
      <div
        class={style.modal}
        onClick={close}
        style={this.calculateStyle(open)}
      >
        <div class={style.modalContainer} onClick={e => e.stopPropagation()}>
          <div class={style.modalChildren}>{children}</div>
        </div>
      </div>
    );
  }
}

import m from "mithril";
import flash from "../models/flash";

export default {
  view(vnode) {
    if (flash.flashMessage !== "") {
      var msg = flash.flashMessage;
      flash.flashMessage = "";

      return m(
        "div",
        {
          onchange: function(e) {
            flash.flashMessage[e.target.name] = e.target.value;
          },
          class: "flex items-center justify-center pa4 bg-lightest-red red"
        },
        m("span", { class: "lh-title ml3" }, msg)
      );
    } else {
      return;
    }
  }
};

import m from "mithril";
import * as Mithril from "mithril";
import nav from "./nav";
import footer from "./footer";

export default {
  view(vnode) {
    return m("div", [
      m(nav),
      m("div", { class: "fl w-100 pa2 h-auto" }, [
        m("p", "welcome to hydrocarbon, an internet reader for dinosaurs")
      ]),
      m(footer)
    ]);
  }
};

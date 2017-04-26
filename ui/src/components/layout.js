import m from "mithril";
import nav from "./nav";
import footer from "./footer";
import flash from "./flash";

export default {
  view(vnode) {
    return m("div", {class: "min-vh-100"}, [
      m(nav),
      m(flash),
      m("div", { class: "fl w-100 pa2 h-auto" }, vnode.children),
      m(footer)
    ]);
  }
};

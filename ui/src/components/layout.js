import m from "mithril";
import nav from "./nav";
import footer from "./footer";

export default {
  view(vnode) {
    return m("div", {class: "min-vh-100"}, [
      m(nav),
      m("div", { class: "fl w-100 pa2 h-auto" }, vnode.children),
      m(footer)
    ]);
  }
};

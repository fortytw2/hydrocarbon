import m from "mithril";
import user from "../models/user";

export default {
  oncreate(vnode) {
    user.activateToken(m.route.param("key"));
  },
  view(vnode) {
    return m.route.set("/");
  }
};

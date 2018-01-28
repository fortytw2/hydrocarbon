import { h, Component } from "preact";
import Drawer from "preact-material-components/Drawer";
import List from "preact-material-components/List";
import "preact-material-components/Drawer/style.css";
import "preact-material-components/List/style.css";
import { route } from "preact-router";

import style from "./style";

export default class FolderView extends Component {
  getContent(feedID) {
    if (feedID === undefined) {
      return <h1> No Feed ID </h1>;
    }
    return <h2> feed content for {feedID} </h2>;
  }

  linkTo = path => () => {
    route(path);
  };

  render({ id, feedID }, {}) {
    return (
      <div class={style.content}>
        <Drawer.PermanentDrawer spacer={false}>
          <Drawer.PermanentDrawerContent>
            <List>
              <List.LinkItem
                onClick={this.linkTo("/folders/" + id + "/feedID")}
              >
                Ars Technical
              </List.LinkItem>
              <List.LinkItem>Ars Technical</List.LinkItem>
              <List.LinkItem>Ars Technical</List.LinkItem>
              <List.LinkItem>Ars Technical</List.LinkItem>
              <List.LinkItem>Ars Technical</List.LinkItem>
            </List>
          </Drawer.PermanentDrawerContent>
        </Drawer.PermanentDrawer>
        <div>{this.getContent(feedID)}</div>
      </div>
    );
  }
}

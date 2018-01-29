import { h, Component } from "preact";
import Drawer from "preact-material-components/Drawer";
import List from "preact-material-components/List";
import "preact-material-components/Drawer/style.css";
import "preact-material-components/List/style.css";
import Feed from "../feed";

import { route } from "preact-router";

import style from "./style";

export default class FolderView extends Component {
  getContent(feedID) {
    if (feedID === undefined) {
      return <h1> No Feed ID </h1>;
    }
    return <Feed id={feedID} />;
  }

  linkTo = path => () => {
    route(path);
  };

  render({ id, feedID, feeds }, {}) {
    if (feeds === undefined || feeds.length === 0) {
      return (
        <div class={style.content}>
          <h3>No Feeds Found, try adding one? </h3>
        </div>
      );
    }

    return (
      <div class={style.content}>
        <Drawer.PermanentDrawer spacer={false}>
          <Drawer.PermanentDrawerContent>
            <List>
              {feeds.map(f => {
                return (
                  <List.LinkItem
                    onClick={this.linkTo("/folders/" + id + "/" + f.id)}
                  >
                    {f.title}
                  </List.LinkItem>
                );
              })}
            </List>
          </Drawer.PermanentDrawerContent>
        </Drawer.PermanentDrawer>
        <div>{this.getContent(feedID)}</div>
      </div>
    );
  }
}

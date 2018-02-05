import { h, Component } from "preact";
import Drawer from "preact-material-components/Drawer";
import List from "preact-material-components/List";
import "preact-material-components/Drawer/style.css";
import "preact-material-components/List/style.css";
import Dialog from "preact-material-components/Dialog";
import "preact-material-components/Dialog/style.css";
import Button from "preact-material-components/Button";
import "preact-material-components/Button/style.css";
import PostList from "../postlist";

import { route } from "preact-router";

import style from "./style";

export default class FeedList extends Component {
  constructor(props) {
    super(props);
    this.setState({
      newFeedPlugin: "",
      newFeedURL: ""
    });
  }

  getContent(feedID) {
    if (feedID === undefined) {
      return (
        <div style="padding: 24px;">
          <h1> No Feed ID </h1>
        </div>
      );
    }
    return <PostList id={feedID} />;
  }

  linkTo = path => () => {
    route(path);
  };

  openWizard = () => {
    this.dialog.MDComponent.show();
  };

  updateUrl = e => {
    e.preventDefault();

    let url = e.target.value;
    this.setState({
      newFeedURL: url
    });
  };

  updatePlugin = e => {
    e.preventDefault();

    let plugin = e.target.value;
    this.setState({
      newFeedPlugin: plugin
    });
  };

  submitNewFeed = e => {
    e.preventDefault();

    console.log("lol can't make feeds rn");
  };

  dialogRef = dialog => (this.dialog = dialog);

  render({ id, feedID, feeds }, { newFeedPlugin, newFeedURL }) {
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
            <List.Item onClick={this.openWizard}>Add Feed to Folder</List.Item>
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
        <Dialog ref={this.dialogRef}>
          <Dialog.Header>Add Feed</Dialog.Header>
          <Dialog.Body>
            <div>
              <input
                type="text"
                placeholder="example url"
                value={newFeedURL}
                onChange={this.updateUrl}
              />
              <input
                type="text"
                placeholder="example plugin"
                value={newFeedPlugin}
                onChange={this.updatePlugin}
              />
            </div>
          </Dialog.Body>
          <Dialog.Footer>
            <Dialog.FooterButton accept onClick={this.submitNewFeed}>
              Create Feed
            </Dialog.FooterButton>
          </Dialog.Footer>
        </Dialog>
        <div>{this.getContent(feedID)}</div>
      </div>
    );
  }
}

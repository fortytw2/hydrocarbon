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

import styles from "./style";

export default class FeedList extends Component {
  constructor(props) {
    super(props);
    this.setState({
      loading: true,
      newFeedPlugin: "",
      newFeedURL: "",
      feeds: []
    });
  }

  componentDidMount(props) {
    this.componentWillReceiveProps();
  }

  componentWillReceiveProps(nextProps) {
    if (this.props.folderID === nextProps.folderID) {
      return;
    }

    this.setState({
      loading: true
    });
    let key = window.localStorage.getItem("hydrocarbon-key");

    fetch(window.baseURL + "/v1/feed/list", {
      method: "POST",
      headers: {
        "x-hydrocarbon-key": key
      },
      body: JSON.stringify({
        folder_id: nextProps.folderID
      })
    })
      .then(res => {
        if (res.ok) {
          return res.json();
        }
      })
      .then(json => {
        this.setState({
          loading: false,
          feeds: json
        });
      });
  }

  getContent(feedID) {
    if (typeof feedID === "undefined") {
      return (
        <div style="padding: 24px;">
          <h1> No Feed ID </h1>
        </div>
      );
    }
    return <PostList feedID={feedID} />;
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
    this.setState({ newFeedURL: url });
  };

  updatePlugin = e => {
    e.preventDefault();

    let plugin = e.target.value;
    this.setState({ newFeedPlugin: plugin });
  };

  submitNewFeed = e => {
    e.preventDefault();

    let key = window.localStorage.getItem("hydrocarbon-key");
    let fURL = this.state.newFeedURL;
    let fPlugin = this.state.newFeedPlugin;
    fetch(window.baseURL + "/v1/feed/create", {
      method: "POST",
      headers: {
        "x-hydrocarbon-key": key
      },
      body: JSON.stringify({
        url: fURL,
        plugin: fPlugin,
        folder_id: this.props.folderID
      })
    })
      .then(res => {
        if (res.ok) {
          return res.json();
        }
      })
      .then(json => {
        let f = this.state.feeds;
        f = f.concat({
          id: json.id,
          title: fURL,
          url: fURL,
          plugin: fPlugin
        });
        this.setState({
          feeds: f,
          newFeedPlugin: "",
          newFeedURL: ""
        });
      });
  };

  dialogRef = dialog => (this.dialog = dialog);

  render({ folderID, feedID }, { loading, feeds, newFeedPlugin, newFeedURL }) {
    if (feeds === undefined || feeds === null || feeds.length === 0) {
      feeds = [];
    }

    if (loading) {
      return (
        <div class={styles.content}>
          <Drawer.PermanentDrawer spacer={false}>
            <Drawer.PermanentDrawerContent>
              <List.Item onClick={this.openWizard}>
                Add Feed to Folder
              </List.Item>
              <List>
                <List.Item>Loading...</List.Item>
              </List>
            </Drawer.PermanentDrawerContent>
          </Drawer.PermanentDrawer>
        </div>
      );
    }

    return (
      <div class={styles.content}>
        <Drawer.PermanentDrawer spacer={false}>
          <Drawer.PermanentDrawerContent>
            <List.Item onClick={this.openWizard}>Add Feed to Folder</List.Item>
            <List>
              {feeds.map(f => {
                if (f.id === feedID) {
                  return (
                    <a
                      onClick={this.linkTo("/folders/" + folderID + "/" + f.id)}
                      class="mdc-list-item mdc-list-item--activated"
                    >
                      {f.title}
                    </a>
                  );
                }
                return (
                  <List.LinkItem
                    onClick={this.linkTo("/folders/" + folderID + "/" + f.id)}
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

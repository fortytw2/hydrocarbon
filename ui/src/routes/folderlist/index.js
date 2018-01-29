import { h, Component } from "preact";
import Drawer from "preact-material-components/Drawer";
import List from "preact-material-components/List";
import "preact-material-components/Drawer/style.css";
import "preact-material-components/List/style.css";
import styles from "./style";
import { route } from "preact-router";

import FolderView from "../folderview";

export default class FolderList extends Component {
  constructor(props) {
    super(props);

    this.setState({
      loading: true,
      folders: []
    });
  }

  componentDidMount() {
    let key = window.localStorage.getItem("hydrocarbon-key");

    fetch("/v1/folder/list", {
      method: "POST",
      headers: {
        "x-hydrocarbon-key": key
      }
    })
      .then(res => {
        if (res.ok) {
          return res.json();
        }
      })
      .then(json => {
        this.setState({ loading: false, folders: json });
      });
  }

  linkTo = path => () => {
    route(path);
  };

  goToFolder(id) {
    return this.linkTo("/folders/" + id);
  }

  getFeedsForFolder = id => {
    let feeds = this.state.folders.map(e => {
      if (e.id === id) {
        return e.feeds;
      }
    });
    return feeds[0];
  };

  render({ id, feedID }, { loading, folders }) {
    if (loading) {
      return <div class={styles.content}>Loading Folders...</div>;
    }

    return (
      <div class={styles.content}>
        <Drawer.PermanentDrawer spacer={false}>
          <Drawer.PermanentDrawerContent>
            <List>
              {folders.map(f => {
                return (
                  <List.LinkItem onClick={this.goToFolder(f.id)}>
                    {f.title} {f.unread}
                  </List.LinkItem>
                );
              })}
            </List>
          </Drawer.PermanentDrawerContent>
        </Drawer.PermanentDrawer>
        <div>
          <FolderView
            id={id}
            feeds={this.getFeedsForFolder(id)}
            feedID={feedID}
          />
        </div>
      </div>
    );
  }
}

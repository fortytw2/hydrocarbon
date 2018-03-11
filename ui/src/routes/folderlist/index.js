import { h, Component } from "preact";
import Drawer from "preact-material-components/Drawer";
import List from "preact-material-components/List";
import Dialog from "preact-material-components/Dialog";
import "preact-material-components/Drawer/style.css";
import "preact-material-components/List/style.css";
import "preact-material-components/Dialog/style.css";
import Button from "preact-material-components/Button";
import "preact-material-components/Button/style.css";
import styles from "./style";
import { route } from "preact-router";

import FeedList from "../feedlist";

export default class FolderList extends Component {
  constructor(props) {
    super(props);

    this.setState({ loading: true, folders: [], newFolderName: "" });
  }

  componentDidMount() {
    this.loadFolders();
  }

  loadFolders = () => {
    let key = window.localStorage.getItem("hydrocarbon-key");

    fetch(window.baseURL + "/v1/folder/list", {
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
  };

  linkTo = path => () => {
    route(path);
  };

  goToFolder(id) {
    return this.linkTo("/folders/" + id);
  }

  openWizard = () => {
    this.dialog.MDComponent.show();
  };

  updateNewFolder = e => {
    e.preventDefault();

    let newFolderName = e.target.value;
    this.setState({ newFolderName });
  };

  submitNewFolder = e => {
    e.preventDefault();

    let key = window.localStorage.getItem("hydrocarbon-key");
    let fName = this.state.newFolderName;
    fetch(window.baseURL + "/v1/folder/create", {
      method: "POST",
      headers: {
        "x-hydrocarbon-key": key
      },
      body: JSON.stringify({
        name: fName
      })
    })
      .then(res => {
        if (res.ok) {
          return res.json();
        }
      })
      .then(json => {
        let f = this.state.folders;
        f = f.concat({ id: json.id, title: fName });
        this.setState({
          folders: f,
          newFolderName: ""
        });
      });
  };

  dialogRef = dialog => (this.dialog = dialog);

  getList(folderID, feedID) {
    if (!(folderID.length >= 2)) {
      return <div>No folder id</div>;
    }
    return <FeedList folderID={folderID} feedID={feedID} />;
  }

  render({ id, feedID }, { loading, folders, newFolderName }) {
    if (loading) {
      return <div class={styles.content}>Loading Folders...</div>;
    }

    return (
      <div class={styles.content}>
        <Drawer.PermanentDrawer spacer={false}>
          <Drawer.PermanentDrawerContent>
            <List>
              <List.Item onClick={this.openWizard}>Add Folder</List.Item>
              {folders.map(f => {
                if (f.id === id) {
                  return (
                    <a
                      class="mdc-list-item mdc-list-item--activated"
                      onClick={this.goToFolder(f.id)}
                    >
                      {f.title} {f.unread}
                    </a>
                  );
                }
                return (
                  <List.LinkItem onClick={this.goToFolder(f.id)}>
                    {f.title} {f.unread}
                  </List.LinkItem>
                );
              })}
            </List>
          </Drawer.PermanentDrawerContent>
        </Drawer.PermanentDrawer>
        <Dialog ref={this.dialogRef}>
          <Dialog.Header>Add Folder</Dialog.Header>
          <Dialog.Body>
            <div>
              <input
                type="text"
                placeholder="example folder"
                value={newFolderName}
                onChange={this.updateNewFolder}
              />
            </div>
          </Dialog.Body>
          <Dialog.Footer>
            <Dialog.FooterButton accept onClick={this.submitNewFolder}>
              Create Folder
            </Dialog.FooterButton>
          </Dialog.Footer>
        </Dialog>
        <div>{this.getList(id, feedID)}</div>
      </div>
    );
  }
}

import { h, Component } from "preact";
import { route } from "preact-router";
import Toolbar from "preact-material-components/Toolbar";
import Drawer from "preact-material-components/Drawer";
import List from "preact-material-components/List";
import Dialog from "preact-material-components/Dialog";
import Switch from "preact-material-components/Switch";
import "preact-material-components/Switch/style.css";
import "preact-material-components/Dialog/style.css";
import "preact-material-components/Drawer/style.css";
import "preact-material-components/List/style.css";
import "preact-material-components/Toolbar/style.css";
import style from "./style";

export default class Layout extends Component {
  closeDrawer() {
    this.state = {
      darkThemeEnabled: false
    };
  }

  openSettings = () => this.dialog.MDComponent.show();

  drawerRef = drawer => (this.drawer = drawer);
  dialogRef = dialog => (this.dialog = dialog);

  toggleDarkTheme = () => {
    this.setState(
      {
        darkThemeEnabled: !this.state.darkThemeEnabled
      },
      () => {
        if (this.state.darkThemeEnabled) {
          document.body.classList.add("mdc-theme--dark");
        } else {
          document.body.classList.remove("mdc-theme--dark");
        }
      }
    );
  };

  feedSections(loggedIn) {
    if (loggedIn) {
      return (
        <Toolbar.Section align-start>
          <Toolbar.Title style="font-size: 14px;">
            <a class={style.toolbarlink} href="/folders">
              Feed
            </a>
          </Toolbar.Title>
        </Toolbar.Section>
      );
    } else {
      return <div />;
    }
  }

  loginSection(loggedIn) {
    if (loggedIn) {
      return (
        <Toolbar.Section align-end>
          <Toolbar.Title style="font-size: 14px;">
            <a class={style.toolbarlink} href="/logout">
              Logout
            </a>
          </Toolbar.Title>
        </Toolbar.Section>
      );
    } else {
      return (
        <Toolbar.Section align-end>
          <Toolbar.Title style="font-size: 14px;">
            <a class={style.toolbarlink} href="/login">
              Login
            </a>
          </Toolbar.Title>
        </Toolbar.Section>
      );
    }
  }

  render({ loggedIn }, {}) {
    return (
      <div class={style.layout}>
        <Toolbar className="toolbar">
          <Toolbar.Row>
            <Toolbar.Section align-start style="flex-grow: 0.35;">
              <Toolbar.Title>
                <a class={style.toolbarlink} href="/">
                  Hydrocarbon
                </a>
              </Toolbar.Title>
            </Toolbar.Section>
            {this.feedSections(loggedIn)}
            {this.loginSection(loggedIn)}
          </Toolbar.Row>
        </Toolbar>
        <div class={style.content}>{this.props.children}</div>
        <Dialog ref={this.dialogRef}>
          <Dialog.Header>Settings</Dialog.Header>
          <Dialog.Body>
            <div>
              Enable dark theme <Switch onClick={this.toggleDarkTheme} />
            </div>
          </Dialog.Body>
          <Dialog.Footer>
            <Dialog.FooterButton accept>okay</Dialog.FooterButton>
          </Dialog.Footer>
        </Dialog>
      </div>
    );
  }
}

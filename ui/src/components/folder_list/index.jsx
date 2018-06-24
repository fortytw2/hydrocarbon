import { h, Component } from "preact";
import { bind } from "decko";
import { createFeed, listPlugins } from "@/http";
import { Link, route } from "preact-router";

import style from "./style.css";
import textboxStyle from "@/styles/textbox.css";

import LibraryAddIcon from "@/assets/library_add.svg";
import RemoveCircleIcon from "@/assets/remove_circle.svg";
import AddCircleIcon from "@/assets/add_circle.svg";

const initialState = {
  collapsedFolders: {}
};

export default class FolderList extends Component {
  constructor(props) {
    super(props);
    this.setState(initialState);
  }

  @bind
  folderLink(folderId) {
    return `/feed/${folderId}`;
  }

  @bind
  feedLink(folderId, feedId) {
    return `/feed/${folderId}/${feedId}`;
  }

  @bind
  addFeed(folderId) {
    return e => {
      this.props.openNewFeedModal(e, folderId);
    };
  }

  @bind
  toggleFolderCollapse(folderId) {
    return e => {
      e.preventDefault();
      if (this.state.collapsedFolders[folderId]) {
        this.setState({
          collapsedFolders: {
            ...this.state.collapsedFolders,
            [folderId]: false
          }
        });
      } else {
        this.setState({
          collapsedFolders: {
            ...this.state.collapsedFolders,
            [folderId]: true
          }
        });
      }
    };
  }

  @bind
  renderFolder(folderId, folder, feedId) {
    let collapseIcon = RemoveCircleIcon;
    let subListClass = style.folderSubList;
    if (this.state.collapsedFolders[folder.id]) {
      subListClass = [subListClass, style.collapsedFolder].join(" ");
      collapseIcon = AddCircleIcon;
    }

    return (
      <li class={style.folder}>
        <span class={style.folderTitle}>
          <span class={style.folderSubTitle}>
            <a onClick={this.toggleFolderCollapse(folder.id)}>
              <img class={style.icon} src={collapseIcon} />
            </a>
            {folder.title}
          </span>
          <a onClick={this.addFeed(folder.id)}>
            <img class={style.icon} src={LibraryAddIcon} />
          </a>
        </span>
        <ol class={subListClass}>
          {this.renderFeeds(folder.id, feedId, folder.feeds)}
        </ol>
      </li>
    );
  }

  @bind
  renderFeeds(folderId, feedId, feeds) {
    if (feeds === null || feeds.length === 0) {
      return;
    }
    return feeds.map(f => this.renderFeed(folderId, feedId, f));
  }

  @bind
  renderFeed(folderId, feedId, feed) {
    let listClass = style.feed;
    if (feedId === feed.id) {
      listClass = [style.feed, style.activeFeed].join(" ");
    }

    return (
      <li class={listClass}>
        <a href={this.feedLink(folderId, feed.id)}>{feed.title}</a>
      </li>
    );
  }

  render({ folderId, folders, feedId }, { collapsedFolders }) {
    return (
      <div>
        <div class={style.searchBox}>
          <input class={textboxStyle.input} type="text" placeholder="filter" />
        </div>
        <ol class={style.folderListInside}>
          {folders.map(f => this.renderFolder(folderId, f, feedId))}
        </ol>
      </div>
    );
  }
}

import { h, Component } from "preact";
import { bind } from "decko";
import { createFeed, listPlugins } from "@/http";
import { Link } from "preact-router";

import style from "./style.css";
import textboxStyle from "@/styles/textbox.css";

export default class FolderList extends Component {
  @bind
  folderLink(folderId) {
    return `/feed/${folderId}`;
  }

  @bind
  feedLink(folderId, feedId) {
    return `/feed/${folderId}/${feedId}`;
  }

  @bind
  renderFolder(folder, feedId) {
    return (
      <li class={style.folder}>
        <Link
          tabIndex="0"
          activeClassName={style.activeLink}
          href={this.folderLink(folder.id)}
        >
          {folder.title}
        </Link>
        <ol class={style.folderSubList}>
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

  render({ folderId, folders, feedId }, {}) {
    return (
      <div>
        <div class={style.searchBox}>
          <input class={textboxStyle.input} type="text" placeholder="filter" />
        </div>
        <ol class={style.folderListInside}>
          {folders.map(f => this.renderFolder(f, feedId))}
        </ol>
      </div>
    );
  }
}

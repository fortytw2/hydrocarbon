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
  renderFolder(folder) {
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
          {this.renderFeeds(folder.id, folder.feeds)}
        </ol>
      </li>
    );
  }

  @bind
  renderFeeds(folderId, feeds) {
    if (feeds === null || feeds.length === 0) {
      return;
    }
    return feeds.map(f => this.renderFeed(folderId, f));
  }

  @bind
  renderFeed(folderId, feed) {
    return (
      <li class={style.feed}>
        <Link
          tabIndex="0"
          activeClassName={style.activeLink}
          href={this.feedLink(folderId, feed.id)}
        >
          {feed.title}
        </Link>
      </li>
    );
  }

  render({ folderId, folders }, {}) {
    return (
      <div>
        <div class={style.searchBox}>
          <input class={textboxStyle.input} type="text" placeholder="filter" />
        </div>
        <ol class={style.folderListInside}>
          {folders.map(f => this.renderFolder(f))}
        </ol>
      </div>
    );
  }
}

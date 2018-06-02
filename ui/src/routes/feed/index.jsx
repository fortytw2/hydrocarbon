import { h, Component } from "preact";
import { Link } from "preact-router";
import style from "./style.css";
import textboxStyle from "@/styles/textbox.css";

export default class Feed extends Component {
  render({ apiKey }, {}) {
    return (
      <div class={style.feedContainer}>
        <div class={style.folderList}>
          <div class={style.editBox}>
            <button class={style.editButton}>New Folder</button>
            <button class={style.editButton}>Edit</button>
          </div>
          <div class={style.searchBox}>
            <input
              class={textboxStyle.input}
              type="text"
              placeholder="filter"
            />
          </div>
          <ol class={style.folderListInside}>
            <li class={style.folder}>
              <Link tabIndex="0" activeClassName={style.active} href="/feed/1">
                Folder 1
              </Link>
              <ol class={style.folderSubList}>
                <li class={style.feed}>
                  <Link
                    tabIndex="0"
                    activeClassName={style.active}
                    href="/feed/1/1"
                  >
                    Feed 1
                  </Link>
                </li>
                <li class={style.feed}>
                  <Link
                    tabIndex="0"
                    activeClassName={style.active}
                    href="/feed/1/2"
                  >
                    Feed 2
                  </Link>
                </li>
                <li class={style.feed}>
                  <Link
                    tabIndex="0"
                    activeClassName={style.active}
                    href="/feed/1/3"
                  >
                    Feed 3
                  </Link>
                </li>
                <li class={style.feed}>
                  <Link
                    tabIndex="0"
                    activeClassName={style.active}
                    href="/feed/1/4"
                  >
                    Feed 4
                  </Link>
                </li>
                <li class={style.feed}>
                  <Link
                    tabIndex="0"
                    activeClassName={style.active}
                    href="/feed/1/5"
                  >
                    Feed 5
                  </Link>
                </li>
                <li class={style.feed}>
                  <Link
                    tabIndex="0"
                    activeClassName={style.active}
                    href="/feed/1/6"
                  >
                    Feed 6
                  </Link>
                </li>
              </ol>
            </li>
            <li class={style.folder}>Folder 2</li>
            <li class={style.folder}>Folder 3</li>
            <li class={style.folder}>Folder 4</li>
            <li class={style.folder}>Folder 5</li>
            <li class={style.folder}>Folder 6</li>
          </ol>
        </div>
        <div class={style.postList}>
          <ol>
            <li>Post 1</li>
            <li>Post 2</li>
            <li>Post 3</li>
            <li>Post 4</li>
            <li>Post 5</li>
            <li>Post 6</li>
          </ol>
        </div>
        <div class={style.postContent}>
          <p>BLAH BLAH BLAH TEXT</p>
        </div>
      </div>
    );
  }
}

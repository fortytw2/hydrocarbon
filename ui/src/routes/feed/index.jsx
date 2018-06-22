import { h, Component } from "preact";
import { Link } from "preact-router";
import { bind } from "decko";

import Modal from "@/components/modal";
import FolderList from "@/components/folder_list";
import CreateFolderForm from "@/components/create_folder_form";
import CreateFeedForm from "@/components/create_feed_form";

import { listFolders, listFeeds } from "@/http";

import style from "./style.css";
import textboxStyle from "@/styles/textbox.css";

const initialState = {
  newFeedModal: false,
  newFolderModal: false,

  folders: []
};

export default class Feed extends Component {
  constructor(props) {
    super(props);

    this.setState(initialState);
  }

  async componentDidMount(props) {
    try {
      let folders = await listFolders({ apiKey: this.props.apiKey });

      folders = await Promise.all(
        folders.map(async f => {
          const feeds = await listFeeds({
            apiKey: this.props.apiKey,
            folderId: f.id
          });
          return {
            ...f,
            feeds
          };
        })
      );

      this.setState({ folders: folders });
    } catch (e) {
      console.log(e);
    }
  }

  @bind
  openNewFeedModal(e) {
    e.preventDefault();
    this.setState({ newFeedModal: true });
  }

  @bind
  closeNewFeedModal(e) {
    e.preventDefault();
    this.setState({ newFeedModal: false });
  }

  @bind
  submitFeed({ name, id }) {
    this.setState({ newFeedModal: false });
    // TODO: add feed to lcoa list
  }

  @bind
  openNewFolderModal(e) {
    e.preventDefault();
    this.setState({ newFolderModal: true });
  }

  @bind
  closeNewFolderModal(e) {
    e.preventDefault();
    this.setState({ newFolderModal: false });
  }

  @bind
  submitFolder({ name, id }) {
    this.setState({ newFolderModal: false });
    // TODO: add folder to local list
  }

  render({ apiKey, folderId }, { newFeedModal, newFolderModal, folders }) {
    return (
      <div class={style.feedContainer}>
        <Modal close={this.closeNewFeedModal} open={newFeedModal}>
          <CreateFeedForm
            onSubmit={this.submitFeed}
            apiKey={apiKey}
            folderId={folderId}
          />
        </Modal>

        <Modal close={this.closeNewFolderModal} open={newFolderModal}>
          <CreateFolderForm onSubmit={this.submitFolder} apiKey={apiKey} />
        </Modal>

        <div class={style.folderList}>
          <div class={style.editBox}>
            <button class={style.editButton} onClick={this.openNewFeedModal}>
              +Feed
            </button>
            <button class={style.editButton} onClick={this.openNewFolderModal}>
              +Folder
            </button>
          </div>
          <FolderList folders={folders} folderId={folderId} />
        </div>
        <div class={style.postList}>
          <ol class={style.postListInside}>
            <Link
              tabIndex="0"
              activeClassName={style.activeLink}
              href="/feed/1/6/1"
            >
              <li class={style.post}>
                <span class={style.postTitle}>001 Good Morning Brother</span>
                <span class={style.postTime}>14:06</span>
              </li>
            </Link>
            <li class={style.unreadPost}>
              <Link
                tabIndex="0"
                activeClassName={style.activeLink}
                href="/feed/1/6/2"
              >
                *Post 2
              </Link>
            </li>
            <li class={style.post}>
              <Link
                tabIndex="0"
                activeClassName={style.activeLink}
                href="/feed/1/6/3"
              >
                Post 3
              </Link>
            </li>
            <li class={style.post}>
              <Link
                tabIndex="0"
                activeClassName={style.activeLink}
                href="/feed/1/6/4"
              >
                Post 4
              </Link>
            </li>
            <li class={style.post}>
              <Link
                tabIndex="0"
                activeClassName={style.activeLink}
                href="/feed/1/6/5"
              >
                Post 5
              </Link>
            </li>
            <li class={style.post}>
              <Link
                tabIndex="0"
                activeClassName={style.activeLink}
                href="/feed/1/6/6"
              >
                Post 6
              </Link>
            </li>
          </ol>
        </div>
        <div class={style.postContent}>
          <div class={style.postInnerContent}>
            <div class={style.postHeader}>
              <h2 class={style.postTitleHeader}>
                Mother of Learning - nobody103
              </h2>
            </div>
            <div class={style.postSubHeader}>
              <h4>001 Good Morning Brother</h4>
              <h4>14:06 Monday, Feb 4, 2018</h4>
            </div>
            <div class={style.postBody} />
          </div>
        </div>
      </div>
    );
  }
}

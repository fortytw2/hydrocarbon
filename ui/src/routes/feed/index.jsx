import { h, Component } from "preact";
import { Link } from "preact-router";
import { bind } from "decko";

import Modal from "@/components/modal";
import FolderList from "@/components/folder_list";
import PostList from "@/components/post_list";
import CreateFolderForm from "@/components/create_folder_form";
import CreateFeedForm from "@/components/create_feed_form";

import { listFolders, listFeeds } from "@/http";

import style from "./style.css";
import textboxStyle from "@/styles/textbox.css";

import LibraryAddIcon from "@/assets/library_add.svg";

const initialState = {
  newFeedModal: false,
  newFeedModalFolderId: null,
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

      // TODO(fortytw2): push this down to the API layer
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
  openNewFeedModal(e, folderId) {
    e.preventDefault();
    this.setState({ newFeedModal: true, newFeedModalFolderId: folderId });
  }

  @bind
  closeNewFeedModal(e) {
    e.preventDefault();
    this.setState({ newFeedModal: false });
  }

  @bind
  submitFeed({ title, id, folderId }) {
    const lastf = this.state.folders.find(f => f.id === folderId);

    const f = lastf;
    if (f.feeds === null) {
      f.feeds = [{ title, id }];
    } else {
      f.feeds.push({ title, id });
    }

    const oldFolders = this.state.folders.filter(f => f.id !== folderId);
    this.setState({ newFeedModal: false, folders: [...oldFolders, f] });
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
  submitFolder({ title, id }) {
    const folders = [{ id, title, feeds: [] }, ...this.state.folders];

    this.setState({ newFolderModal: false, folders: folders });
  }

  render(
    { apiKey, folderId, feedId, postId },
    { newFeedModal, newFeedModalFolderId, newFolderModal, folders }
  ) {
    return (
      <div class={style.feedContainer}>
        <Modal close={this.closeNewFeedModal} open={newFeedModal}>
          <CreateFeedForm
            onSubmit={this.submitFeed}
            apiKey={apiKey}
            folderId={newFeedModalFolderId}
          />
        </Modal>

        <Modal close={this.closeNewFolderModal} open={newFolderModal}>
          <CreateFolderForm onSubmit={this.submitFolder} apiKey={apiKey} />
        </Modal>

        <div class={style.folderList}>
          <div class={style.editBox}>
            <button class={style.editButton} onClick={this.openNewFolderModal}>
              <span class={style.editText}>Add Folder</span>
              <img class={style.icon} src={LibraryAddIcon} />
            </button>
          </div>
          <FolderList
            folders={folders}
            folderId={folderId}
            feedId={feedId}
            openNewFeedModal={this.openNewFeedModal}
          />
        </div>

        <PostList
          apiKey={apiKey}
          folderId={folderId}
          feedId={feedId}
          postId={postId}
        />
      </div>
    );
  }
}

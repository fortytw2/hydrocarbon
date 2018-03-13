export const listFeeds = ({ folderID }) => {
  let key = window.localStorage.getItem("hydrocarbon-key");

  return fetch(window.baseURL + "/v1/feed/list", {
    method: "POST",
    headers: {
      "x-hydrocarbon-key": key
    },
    body: JSON.stringify({
      folder_id: folderID
    })
  }).then(res => {
    if (res.ok) {
      return res.json();
    }
  });
};

export const createFeed = ({ url, plugin, folderID }) => {
  let key = window.localStorage.getItem("hydrocarbon-key");

  return fetch(window.baseURL + "/v1/feed/create", {
    method: "POST",
    headers: {
      "x-hydrocarbon-key": key
    },
    body: JSON.stringify({
      url: url,
      plugin: plugin,
      folder_id: folderID
    })
  }).then(res => {
    if (res.ok) {
      return res.json();
    }
  });
};

export const listFolders = () => {
  let key = window.localStorage.getItem("hydrocarbon-key");

  return fetch(window.baseURL + "/v1/folder/list", {
    method: "POST",
    headers: {
      "x-hydrocarbon-key": key
    }
  }).then(res => {
    if (res.ok) {
      return res.json();
    }
  });
};

export const createFolder = ({ name }) => {
  let key = window.localStorage.getItem("hydrocarbon-key");

  return fetch(window.baseURL + "/v1/folder/create", {
    method: "POST",
    headers: {
      "x-hydrocarbon-key": key
    },
    body: JSON.stringify({
      name: name
    })
  }).then(res => {
    if (res.ok) {
      return res.json();
    }
  });
};

export const listPosts = ({ feedID }) => {
  let key = window.localStorage.getItem("hydrocarbon-key");

  return fetch(window.baseURL + "/v1/post/list", {
    method: "POST",
    headers: {
      "x-hydrocarbon-key": key
    },
    body: JSON.stringify({
      feed_id: feedID
    })
  });
};

export const createKey = ({ token }) => {
  return fetch(window.baseURL + "/v1/key/create", {
    method: "POST",
    body: JSON.stringify({
      token: token
    })
  }).then(response => {
    if (response.ok) {
      return response.json();
    }
  });
};

export const requestToken = ({ email }) => {
  return fetch(window.baseURL + "/v1/token/create", {
    method: "POST",
    body: JSON.stringify({
      email: email
    })
  }).then(response => {
    if (response.ok) {
      return response.json();
    }
  });
};

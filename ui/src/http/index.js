export const listFeeds = ({ folderId, apiKey }) => {
  return fetch("/v1/feed/list", {
    method: "POST",
    headers: {
      "x-hydrocarbon-key": apiKey
    },
    body: JSON.stringify({
      folder_id: folderId
    })
  })
    .then(res => {
      if (res.ok) {
        return res.json();
      }
    })
    .then(json => {
      if (json.status === "error") {
        throw json.error;
      }
      return json.data;
    });
};

export const createFeed = ({ url, folderId, apiKey }) => {
  return fetch("/v1/feed/create", {
    method: "POST",
    headers: {
      "x-hydrocarbon-key": apiKey
    },
    body: JSON.stringify({
      url: url,
      folder_id: folderId
    })
  })
    .then(res => {
      if (res.ok) {
        return res.json();
      }
    })
    .then(json => {
      if (json.status === "error") {
        throw json.error;
      }
      return json.data;
    });
};

export const listFolders = ({ apiKey }) => {
  return fetch("/v1/folder/list", {
    method: "POST",
    headers: {
      "x-hydrocarbon-key": apiKey
    }
  })
    .then(res => {
      if (res.ok) {
        return res.json();
      }
    })
    .then(json => {
      if (json.status === "error") {
        throw json.error;
      }
      return json.data;
    });
};

export const listPlugins = ({ apiKey }) => {
  return fetch("/v1/plugin/list", {
    method: "POST",
    headers: {
      "x-hydrocarbon-key": apiKey
    }
  })
    .then(res => {
      if (res.ok) {
        return res.json();
      }
    })
    .then(json => {
      if (json.status === "error") {
        throw json.error;
      }
      return json.data.plugins;
    });
};

export const markRead = ({ apiKey, postId }) => {
  return fetch("/v1/post/read", {
    body: JSON.stringify({
      post_id: postId
    }),
    method: "POST",
    headers: {
      "x-hydrocarbon-key": apiKey
    }
  })
    .then(res => {
      if (res.ok) {
        return res.json();
      }
    })
    .then(json => {
      if (json.status === "error") {
        throw json.error;
      }
      return true;
    });
};

export const createFolder = ({ name, apiKey }) => {
  return fetch("/v1/folder/create", {
    method: "POST",
    headers: {
      "x-hydrocarbon-key": apiKey
    },
    body: JSON.stringify({
      name: name
    })
  })
    .then(res => {
      if (res.ok) {
        return res.json();
      }
    })
    .then(json => {
      if (json.status === "error") {
        throw json.error;
      }
      return json.data;
    });
};

export const listPosts = ({ feedId, apiKey }) => {
  return fetch("/v1/post/list", {
    method: "POST",
    headers: {
      "x-hydrocarbon-key": apiKey
    },
    body: JSON.stringify({
      feed_id: feedId
    })
  })
    .then(response => {
      if (response.ok) {
        return response.json();
      }
    })
    .then(json => {
      if (json.status === "error") {
        throw json.error;
      }
      return json.data;
    });
};

export const createKey = ({ token }) => {
  return fetch("/v1/key/create", {
    method: "POST",
    body: JSON.stringify({
      token: token
    })
  })
    .then(response => {
      if (response.ok) {
        return response.json();
      }
    })
    .then(json => {
      if (json.status === "error") {
        throw json.error;
      }
      return json;
    });
};

export const verifyKey = apiKey => {
  return fetch("/v1/key/verify", {
    method: "POST",
    headers: {
      "x-hydrocarbon-key": apiKey
    }
  })
    .then(response => {
      if (response.ok) {
        return response.json();
      }
    })
    .then(json => {
      if (json.status === "error") {
        throw json.error;
      }
      return json;
    });
};

export const requestToken = ({ email }) => {
  return fetch("/v1/token/create", {
    method: "POST",
    body: JSON.stringify({
      email: email
    })
  })
    .then(response => {
      if (response.ok) {
        return response.json();
      }
    })
    .then(json => {
      if (json.status === "error") {
        throw json.error;
      }
      return json;
    });
};

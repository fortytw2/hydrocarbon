export const listFeeds = ({ folderID }) => {
  return fetch("/v1/feed/list", {
    method: "POST",
    headers: {
      "x-hydrocarbon-key": key
    },
    body: JSON.stringify({
      folder_id: folderID
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
      return json;
    });
};

export const createFeed = ({ url, plugin, folderID }) => {
  return fetch("/v1/feed/create", {
    method: "POST",
    headers: {
      "x-hydrocarbon-key": key
    },
    body: JSON.stringify({
      url: url,
      plugin: plugin,
      folder_id: folderID
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
      return json;
    });
};

export const listFolders = () => {
  return fetch("/v1/folder/list", {
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
      if (json.status === "error") {
        throw json.error;
      }
      return json;
    });
};

export const listPlugins = () => {
  return fetch("/v1/plugin/list", {
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
      if (json.status === "error") {
        throw json.error;
      }
      return json;
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
      return json;
    });
};

export const listPosts = ({ feedID }) => {
  return fetch("/v1/post/list", {
    method: "POST",
    headers: {
      "x-hydrocarbon-key": key
    },
    body: JSON.stringify({
      feed_id: feedID
    })
  }).then(json => {
    if (json.status === "error") {
      throw json.error;
    }
    return json;
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

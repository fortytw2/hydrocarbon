import {
  NOTIFICATION_LEVEL_INFO,
  NOTIFICATION_LEVEL_WARNING,
} from "../state/notifications/types";

import { Store } from "../state/store";
import { activateApiKey } from "../state/login/actions";
import { addNotification } from "../state/notifications/actions";
import { ajax } from "nanoajax";

export function RequestLoginToken(email) {
  ajax(
    {
      url: "/api/token/request",
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        email: email
      })
    },
    function(code, responseText, request) {
      if (code === 0) {
        Store.dispatch(
          addNotification(
            NOTIFICATION_LEVEL_INFO,
            "unable to connect to internet"
          )
        );

        return;
      }
      if (code !== 200) {
        alert("something terrible has happened", responseText);
        return;
      }

      var parsed = JSON.parse(responseText);
      if (parsed.status === "error") {
        Store.dispatch(
          addNotification(NOTIFICATION_LEVEL_WARNING, parsed.error)
        );
      } else if (parsed.status === "success") {
        Store.dispatch(addNotification(NOTIFICATION_LEVEL_INFO, parsed.note));
      }
    }
  );
}

export function ActivateLoginToken(dispatch, token) {
  ajax(
    {
      url: "/api/token/activate",
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        token: token
      })
    },
    function(code, responseText, request) {
      if (code === 0) {
        dispatch(
          addNotification(
            NOTIFICATION_LEVEL_INFO,
            "unable to connect to internet"
          )
        );

        return;
      }
      if (code !== 200) {
        alert("something terrible has happened", responseText);
        return;
      }

      var parsed = JSON.parse(responseText);
      if (parsed.status === "error") {
        dispatch(addNotification(NOTIFICATION_LEVEL_WARNING, parsed.error));
      } else if (parsed.status === "success") {
        dispatch(activateApiKey(parsed.email, parsed.key));
      }
    }
  );
}

const url = {
  instruments: "/rpc/instruments",
  candles: "/rpc/candles",
  signals: "/rpc/signals",
  subscribe: "/rpc/auth/subscribe",
  subscriptions: "/rpc/auth/subscriptions",
};

function removeSubscription({ idToken, instrumentID }) {
  return fetch(`${url.subscriptions}/${instrumentID}`, {
    credentials: "include",
    mode: "cors",
    headers: {
      Authorization: `Bearer ${idToken}`
    },
    method: "DELETE",
  })
    .then((response) => {
      if (response.ok) {
        return response;
      }
      throw new Error(`${response.statusText} : ${response.status}`);
    });
}

function retrieveSubscriptions(idToken) {
  return fetch(url.subscriptions, {
    credentials: "include",
    mode: "cors",
    headers: {
      Authorization: `Bearer ${idToken}`
    },
  })
    .then((response) => {
      if (response.ok) {
        return response;
      }
      throw new Error(`${response.statusText} : ${response.status}`);
    })
    .then(resp => resp.json());
}

function retrieveInstruments() {
  return fetch(url.instruments)
    .then((response) => {
      if (response.ok) {
        return response;
      }
      throw new Error(`${response.statusText} : ${response.status}`);
    })
    .then(resp => resp.json());
}

function retrieveCandles(instrumentID) {
  return fetch(`${url.candles}/${instrumentID}`)
    .then((response) => {
      if (response.ok) {
        return response;
      }
      throw new Error(`${response.statusText} : ${response.status}`);
    })
    .then(resp => resp.json());
}

function retrieveSignals(instrumentID, period, buyLimit, sellLimit) {
  return fetch(`${url.signals}/${instrumentID}/period/${period}/buyLimit/${buyLimit}/sellLimit/${sellLimit}`)
    .then((response) => {
      if (response.ok) {
        return response;
      }
      throw new Error(`${response.statusText} : ${response.status}`);
    })
    .then(resp => resp.json());
}

function addSubscription({ idToken, instrument, period, buyLimit, sellLimit }) {
  return fetch(`${url.subscribe}/period/${period}/buyLimit/${buyLimit}/sellLimit/${sellLimit}`, {
    method: "POST",
    body: JSON.stringify(instrument),
    credentials: "include",
    mode: "cors",
    headers: {
      Authorization: `Bearer ${idToken}`
    },
  })
    .then((response) => {
      if (response.ok) {
        return response;
      }
      throw new Error(`${response.statusText} : ${response.status}`);
    });
}

export {
  retrieveInstruments,
  retrieveCandles,
  retrieveSignals,
  addSubscription,
  retrieveSubscriptions,
  removeSubscription
};
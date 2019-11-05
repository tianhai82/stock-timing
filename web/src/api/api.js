const url = {
  instruments: "/rpc/instruments",
  candles: "/rpc/candles",
  signals: "/rpc/signals",
  subscribe: "/rpc/auth/subscribe",
};

function retrieveInstruments() {
  return fetch(url.instruments)
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

function retrieveSignals(instrumentID) {
  return fetch(`${url.signals}/${instrumentID}`)
    .then((response) => {
      if (response.ok) {
        return response;
      }
      throw new Error(`${response.statusText} : ${response.status}`);
    })
    .then(resp => resp.json());
}

function addSubscription({ idToken, instrument }) {
  return fetch(url.subscribe, {
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
  addSubscription
};
const url = {
  instruments: "https://stock-timing.appspot.com/rpc/instruments",
  candles: "https://stock-timing.appspot.com/rpc/candles",
  signals: "https://stock-timing.appspot.com/rpc/signals",
  subscribe: "https://stock-timing.appspot.com/rpc/auth/subscribe",
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
      'Content-Type': 'application/json',
      idToken,
    },
  })
    .then((response) => {
      if (response.ok) {
        return response;
      }
      throw new Error(`${response.statusText} : ${response.status}`);
    })
    .then(resp => resp.json())
}

export {
  retrieveInstruments,
  retrieveCandles,
  retrieveSignals,
  addSubscription
};
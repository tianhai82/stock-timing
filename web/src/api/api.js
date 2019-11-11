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

function retrieveSignals(instrumentID, period) {
  return fetch(`${url.signals}/${instrumentID}/period/${period}`)
    .then((response) => {
      if (response.ok) {
        return response;
      }
      throw new Error(`${response.statusText} : ${response.status}`);
    })
    .then(resp => resp.json());
}

function addSubscription({ idToken, instrument, period }) {
  return fetch(`${url.subscribe}/period/${period}`, {
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
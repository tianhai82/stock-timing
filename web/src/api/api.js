const url = {
  instruments: "https://stock-timing.appspot.com/rpc/instruments",
  candles: "https://stock-timing.appspot.com/rpc/candles",
  signals: "https://stock-timing.appspot.com/rpc/signals",
};

function retrieveInstruments() {
  return fetch(url.instruments)
    .then(resp => resp.json());
}

function retrieveCandles(instrumentID) {
  return fetch(`${url.candles}/${instrumentID}`)
    .then(resp => resp.json());
}

function retrieveSignals(instrumentID) {
  return fetch(`${url.signals}/${instrumentID}`)
    .then(resp => resp.json());
}

export {
  retrieveInstruments,
  retrieveCandles,
  retrieveSignals,
};
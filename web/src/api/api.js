const url = {
  instruments: "https://stock-timing.appspot.com/rpc/instruments",
  candles: "https://stock-timing.appspot.com/rpc/candles",
};

function retrieveInstruments() {
  return fetch(url.instruments, {
    method: "POST"
  }).then(resp => resp.json());
}

function retrieveCandles(instrumentID) {
  return fetch(`${url.candles}/${instrumentID}`, {
    method: "POST"
  }).then(resp => resp.json());
}

export {
  retrieveInstruments,
  retrieveCandles,
};
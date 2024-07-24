window.addEventListener("htmx:configRequest", (evt) => {
  if (evt.detail.verb == "get" && evt.detail.path == "/dynamicManyArts") {
    evt.detail.path = Alpine.store("manyArtsURL");
  }
});

// ---------------------------------------- //

let dontChange = false;
let defaultReq = {
  search: "",
  filter: {
    tags: [],
    minPrice: null,
    maxPrice: null,
  },
  sort: {
    by: "",
    asc: false,
  },
  pagination: {
    page: 1,
    limit: 20,
  },
};

document.addEventListener("alpine:init", () => {
  const body = document.querySelector("body");
  const event = new Event("findManyArts", { bubbles: true });

  Alpine.store("req", getParamReq() || copy(defaultReq));
  window.onpopstate = () => {
    dontChange = true;
    Alpine.store("req", getParamReq() || copy(defaultReq));
  };

  Alpine.store("manyArtsURL", "/api/arts");
  Alpine.store("total", 0);
  Alpine.effect(() => {
    JSON.stringify(Alpine.store("req"));
    JSON.stringify(Alpine.store("manyArtsURL"));

    body.dispatchEvent(event);
    htmx.trigger("body", "findManyArts");
  });
});

function pushHistory(req) {
  if (dontChange) {
    dontChange = false;
    return;
  }

  const url = new URL(window.location.href);
  url.searchParams.set("req", JSON.stringify(req));
  window.history.pushState(null, document.title, url.toString());
}

function getParamReq() {
  let params = new URLSearchParams(window.location.search);
  return (paramReq = JSON.parse(params.get("req")));
}

function requestBody() {
  let req = Alpine.store("req");

  if (req.filter.minPrice == null) {
    req.filter.minPrice = -1;
  }

  if (req.filter.maxPrice == null) {
    req.filter.maxPrice = -1;
  }

  return JSON.stringify(req);
}

function copy(obj) {
  let clone = Object.assign({}, obj);
  return clone;
}

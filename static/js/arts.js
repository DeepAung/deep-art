let req = {
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

  Alpine.store("req", req);
  Alpine.store("total", 0);
  Alpine.effect(() => {
    JSON.stringify(Alpine.store("req"));
    body.dispatchEvent(event);
    htmx.trigger("body", "findManyArts");
  });
});

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

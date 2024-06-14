document.addEventListener("alpine:init", () => {
	Alpine.store("req", {
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

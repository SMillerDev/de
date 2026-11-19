import { GeoJSON, LeafletMap, Path, TileLayer } from "leaflet";

const map = new LeafletMap("map", {
	zoomControl: false,
	scrollWheelZoom: false,
}).setView([49.574, 30.827], 6);
new TileLayer("https://tile.openstreetmap.org/{z}/{x}/{y}.png", {
	maxZoom: 19,
	attribution:
		'&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>',
}).addTo(map);

function updateMap() {
	clearMap();
	const value = document.querySelector("textarea.answer").value;
	try {
		const json = JSON.parse(value);
		const location = new GeoJSON(json).addTo(map);
		map.fitBounds(location.getBounds());
		document.querySelector("#invalid-helper").innerHTML = "";
		document.querySelector("textarea.answer").ariaInvalid = false;
		document.querySelector("input[type=submit]").disabled = false;
	} catch (error) {
		document.querySelector("#invalid-helper").innerHTML = error;
		document.querySelector("textarea.answer").ariaInvalid = true;
		document.querySelector("input[type=submit]").disabled = true;
	}
}

window.addEventListener("load", () => {
	updateMap();
	const field = document.querySelector("textarea.answer");
	console.log(field);
	if (field === undefined) {
		console.error("textarea not found!");
		return;
	}
	field.addEventListener("onkeyup", updateMap);
	field.addEventListener("focusout", updateMap);
});

function clearMap() {
	map.eachLayer((layer) => {
		if (layer instanceof Path) map.removeLayer(layer);
	});
}

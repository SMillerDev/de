import {
  LeafletMap,
  TileLayer,
  GeoJSON,
} from "leaflet";

const map = new LeafletMap("map", { zoomControl: false, scrollWheelZoom: false }).setView([49.574, 30.827], 6);
new TileLayer("https://tile.openstreetmap.org/{z}/{x}/{y}.png", {
  maxZoom: 19,
  attribution:
    '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>',
}).addTo(map);

function updateMap() {
  let value = document.querySelector("textarea.prompt").value;
  try {
    let json = JSON.parse(value);
    let location = new GeoJSON(json).addTo(map);
    document.querySelector("textarea.prompt").ariaInvalid = false;
    map.fitBounds(location.getBounds());
  } catch (error) {
    document.querySelector("#invalid-helper").innerHTML = error;
    document.querySelector("textarea.prompt").ariaInvalid = true;
  }
}

updateMap()
document.querySelector("textarea.prompt").addEventListener('onkeyup', updateMap)

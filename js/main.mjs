import {
  LeafletMap,
  TileLayer,
  GeoJSON,
  Marker,
  Polyline,
  Path,
} from "leaflet";

window.gameData = [];
window.question = 0;

const map = new LeafletMap("map").setView([49.574, 30.827], 6);
new TileLayer("https://tile.openstreetmap.org/{z}/{x}/{y}.png", {
  maxZoom: 19,
  attribution:
    '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>',
}).addTo(map);

// new TileLayer('https://server.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}', {
//     attribution: 'Tiles &copy; Esri &mdash; Source: Esri, i-cubed, USDA, USGS, AEX, GeoEye, Getmapping, Aerogrid, IGN, IGP, UPR-EGP, and the GIS User Community'
// }).addTo(map);

async function getData(path) {
  const url = `/api/v1/${path}`;
  try {
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error(`Response status: ${response.status}`);
    }

    const result = await response.json();
    console.log(result);
    return result;
  } catch (error) {
    console.error(error.message);
    return null;
  }
}

function distanceString(distance) {
  var unit = "meters";
  if (distance > 1000) {
    distance = distance / 1000;
    unit = "kilometers";
  }
  return `You clicked ${Math.floor(distance)} ${unit} away from the location`;
}

function onMapClick(e) {
  if (!window.gameData) {
    return;
  }
  console.debug("Clicked at " + e.latlng);

  let location = new GeoJSON(window.gameData[window.question].answer).addTo(
    map,
  );
  map.fitBounds(location.getBounds());
  if (location.getBounds().contains(e.latlng)) {
    document.querySelector(".prompt main p").textContent = "You found it! 🎉";
    window.gameData[window.question].result = { distance: 0 };

    gameEnded();
    return;
  }

  new Marker(e.latlng).addTo(map);
  let center = location.getBounds().getCenter();
  let distance = e.latlng.distanceTo(center);
  document.querySelector(".prompt main p").textContent = distanceString(distance);
  window.gameData[window.question].result = { distance: distance };

  new Polyline([e.latlng, center], { color: "grey", dashArray: "4" }).addTo(
    map,
  );

  gameEnded();
}

function gameEnded() {
  if (window.question < window.gameData.length) {
    document.querySelector(".next").style.display = "block";
    return false;
  }
  document.querySelector(".next").style.display = "none";
  document.querySelector(".prompt main p").textContent = "The end";
  return true;
}

function initQuestion() {
  map.fitBounds(baseLocation.getBounds());
  document.querySelector("button.next").style.display = "none";
  document.querySelector(".prompt main p").textContent =
  window.gameData[window.question].prompt;
}

async function startGame() {
  let data = await getData("today.json");
  if (data === null) {
    return;
  }
  window.gameData = data["data"];
  console.debug("started game");
  document.querySelector(".start").style.display = "none";
  initQuestion();
}

async function nextQuestion() {
  if (gameEnded()) {
    return;
  }
  window.question++;
  map.eachLayer(function (layer) {
    if (layer instanceof Marker) map.removeLayer(layer);
    if (layer instanceof Path) map.removeLayer(layer);
  });
  initQuestion();
}
let base = await getData("base.geo.json");
let baseLocation = new GeoJSON(base);
map.fitBounds(baseLocation.getBounds());
map.on("click", onMapClick);
document.querySelector("button.start").addEventListener("click", startGame);
document.querySelector("button.next").addEventListener("click", nextQuestion);

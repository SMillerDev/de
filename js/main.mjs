import {
  LeafletMap,
  TileLayer,
  GeoJSON,
  Marker,
  Polyline,
  Path,
  LatLng,
  Icon,
} from "leaflet";

window.gameData = {};
window.question = -1;

const mapIcon = new Icon({
    iconUrl: '/img/pin.png',
    // iconRetinaUrl: 'my-icon@2x.png',
    iconSize: [30, 30],
    iconAnchor: [15, 30],
});

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

    return await response.json();
  } catch (error) {
    console.error(error.message);
    return null;
  }
}

function toggleButton(className, show) {
  document.querySelector("button." + className).style.display = (show ? "block" : "none");
}

function setText(text) {
  document.querySelector(".prompt main p").innerHTML = text;
}

function distanceString(distance, name = 'the location') {
  if (distance === 0) {
    return `🎉 You found ${name}!`;
  }
  var unit = "meters";
  if (distance > 1000) {
    distance = distance / 1000;
    unit = "kilometers";
  }
  return `You clicked ${Math.floor(distance)} ${unit} away from ${name}`;
}

function getAnswers() {
  let answerString = localStorage.getItem('answers') || '{}';
  return JSON.parse(answerString);
}

function recordAnswer(distance, latlng) {
  let answers = getAnswers();
  fetch(`/api/v1/result?id=${getCurrentQuestionId()}&lat=${latlng.lat}&lng=${latlng.lng}`);
  answers[getCurrentQuestionId()] = { distance: distance, lat: latlng.lat, lng: latlng.lng };
  localStorage.setItem('answers', JSON.stringify(answers));
}

function drawLocation(location) {
  return new GeoJSON(location).addTo(map);
}

function drawPickDiff(target, pick) {
  new Marker(pick, {icon: mapIcon}).addTo(map);
  new Polyline([pick, target], { color: "grey", dashArray: "4" }).addTo(map);
}

function onMapClick(e) {
  if (document.querySelector("button.next").style.display !== 'none') { return }
  console.debug("Clicked at " + e.latlng);
  let location = drawLocation(getCurrentQuestion().answer);
  map.fitBounds(location.getBounds());
  if (location.getBounds().contains(e.latlng)) {
    setText(distanceString(0));
    recordAnswer(0, e.latlng);
    toggleButton('next', true);

    return;
  }

  let center = location.getBounds().getCenter();
  let distance = e.latlng.distanceTo(center);
  drawPickDiff(center, e.latlng);
  setText(distanceString(distance));
  recordAnswer(distance, e.latlng);
  toggleButton('next', true);
}

function gameEnded() {
  console.debug(`Question ${window.question + 1} / ${window.questionIds.length}`);
  if (window.question < (window.questionIds.length - 1)) {
    return false;
  }
  console.debug(`Marking end!`);
  clearMap();
  let answerStrings = [];
  let answers = getAnswers();
  for (let id in answers) {
    let gameAnswer = window.gameData[id].answer;
    answerStrings.push(`${distanceString(answers[id].distance, gameAnswer.properties.name)}<br/>`)
    let location = drawLocation(gameAnswer);
    drawPickDiff(location.getBounds().getCenter(), new LatLng(answers[id].lat, answers[id].lng));
  }
  map.fitBounds(baseLocation.getBounds());
  toggleButton('next', false);
  let answerValues = Object.values(answers);
  setText(`<em>The end, you got ${answerValues.filter(answer => answer.distance === 0).length}/${answerValues.length}</em><br/><br/> ${answerStrings.join('\n')}`);
  return true;
}

function getCurrentQuestionId() {
  return window.questionIds[window.question];
}
function getCurrentQuestion() {
  let key = getCurrentQuestionId();
  return window.gameData[key];
}

function clearMap() {
  map.eachLayer(function (layer) {
    if (layer instanceof Marker) map.removeLayer(layer);
    if (layer instanceof Path) map.removeLayer(layer);
  });
}

function initQuestion() {
  localStorage.setItem('lastPlay', new Date().toISOString());
  map.fitBounds(baseLocation.getBounds());
  toggleButton('next', false);
  setText(getCurrentQuestion().prompt);
}

async function onNextClick() {
  if (gameEnded()) {
    return;
  }
  window.question++;
  clearMap();
  initQuestion();
}

let base = await getData("base.geo.json");
let baseLocation = new GeoJSON(base);
map.fitBounds(baseLocation.getBounds());

let data = await getData("today.json");
if (data === null) {
  alert('failed to load data!')
}
window.gameData = data["data"];
window.questionIds = Object.getOwnPropertyNames(window.gameData);

let diff = window.questionIds.filter(x => !Object.getOwnPropertyNames(getAnswers()).includes(x));
if (diff.length === 0) {
  window.question = window.questionIds.length;
  gameEnded();
}

map.on("click", onMapClick);
document.querySelector("button.next").addEventListener("click", onNextClick);
localStorage.setItem('lastPlay', new Date().toISOString());

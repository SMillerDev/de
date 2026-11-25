import {
	GeoJSON,
	Icon,
	LatLng,
	LeafletMap,
	Marker,
	Path,
	Polyline,
  TileLayer,
	FeatureGroup
} from "leaflet";

window.gameData = {};
window.question = -1;

const mapIcon = new Icon({
	iconUrl: "/img/pin.png",
	// iconRetinaUrl: 'my-icon@2x.png',
	iconSize: [30, 30],
	iconAnchor: [15, 30],
});

const map = new LeafletMap("map", { zoomControl: false }).setView(
	[49.574, 30.827],
	6,
);

let layer = null;
if (localStorage.getItem("mode") !== "hard") {
  layer = new TileLayer("https://tiles.openstreetmap.org.ua/osm/{z}/{x}/{y}.png", {
    maxZoom: 19,
    attribution:
      '&copy; <a href="https://www.openstreetmap.org.ua/copyright">OpenStreetMap</a>',
  })
} else {
  layer = new TileLayer('https://server.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}', {
    attribution: 'Tiles &copy; Esri &mdash; Source: Esri, i-cubed, USDA, USGS, AEX, GeoEye, Getmapping, Aerogrid, IGN, IGP, UPR-EGP, and the GIS User Community'
  })
}
layer.addTo(map);

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

function toggleButton(show) {
	document.querySelector(`article.prompt div`).style.display = show
		? "block"
		: "none";
}

function setText(text) {
	document.querySelector(".prompt main p").innerHTML = text;
}

function distanceString(distance, name = "the location") {
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
	const answerString = localStorage.getItem("answers") || "{}";
	return JSON.parse(answerString);
}

function recordAnswer(distance, latlng) {
	const answers = getAnswers();
	fetch(
		`/api/v1/result?id=${getCurrentQuestionId()}&lat=${latlng.lat}&lng=${latlng.lng}`,
		{ method: "POST" },
	);
	answers[getCurrentQuestionId()] = {
		distance: distance,
		lat: latlng.lat,
		lng: latlng.lng,
	};
	localStorage.setItem("answers", JSON.stringify(answers));
}

function setProgress() {
	const progress = `${Object.keys(getAnswers()).length} / ${window.questionIds.length}`;
	console.debug(`Question ${progress}`);
	document.querySelector("nav li.progress").innerHTML = progress;

	var distance = 0;
	const answers = getAnswers();
	for (const id in answers) {
		distance = distance + answers[id].distance;
	}
	document.querySelector("nav li.score").style.display = "block";
	document.querySelector("nav li.score").innerHTML =
		`Distance: ${Math.floor(distance / 1000)} km`;
}

function drawLocation(location) {
	return new GeoJSON(location);
}

function drawPickDiff(target, pick) {
	let marker = new Marker(pick, { icon: mapIcon });
  let line = new Polyline([pick, target], { color: "grey", dashArray: "4" });
  return new FeatureGroup([line, marker]);
}

function onMapClick(e) {
	if (document.querySelector("article.prompt div").style.display !== "none") {
		return;
	}
	console.debug(`Clicked at ${e.latlng}`);
	const location = drawLocation(getCurrentQuestion().answer).addTo(map);
	if (location.getBounds().contains(e.latlng)) {
		setText(distanceString(0));
		recordAnswer(0, e.latlng);
		toggleButton(true);

		return;
	}

	const center = location.getBounds().getCenter();
	const distance = e.latlng.distanceTo(center);
  let answerLayer = drawPickDiff(center, e.latlng);
  answerLayer.addLayer(location);
  answerLayer.addTo(map);
  map.flyToBounds(answerLayer)
	setText(distanceString(distance));
	recordAnswer(distance, e.latlng);
	toggleButton(true);
}

function shouldRestart(date) {
	const dateStorage = localStorage.getItem("lastPlay");
	if (dateStorage == null) {
		return;
	}
	const playDate = new Date(`${date}T00:00:00Z`);
	const dataDate = new Date(dateStorage);
	if (playDate > dataDate) {
		localStorage.setItem("answers", "{}");
	}
}

function gameEnded() {
	if (window.question < window.questionIds.length - 1) {
		return false;
	}
  console.debug(`Marking end!`);
  document.querySelector("button.next").removeEventListener("click", onNextClick);
	clearMap();
	const answerStrings = [];
	const answers = getAnswers();
	for (const id in window.gameData) {
    const gameAnswer = window.gameData[id].answer;
    let wikiText = Object.hasOwn(window.gameData[id].wikis, 'en') ? ` - <a class="wiki secondary" href="${window.gameData[id].wikis.en}">i</a>` : '';
		answerStrings.push(
			`${window.questionIds.indexOf(id) + 1}. ${distanceString(answers[id].distance, gameAnswer.properties.name)} ${wikiText}<br/>`,
		);
		const location = drawLocation(gameAnswer);
		let answerLayer = drawPickDiff(
			location.getBounds().getCenter(),
			new LatLng(answers[id].lat, answers[id].lng),
    );
    answerLayer.addLayer(location)
    answerLayer.addTo(map);
	}
	map.fitBounds(baseLocation.getBounds());
	toggleButton(false);
	const answerValues = Object.values(answers);
	setText(
		`<em>The end, you got ${answerValues.filter((answer) => answer.distance === 0).length}/${answerValues.length}</em><br/><br/> ${answerStrings.join("\n")}`,
	);
	return true;
}

function getCurrentQuestionId() {
	return window.questionIds[window.question];
}
function getCurrentQuestion() {
	const key = getCurrentQuestionId();
	return window.gameData[key];
}

function clearMap() {
	map.eachLayer((layer) => {
    if (layer instanceof FeatureGroup) {
      layer.clearLayers();
      map.removeLayer(layer);
    }
	});
}

function initQuestion() {
	localStorage.setItem("lastPlay", new Date().toISOString());
	map.fitBounds(baseLocation.getBounds());
	toggleButton(false);
	setProgress();
	setText(getCurrentQuestion().prompts.en);
}

async function onNextClick() {
	if (gameEnded()) {
		return;
	}
	window.question++;
	clearMap();
	initQuestion();
}

const base = await getData("base.geo.json");
const baseLocation = new GeoJSON(base);
map.fitBounds(baseLocation.getBounds());

const data = await getData("today.json");
if (data === null) {
	alert("failed to load data!");
}
shouldRestart(data.playDate);
window.gameData = data.data;
window.questionIds = Object.getOwnPropertyNames(window.gameData);
if (window.questionIds.length === 0) {
	setText("No questions today!");
} else {
	const diff = window.questionIds.filter(
		(x) => !Object.getOwnPropertyNames(getAnswers()).includes(x),
	);
	if (diff.length === 0) {
		window.question = window.questionIds.length;
		gameEnded();
	}
}

function isEasy() {
  return localStorage.getItem("mode") !== "hard"
}

function onDifficultyClick() {
  let easy = isEasy();
  switch (easy) {
    case false:
      localStorage.setItem("mode", "easy");
      break;
    default:
      localStorage.setItem("mode", "hard");
      break;
  }
  location.reload();
}

setProgress();
map.on("click", onMapClick);
document.querySelector("button.difficulty").innerHTML = isEasy() ? "🌐" : "🌍"
document.querySelector("button.difficulty").addEventListener("click", onDifficultyClick);
document.querySelector("button.next").addEventListener("click", onNextClick);

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

function clearMap() {
	map.eachLayer((layer) => {
		if (layer instanceof Path) map.removeLayer(layer);
	});
}

async function osmApi(query) {
  let params = {
    q: query,
    format: 'json'
  };
  let url = `https://nominatim.openstreetmap.org/search?`;
  const response = await fetch(
    url + new URLSearchParams(params).toString(), {
      mode: 'cors',
      method: 'GET',
  });
  if (!response.ok) {
    return null;
  }

  return response.json();
}

async function getGeoJSON(type, id) {
  if (type !== 'way') { return ''; }
  let response = await fetch(`https://api.openstreetmap.org/api/0.6/${type}/${id}/full.json`);
  let data = await response.json();
  let coords = {};
  for (let elementKey in data.elements) {
    let element = data.elements[elementKey];
    if (element.type !== 'node') { continue; }
    coords[element.id] = [element.lon, element.lat];
  }
  let way = data.elements.filter(element => element.type === 'way')[0];
  let path = [];
  for (let nodeKey in way.nodes) {
    let nodeId = way.nodes[nodeKey];
    path.push(coords[nodeId]);
  }
  return {
    type: "Feature",
    properties: {
      osm_id: id,
      osm_type: type
    },
    geometry: {
      type: "Polygon",
      coordinates: [path]
    }
  };
}

async function wikiApi(query, lang) {
  query['format'] = 'json';
  query['origin'] = '*';
  let headers = {'Api-User-Agent': 'Hra-de admin/0.1 (https://smillerdev-de.sliplane.app/)'}
  let url = `https://${lang}.wikipedia.org/w/api.php?`;
  const response = await fetch(
    url + new URLSearchParams(query).toString(), {
      mode: 'cors',
      method: 'GET',
      headers: new Headers(headers),
  });
  if (!response.ok) {
    return null;
  }

  return response.json();
}

async function getExcerpt(title, lang) {
  let excerptData = await wikiApi({
    titles: title,
    prop: 'extracts',
    action: 'query',
    exsectionformat: 'plain',
    explaintext: 'true',
    exintro: 'true',
    exsentences: '2',
  }, lang);
  if (excerptData == null) {
    return;
  }

  return Object.values(excerptData.query.pages)[0].extract;
}

async function getAlternateTitles(name, lang) {
  let languageData = await wikiApi({
    titles: name,
    lllimit: 500,
    prop: 'langlinks',
    action: 'query',
  }, lang);
  if (languageData == null) {
    return;
  }
  let returnValue = {};
  returnValue[lang] = name;
  let links = Object.values(languageData.query.pages)[0].langlinks;
  for (let item in links) {
    returnValue[links[item].lang] = links[item]['*'];
  }

  return returnValue;
}

function setIfNotPresent(selector, value) {
  let existingValue = document.querySelector(selector).value;
  if (existingValue !== '') { return; }
  console.log(existingValue);
  document.querySelector(selector).value = value;
}

async function fillFields() {
  let name = document.querySelector("input.new-name").value;
  let titles = await getAlternateTitles(name, 'en');

  for (let langKey in window.supportedLanguages) {
    let lang = window.supportedLanguages[langKey];
    let excerpt = await getExcerpt(titles[lang], lang);
    if (excerpt !== null) {
      setIfNotPresent(`input[name=prompt-${lang}`, excerpt);
      setIfNotPresent(`input[name=wiki-${lang}`, `https://${lang}.wikipedia.org/wiki/${titles[lang].replaceAll(' ', '_')}`);
    }
  }

  let osm = await osmApi(name);
  let osmList = osm.filter(item => item.osm_type == 'way')
  if (osmList.length < 1) { return; }
  let geojson = await getGeoJSON(osmList[0].osm_type, osmList[0].osm_id);
  if (geojson !== '') {
    setIfNotPresent("textarea[name=answer]", JSON.stringify(geojson, null, 2));
    updateMap();
  }
}

updateMap();
const answerField = document.querySelector("textarea.answer");
answerField.addEventListener("onkeyup", updateMap);
answerField.addEventListener("focusout", updateMap);

const nameField = document.querySelector("input.new-name");
if (nameField !== null) {
  nameField.addEventListener("onkeyup", fillFields);
  nameField.addEventListener("focusout", fillFields);
}

import APIClient from "./apiClient.js";

export const map = L.map('map', {
    zoomControl: false
}).setView([6.673175, -1.565423], 16);

L.control.scale({
    position: 'topleft'
}).addTo(map);

L.control.zoom({
    position: 'bottomright'
}).addTo(map);

const markersContainer = [];
const polylineCoordinates = [];

function clearMarkers() {
    for (let m of markersContainer) {
        map.removeLayer(m);
    }
}

const clearButton = L.control({position: 'topright'});

clearButton.onAdd = () => {
    const div = L.DomUtil.create('div');
    div.innerHTML = '<button class="btn btn-danger clear-button">Clear Markers</button>';

    div.querySelector('.clear-button').addEventListener('click', clearMarkers);

    return div;
};

clearButton.addTo(map);

L.control.sidepanel('mySidepanelLeft', {
    tabsPosition: 'left',
    startTab: 'tab-1'
}).addTo(map);

L.tileLayer('http://{s}.google.com/vt/lyrs=s&x={x}&y={y}&z={z}', {
    maxZoom: 20,
    subdomains: ['mt0', 'mt1', 'mt2', 'mt3']
}).addTo(map);

export const submitForm = (routeUrl, formID) => {
    document.addEventListener('DOMContentLoaded', () => {
        const form = document.getElementById(formID);
        const spinner = document.getElementById('spinner');

        form.addEventListener('submit', (e) => {
            e.preventDefault();

            spinner.style.display = 'flex';

            const formData = new FormData(form);
            const jsonData = Object.fromEntries(formData.entries());
            const data = JSON.stringify(jsonData);

            APIClient(routeUrl, 'POST', data, result => {
                spinner.style.display = 'none';
                console.log(result);
                const data = result.paths.map(res => res.geom_geographic);
                const distance = result.distance;
                console.log(data);
                drawPath(data, distance);
            });
        });
    });
};

export const searchFilter = (searchInputID, dropdownListID, data) => {
    const searchInput = document.getElementById(searchInputID);
    const dropdownList = document.getElementById(dropdownListID);

    searchInput.addEventListener("input", function () {
        const value = searchInput.value.toLowerCase();
        dropdownList.innerHTML = "";

        const filteredData = data.filter(item => item.toLowerCase().includes(value));

        filteredData.forEach(item => {
            const li = document.createElement("li");
            li.classList.add("dropdown-item");
            li.textContent = item;
            li.setAttribute("role", "button");
            dropdownList.appendChild(li);

            li.addEventListener("click", () => {
                searchInput.value = item;
                dropdownList.classList.remove("show");
            });
        });

        if (value === '' || filteredData.length === 0) {
            dropdownList.classList.remove("show");
        } else {
            dropdownList.classList.add("show");
        }
    });
}

export const showAlert = (message, alertType) => {
    const alertDiv = document.createElement('div');
    alertDiv.classList.add('alert', `alert-${alertType}`, 'alert-dismissible', 'fade', 'show');
    alertDiv.setAttribute('role', 'alert');

    alertDiv.innerHTML = `
        ${message}
        <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
    `;

    document.body.appendChild(alertDiv);
}

function drawPath(data, distance) {
    var markerDelay = 10;

    data.forEach((wktStr, index) => {
        // setTimeout(() => {
        let wkt = new Wkt.Wkt();
        wkt.read(wktStr);

        let latLng = wkt.toObject().getLatLng();
        let leafletObj = L.circleMarker(latLng, {
            radius: 3,
            fillColor: '#3889bc',
            color: '#3889bc',
            weight: 1.5,
            opacity: 1,
            fillOpacity: 0.8,

        });

        polylineCoordinates.push(latLng);
        markersContainer.push(leafletObj);
        leafletObj.addTo(map);
        L.polyline(polylineCoordinates, {color: '#3889bc', weight: 1, opacity: 1.0, dashedArray: '2, 2'}).addTo(map);


        if (index === 0) {
            map.flyTo(leafletObj.getLatLng(), 15, {duration: 1});
        }
        // }, index * markerDelay);
    });

    setTimeout(() => {
        const group = L.featureGroup(markersContainer);
        map.fitBounds(group.getBounds(), {padding: [10, 10]});

        if (polylineCoordinates.length > 0) {
            var lastPoint = polylineCoordinates[polylineCoordinates.length - 1];
            var distancePopup = L.popup()
                .setLatLng(lastPoint)
                .setContent(`<div style="font-size: 16px; font-weight: bold; color: #333;">Approximately ${distance}m walk</div>`)
                .openOn(map);
        }
    }, data.length * markerDelay);
}

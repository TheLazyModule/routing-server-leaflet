import APIClient from "./apiClient.js";


export const map = L.map('map', {
    zoomControl: false
}).setView([6.673175, -1.565423], 20);
L.control.scale({
    position: 'topleft'
}).addTo(map);
L.control.zoom({
    position: 'bottomright'
}).addTo(map);


const markersContainer = [];

function clearMarkers() {

    for (let m of markersContainer) {
        map.removeLayer(m);
    }
}

const clearButton = L.control({position: 'topright'});

clearButton.onAdd = () => {
    const div = L.DomUtil.create('div');
    div.innerHTML = '<button class="btn btn-danger clear-button">Clear Markers</button>';

    // Wait until the button is added to the map
    div.querySelector('.clear-button').addEventListener('click', clearMarkers);

    return div;
};

clearButton.addTo(map);

const sidepanelLeft = L.control.sidepanel('mySidepanelLeft', {
    tabsPosition: 'left',
    startTab: 'tab-1'
}).addTo(map);

L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
}).addTo(map);

export const submitForm = (routeUrl, formID) => {
    document.addEventListener('DOMContentLoaded', () => {
        const form = document.getElementById(formID);

        form.addEventListener('submit', (e) => {
            e.preventDefault(); // Prevent the default form submission

            const formData = new FormData(form);
            const jsonData = Object.fromEntries(formData.entries());
            console.log(jsonData)
            const data = JSON.stringify(jsonData)


            APIClient(routeUrl, 'POST', data, result => {
                console.log(result)
                const data = result.paths.map(res => res.point_geom_geographic)
                console.log(data)
                drawPath(data)
            })
        });
    })
}

export const searchFilter = (searchInputID, dropdownListID, data) => {
    const searchInput = document.getElementById(searchInputID);
    const dropdownList = document.getElementById(dropdownListID);

    searchInput.addEventListener("input", function () {
        const value = searchInput.value.toLowerCase();
        dropdownList.innerHTML = ""; // Clear previous results

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

    // Add the alert to the document body
    document.body.appendChild(alertDiv);
}


function drawPath(data) {

    var polylineCoordinates = [];

    data.forEach((wktStr) => {
        let wkt = new Wkt.Wkt();
        wkt.read(wktStr);

        let leafletObj = wkt.toObject();

        leafletObj.addTo(map); // Adds the marker to the map
        markersContainer.push(leafletObj)


        if (leafletObj.getLatLng) {
            polylineCoordinates.push(leafletObj.getLatLng());
        }
    });

    // L.polyline(polylineCoordinates, {color: 'red'}).addTo(map);
}

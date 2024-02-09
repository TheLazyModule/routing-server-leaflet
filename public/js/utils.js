import APIClient from "./apiClient";

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
                var polylineCoordinates = [];

                data.forEach(function (wktStr) {
                    var wkt = new Wkt.Wkt();
                    wkt.read(wktStr);

                    var leafletObj = wkt.toObject();
                    leafletObj.addTo(map); // Adds the marker to the map

                    if (leafletObj.getLatLng) {
                        polylineCoordinates.push(leafletObj.getLatLng());
                    }
                });

                // var polyline = L.polyline(polylineCoordinates, {color: 'red'}).addTo(map);
            })
        });
    })
}

export const searchFilter = (searchInputID, dropdownListID, data) => {
    const searchInput = document.getElementById(searchInputID);
    const dropdownList = document.getElementById(dropdownListID);

    searchInput.addEventListener("input", function () {
        let value = searchInput.value.toLowerCase();
        dropdownList.innerHTML = ""; // Clear previous results

        let filteredData = data.filter(item => item.toLowerCase().includes(value));

        filteredData.forEach(item => {
            let li = document.createElement("li");
            li.classList.add("dropdown-item");
            li.textContent = item;
            li.setAttribute("role", "button");
            dropdownList.appendChild(li);

            li.addEventListener("click", function () {
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

export const setMap = () => {

    const map = L.map('map', {
        zoomControl: false
    }).setView([6.673175, -1.565423], 20);
    L.control.scale({
        position: 'topleft'
    }).addTo(map);
    L.control.zoom({
        position: 'bottomright'
    }).addTo(map);

    const sidepanelLeft = L.control.sidepanel('mySidepanelLeft', {
        tabsPosition: 'left',
        startTab: 'tab-1'
    }).addTo(map);

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
    }).addTo(map);
}
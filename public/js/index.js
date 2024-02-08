import APIClient from "./apiClient";

(function () {
    const api = APIClient('GET', '', result => {

        console.log(result.data)
    })
})

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

L.marker([6.673175, -1.565423], 20).addTo(map)
    .bindPopup('A pretty CSS3 popup.<br> Easily customizable.')
    .openPopup()
    .bindTooltip('A pretty CSS3 tooltip.<br> Easily customizable.');


document.addEventListener("DOMContentLoaded", () => {
    const searchInput = document.getElementById("searchInputFrom");
    const dropdownList = document.getElementById("dropdownListFrom");
    const items = ["Apple", "Banana", "Orange", "Mango", "Grape", "Strawberry"]; // Sample data


    searchInput.addEventListener("input", function () {
        let value = searchInput.value.toLowerCase();
        dropdownList.innerHTML = ""; // Clear previous results

        let filteredItems = items.filter(item => item.toLowerCase().includes(value));

        filteredItems.forEach(item => {
            let li = document.createElement("li");
            li.classList.add("dropdown-item");
            li.textContent = item;
            li.setAttribute("role", "button");
            dropdownList.appendChild(li);

            li.addEventListener("click", function () {
                searchInput.value = item; // Set input value to the selected item's text
                let dropdownElement = new bootstrap.Dropdown(searchInput);
                dropdownElement.hide();
            });
        });

        if (value === '' || filteredItems.length === 0) {
            dropdownList.classList.remove("show");
        } else {
            dropdownList.classList.add("show");
        }
    });
});

document.addEventListener("DOMContentLoaded", () => {
    const searchInput = document.getElementById("searchInputTo");
    const dropdownList = document.getElementById("dropdownListTo");
    const items = ["Apple", "Banana", "Orange", "Mango", "Grape", "Strawberry"]; // Sample data

    searchInput.addEventListener("input", function () {
        let value = searchInput.value.toLowerCase();
        dropdownList.innerHTML = ""; // Clear previous results

        let filteredItems = items.filter(item => item.toLowerCase().includes(value));

        filteredItems.forEach(item => {
            let li = document.createElement("li");
            li.classList.add("dropdown-item");
            li.textContent = item;
            li.setAttribute("role", "button");
            dropdownList.appendChild(li);

            li.addEventListener("click", function () {
                searchInput.value = item; // Set input value to the selected item's text
                let dropdownElement = new bootstrap.Dropdown(searchInput);
                dropdownElement.hide();
            });
        });

        if (value === '' || filteredItems.length === 0) {
            dropdownList.classList.remove("show");
        } else {
            dropdownList.classList.add("show");
        }
    });
});

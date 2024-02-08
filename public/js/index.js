import APIClient from "./apiClient.js";


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


APIClient('GET', result => {
    // const data = result.data.map(res => res.name)
    searchFilter("searchInputFrom", "dropdownListFrom", data)
    searchFilter("searchInputTo", "dropdownListTo", data)

})

const searchFilter = (searchInputID, dropdownListID, data) => {
    document.addEventListener("DOMContentLoaded", () => {
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
                    searchInput.value = item; // Set input value to the selected item's text
                    let dropdownElement = new bootstrap.Dropdown(searchInput);
                    dropdownElement.hide();
                });
            });

            if (value === '' || filteredData.length === 0) {
                dropdownList.classList.remove("show");
            } else {
                dropdownList.classList.add("show");
            }
        });
    });

}

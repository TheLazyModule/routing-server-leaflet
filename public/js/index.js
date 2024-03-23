import APIClient from "./apiClient.js";
import {searchFilter, submitForm} from "./utils.js";


APIClient('/all', 'GET', '', result => {
    let places = result.places.filter(res => {
        return res.name !== null
    })

    let buildings = result.buildings.filter(res => {
        return res.name !== null
    })

    const data = [...places.map(res => res.name), ...buildings.map(res => res.name)]

    searchFilter("searchInputFrom", "dropdownListFrom", data)
    searchFilter("searchInputTo", "dropdownListTo", data)

}, resError => {
})

submitForm('/all/route', 'form')



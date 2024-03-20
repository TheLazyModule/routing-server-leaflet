import APIClient from "./apiClient.js";
import {searchFilter, submitForm} from "./utils.js";


APIClient('/places', 'GET', '', result => {
    let data = result.filter(res => {
        return res.name !== null
    })
    data = data.map(res => res.name)
    // console.log(data)
    searchFilter("searchInputFrom2", "dropdownListFrom2", data)
    searchFilter("searchInputTo2", "dropdownListTo2", data)
}, error => {
    if (error)
        throw new Error(error.message)
})

APIClient('/buildings', 'GET', '', result => {
    let data = result.filter(res => {
        return res.name !== null
    })
    data = data.map(res => res.name)
    // console.log(data)
    searchFilter("searchInputFrom", "dropdownListFrom", data)
    searchFilter("searchInputTo", "dropdownListTo", data)

}, error => {
    if (error)
        throw new Error(error.message)
})

APIClient('/all', 'GET', '', result => {
    let places = result.places.filter(res => {
        return res.name !== null
    })

    let buildings = result.buildings.filter(res => {
        return res.name !== null
    })


}, resError => {
})

submitForm('/buildings/route', 'form')
submitForm('/places/route', 'form2')



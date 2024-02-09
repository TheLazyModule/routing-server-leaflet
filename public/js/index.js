import APIClient from "./apiClient.js";
import {searchFilter, setMap, submitForm} from "./utils";


setMap()


APIClient('/places', 'GET', '', result => {
    let data = result.filter(res => {
        return res.name !== null
    })
    data = data.map(res => res.name)
    console.log(data)
    searchFilter("searchInputFrom", "dropdownListFrom", data)
    searchFilter("searchInputTo", "dropdownListTo", data)

})

APIClient('/buildings', 'GET', '', result => {
    let data = result.filter(res => {
        return res.name !== null
    })
    data = data.map(res => res.name)
    console.log(data)
    searchFilter("searchInputFrom", "dropdownListFrom", data)
    searchFilter("searchInputTo", "dropdownListTo", data)

})

submitForm('/buildings/route', 'form')
submitForm('/places/route', 'form2')



import APIClient from "./apiClient.js";
import {searchFilter, setMap, submitForm} from "./utils.js";


setMap()


APIClient('/places', 'GET', '', result => {
    let data = result.filter(res => {
        return res.name !== null
    })
    data = data.map(res => res.name)
    console.log(data)
    searchFilter("searchInputFrom2", "dropdownListFrom2", data)
    searchFilter("searchInputTo2", "dropdownListTo2", data)
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



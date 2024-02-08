const APIClient = (method, resCallback) => $.ajax({
    url: '/buildings',
    type: method,
    dataType: 'json',
    complete: (xhr) => {
        let response = xhr.responseJSON;
        console.log(response)
        console.log(response)
        resCallback(response.results)
    },
    error: (error) => {
        console.error('Error:', error);
    }
});

export default APIClient;